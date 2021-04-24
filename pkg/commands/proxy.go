package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
	shlex "github.com/flynn/go-shlex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"moul.io/assh/v2/pkg/config"
	"moul.io/assh/v2/pkg/ratelimit"
)

type contextKey string

type gatewayErrorMsg struct {
	gateway string
	err     zap.Field
}

var syncContextKey contextKey = "sync"

var proxyCommand = &cobra.Command{
	Use:     "connect",
	Short:   "Connect to host SSH socket, used by ProxyCommand",
	Example: "Argument is a host.",
	Hidden:  true,
	RunE:    runProxyCommand,
}

// nolint:gochecknoinits
func init() {
	proxyCommand.Flags().BoolP("no-rewrite", "", false, "Do not automatically rewrite outdated configuration")
	proxyCommand.Flags().IntP("port", "p", 0, "SSH destination port")
	proxyCommand.Flags().BoolP("dry-run", "", false, "Only show how assh would connect but don't actually do it")
	_ = viper.BindPFlags(proxyCommand.Flags())
}

func runProxyCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("assh: \"connect\" requires 1 argument. See 'assh connect --help'")
	}

	target := args[0]
	logger().Debug("initializing proxy", zap.String("target", target))

	// dry-run option
	// Setting the 'ASSH_DRYRUN=1' environment variable,
	// so 'assh' can use gateways using sub-SSH commands.
	if viper.GetBool("dry-run") {
		if err := os.Setenv("ASSH_DRYRUN", "1"); err != nil {
			return errors.Wrap(err, "failed to configure environment")
		}
	}
	dryRun := os.Getenv("ASSH_DRYRUN") == "1"

	conf, err := config.Open(viper.GetString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	if err = conf.LoadKnownHosts(); err != nil {
		logger().Debug("Failed to load assh known_hosts", zap.Error(err))
	}

	automaticRewrite := !viper.GetBool("no-rewrite")
	isOutdated, err2 := conf.IsConfigOutdated(target)
	if err2 != nil {
		logger().Warn("Cannot check if ~/.ssh/config is outdated", zap.Error(err))
	}
	if isOutdated {
		if automaticRewrite {
			// BeforeConfigWrite
			type configWriteHookArgs struct {
				SSHConfigPath string
			}
			hookArgs := configWriteHookArgs{
				SSHConfigPath: conf.SSHConfigPath(),
			}

			logger().Debug("Calling BeforeConfigWrite hooks")
			if drivers, err := conf.Defaults.Hooks.BeforeConfigWrite.InvokeAll(hookArgs); err != nil {
				logger().Error("BeforeConfigWrite hook failed", zap.Error(err))
			} else {
				defer drivers.Close()
			}

			// Save
			logger().Debug("The configuration file is outdated, rebuilding it before calling ssh")
			logger().Warn("'~/.ssh/config' has been rewritten.  SSH needs to be restarted.  See https://github.com/moul/assh/issues/122 for more information.")
			logger().Debug("Saving SSH config")
			if err := conf.SaveSSHConfig(); err != nil {
				return errors.Wrap(err, "failed to save SSH config file")
			}

			// AfterConfigWrite
			logger().Debug("Calling AfterConfigWrite hooks")
			if drivers, err := conf.Defaults.Hooks.AfterConfigWrite.InvokeAll(hookArgs); err != nil {
				logger().Error("AfterConfigWrite hook failed", zap.Error(err))
			} else {
				defer drivers.Close()
			}
		} else {
			logger().Warn("The configuration file is outdated; you need to run `assh config build --no-automatic-rewrite > ~/.ssh/config` to stay updated")
		}
	}

	// FIXME: handle complete host with json

	host, err := computeHost(target, viper.GetInt("port"), conf)
	if err != nil {
		return errors.Wrapf(err, "Failed to get host %q", target)
	}
	var w bytes.Buffer
	if err := host.WriteSSHConfigTo(&w); err != nil {
		return errors.Wrap(err, "failed to write ssh config")
	}
	logger().Debug("generated ssh config file", zap.String("buffer", w.String()))

	hostJSON, err2 := json.Marshal(host)
	if err2 != nil {
		logger().Warn("Failed to marshal host", zap.Error(err2))
	} else {
		logger().Debug("Host", zap.String("host", string(hostJSON)))
	}

	logger().Debug("Proxying")
	return proxy(host, conf, dryRun)
}

// nolint:unparam
func computeHost(dest string, portOverride int, conf *config.Config) (*config.Host, error) {
	host := conf.GetHostSafe(dest)

	if portOverride > 0 {
		host.Port = strconv.Itoa(portOverride)
	}

	return host, nil
}

func expandSSHTokens(tokenized string, host *config.Host) string {
	result := tokenized

	// OpenSSH Token Cheatsheet (stolen directly from the man pages)
	//
	// %%    A literal `%'.
	// %C    Shorthand for %l%h%p%r.
	// %d    Local user's home directory.
	// %h    The remote hostname.
	// %i    The local user ID.
	// %L    The local hostname.
	// %l    The local hostname, including the domain name.
	// %n    The original remote hostname, as given on the command line.
	// %p    The remote port.
	// %r    The remote username.
	// %u    The local username.

	// TODO: Expansion of strings like "%%C" and "%C" are equivalent due to the
	//       order that tokens are evaluated.  Should look at how OpenSSH implements
	//       the tokenization behavior.

	// Expand a home directory ~.  Assume nobody is using
	// the ~otheruser syntax.
	homedir := os.ExpandEnv("$HOME")

	if result[0] == '~' {
		result = strings.Replace(result, "~", homedir, 1)
	}
	result = strings.ReplaceAll(result, "%d", homedir)

	result = strings.ReplaceAll(result, "%%", "%")
	result = strings.ReplaceAll(result, "%C", "%l%h%p%r")
	result = strings.ReplaceAll(result, "%h", host.Name())
	result = strings.ReplaceAll(result, "%i", strconv.Itoa(os.Geteuid()))
	result = strings.ReplaceAll(result, "%p", host.Port)

	if hostname, err := os.Hostname(); err == nil {
		result = strings.ReplaceAll(result, "%L", hostname)
	} else {
		result = strings.ReplaceAll(result, "%L", "hostname")
	}

	if host.User != "" {
		result = strings.ReplaceAll(result, "%r", host.User)
	} else {
		if userdata, err := user.Current(); err == nil {
			result = strings.ReplaceAll(result, "%r", userdata.Username)
		} else {
			result = strings.ReplaceAll(result, "%r", "username")
		}
	}

	return result
}

func prepareHostControlPath(host *config.Host) error {
	if !config.BoolVal(host.ControlMasterMkdir) || ("none" == host.ControlPath || "" == host.ControlPath) {
		return nil
	}

	controlPath := expandSSHTokens(host.ControlPath, host)
	controlPathDir := path.Dir(controlPath)
	logger().Debug("Creating control path", zap.String("path", controlPathDir))
	return os.MkdirAll(controlPathDir, 0700)
}

func proxy(host *config.Host, conf *config.Config, dryRun bool) error {
	if err := prepareHostControlPath(host.Clone()); err != nil {
		return errors.Wrap(err, "failed to prepare host control-path")
	}

	if len(host.Gateways) > 0 {
		logger().Debug("Trying gateways", zap.String("gateways", strings.Join(host.Gateways, ", ")))
		var gatewayErrors []gatewayErrorMsg
		for _, gateway := range host.Gateways {
			log.Println(gateway)
			if gateway == "direct" {
				if err := proxyDirect(host, dryRun); err != nil {
					gatewayErrors = append(gatewayErrors, gatewayErrorMsg{
						gateway: "direct", err: zap.Error(err)})
				} else {
					return nil
				}
			} else {
				hostCopy := host.Clone()
				gatewayHost := conf.GetGatewaySafe(gateway)

				if err := prepareHostControlPath(hostCopy); err != nil {
					return errors.Wrap(err, "failed to prepare host control-path")
				}

				// FIXME: dynamically add "-v" flags

				var command string

				// FIXME: detect ssh client version and use netcat if too old
				// for now, the workaround is to configure the ProxyCommand of the host to "nc %h %p"

				if err := hostPrepare(hostCopy, gateway); err != nil {
					return errors.Wrap(err, "failed to prepare host for gateway")
				}

				if hostCopy.ProxyCommand != "" {
					command = "ssh %name -- " + hostCopy.ExpandString(hostCopy.ProxyCommand, gateway)
				} else {
					command = hostCopy.ExpandString("ssh -W %h:%p ", "") + "%name"
				}

				logger().Debug(
					"Using gateway",
					zap.String("gateway", gateway),
					zap.String("command", command),
				)
				if err := runProxy(gatewayHost, command, dryRun); err != nil {
					gatewayErrors = append(gatewayErrors, gatewayErrorMsg{
						gateway: gateway, err: zap.Error(err)})
				} else {
					return nil
				}
			}
		}
		if len(gatewayErrors) > 0 {
			for _, errMsg := range gatewayErrors {
				conType := "gateway"
				if errMsg.gateway == "direct" {
					conType = "connection"
				}
				logger().Error(
					fmt.Sprintf("Failed to use '%s' %s with error:",
						errMsg.gateway, conType), errMsg.err)
			}
		}
		return errors.New("no such available gateway")
	}

	logger().Debug("Connecting without gateway")
	return proxyDirect(host, dryRun)
}

func proxyDirect(host *config.Host, dryRun bool) error {
	if host.ProxyCommand != "" {
		return runProxy(host, host.ProxyCommand, dryRun)
	}
	return proxyGo(host, dryRun)
}

func runProxy(host *config.Host, command string, dryRun bool) error {
	command = host.ExpandString(command, "")
	logger().Debug("ProxyCommand", zap.String("command", command))
	args, err := shlex.Split(command)
	if err != nil {
		return err
	}

	if dryRun {
		return fmt.Errorf("dry-run: Execute %s", args)
	}

	spawn := exec.Command(args[0], args[1:]...) // #nosec
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	return spawn.Run()
}

func hostPrepare(host *config.Host, gateway string) error {
	if host.HostName == "" {
		host.HostName = host.Name()
	}

	if len(host.ResolveNameservers) > 0 {
		logger().Debug(
			"Resolving host",
			zap.String("hostname", host.HostName),
			zap.String("nameservers", strings.Join(host.ResolveNameservers, ", ")),
		)
		// FIXME: resolve using custom dns server
		results, err := net.LookupAddr(host.HostName)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			host.HostName = results[0]
		}
		logger().Debug("Resolved host", zap.String("hostname", host.HostName))
	}

	if host.ResolveCommand != "" {
		command := host.ExpandString(host.ResolveCommand, gateway)
		logger().Debug(
			"Resolving host",
			zap.String("hostname", host.HostName),
			zap.String("resolve-command", host.ResolveCommand),
		)

		args, err := shlex.Split(command)
		if err != nil {
			return err
		}

		cmd := exec.Command(args[0], args[1:]...) // #nosec
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "failed to run resolve-command")
		}

		host.HostName = strings.TrimSpace(stdout.String())
		logger().Debug("Resolved host", zap.String("hostname", host.HostName))
	}
	return nil
}

type exportReadWrite struct {
	written uint64
	err     error
}

// ConnectionStats contains network and timing informations about a connection
type ConnectionStats struct {
	WrittenBytes            uint64
	WrittenBytesHuman       string
	CreatedAt               time.Time
	ConnectedAt             time.Time
	DisconnectedAt          time.Time
	ConnectionDuration      time.Duration
	ConnectionDurationHuman string
	AverageSpeed            float64
	AverageSpeedHuman       string
}

func (c *ConnectionStats) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logger().Error("failed to marshal ConnectionStats", zap.Error(err))
		return ""
	}
	return string(b)
}

// ConnectHookArgs is the struture sent to the hooks and used in Go templates by the hook drivers
type ConnectHookArgs struct {
	Host  *config.Host
	Stats *ConnectionStats
	Error string
}

func (c ConnectHookArgs) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logger().Error("failed to marshal ConnectHookArgs", zap.Error(err))
		return ""
	}
	return string(b)
}

func proxyGo(host *config.Host, dryRun bool) error {
	stats := ConnectionStats{
		CreatedAt: time.Now(),
	}
	connectHookArgs := ConnectHookArgs{
		Host:  host,
		Stats: &stats,
	}

	logger().Debug("Preparing host object")
	if err := hostPrepare(host, ""); err != nil {
		return errors.Wrap(err, "failed to prepare host")
	}

	if dryRun {
		return fmt.Errorf("dry-run: Golang native TCP connection to '%s:%s'", host.HostName, host.Port)
	}

	// BeforeConnect hook
	logger().Debug("Calling BeforeConnect hooks")
	if drivers, err := host.Hooks.BeforeConnect.InvokeAll(connectHookArgs); err != nil {
		logger().Error("BeforeConnect hook failed", zap.Error(err))
	} else {
		defer drivers.Close()
	}

	logger().Debug("Connecting to host", zap.String("hostname", host.HostName), zap.String("port", host.Port))

	// use GatewayConnectTimeout, fallback on ConnectTimeout
	timeout := host.GatewayConnectTimeout
	if host.ConnectTimeout != 0 {
		timeout = host.ConnectTimeout
	}
	if timeout < 0 { // set to 0 to disable
		timeout = 0
	}
	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%s", host.HostName, host.Port),
		time.Duration(timeout)*time.Second,
	)
	if err != nil {
		// OnConnectError hook
		connectHookArgs.Error = err.Error()
		logger().Debug("Calling OnConnectError hooks")
		if drivers, err := host.Hooks.OnConnectError.InvokeAll(connectHookArgs); err != nil {
			logger().Error("OnConnectError hook failed", zap.Error(err))
		} else {
			defer drivers.Close()
		}

		return errors.Wrap(err, "failed to dial")
	}
	logger().Debug(
		"Connected",
		zap.String("hostname", host.HostName),
		zap.String("port", host.Port),
	)
	stats.ConnectedAt = time.Now()

	// OnConnect hook
	logger().Debug("Calling OnConnect hooks")
	if drivers, err := host.Hooks.OnConnect.InvokeAll(connectHookArgs); err != nil {
		logger().Error("OnConnect hook failed", zap.Error(err))
	} else {
		defer drivers.Close()
	}

	// Ignore SIGHUP
	signal.Ignore(syscall.SIGHUP)

	waitGroup := sync.WaitGroup{}
	result := exportReadWrite{}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, syncContextKey, &waitGroup)

	waitGroup.Add(2)

	var reader io.Reader
	var writer io.Writer
	reader = conn
	writer = conn
	if host.RateLimit != "" {
		bytes, err := humanize.ParseBytes(host.RateLimit)
		if err != nil {
			return errors.Wrap(err, "failed to parse rate limit configuration")
		}
		limit := rate.Limit(float64(bytes))
		limiter := rate.NewLimiter(limit, int(bytes))
		reader = ratelimit.NewReader(conn, limiter)
		writer = ratelimit.NewWriter(conn, limiter)
	}

	c1 := readAndWrite(ctx, reader, os.Stdout)
	c2 := readAndWrite(ctx, os.Stdin, writer)
	select {
	case result = <-c1:
		stats.WrittenBytes = result.written
	case result = <-c2:
	}
	if result.err != nil && result.err == io.EOF {
		result.err = nil
	}

	if err := conn.Close(); err != nil {
		return err
	}
	cancel()
	waitGroup.Wait()
	select {
	case res := <-c1:
		stats.WrittenBytes = res.written
	default:
	}

	stats.DisconnectedAt = time.Now()
	stats.ConnectionDuration = stats.DisconnectedAt.Sub(stats.ConnectedAt)
	averageSpeed := float64(stats.WrittenBytes) / stats.ConnectionDuration.Seconds()
	// round duraction
	stats.ConnectionDuration = ((stats.ConnectionDuration + time.Second/2) / time.Second) * time.Second
	stats.AverageSpeed = math.Ceil(averageSpeed*1000) / 1000
	// human
	stats.WrittenBytesHuman = humanize.Bytes(stats.WrittenBytes)
	connectionDurationHuman := humanize.RelTime(stats.DisconnectedAt, stats.ConnectedAt, "", "")
	stats.ConnectionDurationHuman = strings.ReplaceAll(connectionDurationHuman, "now", "0 sec")
	stats.AverageSpeedHuman = humanize.Bytes(uint64(stats.AverageSpeed)) + "/s"

	// OnDisconnect hook
	logger().Debug("Calling OnDisconnect hooks")
	if drivers, err := host.Hooks.OnDisconnect.InvokeAll(connectHookArgs); err != nil {
		logger().Error("OnDisconnect hook failed", zap.Error(err))
	} else {
		defer drivers.Close()
	}

	logger().Debug(
		"Connection finished",
		zap.Uint64("bytes written", stats.WrittenBytes),
		zap.Error(result.err),
	)
	return result.err
}

func readAndWrite(ctx context.Context, r io.Reader, w io.Writer) <-chan exportReadWrite {
	buff := make([]byte, 1024)
	c := make(chan exportReadWrite, 1)

	go func() {
		defer ctx.Value(syncContextKey).(*sync.WaitGroup).Done()

		export := exportReadWrite{}
		for {
			select {
			case <-ctx.Done():
				c <- export
				return
			default:
				nr, err := r.Read(buff)
				if err != nil {
					export.err = err
					c <- export
					return
				}
				if nr > 0 {
					wr, err := w.Write(buff[:nr])
					if err != nil {
						export.err = err
						c <- export
						return
					}
					if wr > 0 {
						export.written += uint64(wr)
					}
				}
			}
		}
	}()
	return c
}

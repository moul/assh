package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"

	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/moul/advanced-ssh-config/pkg/logger"
	"github.com/moul/advanced-ssh-config/pkg/ratelimit"
)

type contextKey string

var syncContextKey contextKey = "sync"

func cmdProxy(c *cli.Context) error {
	logger.Logger.Debugf("assh args: %s", c.Args())

	if len(c.Args()) < 1 {
		logger.Logger.Fatalf("assh: \"connect\" requires 1 argument. See 'assh connect --help'.")
	}

	// dry-run option
	// Setting the 'ASSH_DRYRUN=1' environment variable,
	// so 'assh' can use gateways using sub-SSH commands.
	if c.Bool("dry-run") {
		if err := os.Setenv("ASSH_DRYRUN", "1"); err != nil {
			logger.Logger.Fatalf("Cannot set env var: %v", err)
		}
	}
	dryRun := os.Getenv("ASSH_DRYRUN") == "1"

	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		logger.Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	if err = conf.LoadKnownHosts(); err != nil {
		logger.Logger.Debugf("Failed to load assh known_hosts: %v", err)
	}

	target := c.Args()[0]

	automaticRewrite := !c.Bool("no-rewrite")
	isOutdated, err2 := conf.IsConfigOutdated(target)
	if err2 != nil {
		logger.Logger.Warnf("Cannot check if ~/.ssh/config is outdated: %v", err)
	} else if isOutdated {
		if automaticRewrite {
			// BeforeConfigWrite
			type configWriteHookArgs struct {
				SSHConfigPath string
			}
			hookArgs := configWriteHookArgs{
				SSHConfigPath: conf.SSHConfigPath(),
			}
			logger.Logger.Debugf("Calling BeforeConfigWrite hooks")
			beforeConfigWriteDrivers, err3 := conf.Defaults.Hooks.BeforeConfigWrite.InvokeAll(hookArgs)
			if err3 != nil {
				logger.Logger.Errorf("BeforeConfigWrite hook failed: %v", err3)
			}
			defer beforeConfigWriteDrivers.Close()

			// Save
			logger.Logger.Debugf("The configuration file is outdated, rebuilding it before calling ssh")
			logger.Logger.Warnf("'~/.ssh/config' has been rewritten.  SSH needs to be restarted.  See https://github.com/moul/advanced-ssh-config/issues/122 for more information.")
			logger.Logger.Debugf("Saving SSH config")
			err3 = conf.SaveSSHConfig()
			if err3 != nil {
				logger.Logger.Fatalf("Cannot save SSH config file: %v", err3)
			}

			// AfterConfigWrite
			logger.Logger.Debugf("Calling AfterConfigWrite hooks")
			afterConfigWriteDrivers, err3 := conf.Defaults.Hooks.AfterConfigWrite.InvokeAll(hookArgs)
			if err3 != nil {
				logger.Logger.Errorf("AfterConfigWrite hook failed: %v", err3)
			}
			defer afterConfigWriteDrivers.Close()

		} else {
			logger.Logger.Warnf("The configuration file is outdated; you need to run `assh config build --no-automatic-rewrite > ~/.ssh/config` to stay updated")
		}
	}

	// FIXME: handle complete host with json

	host, err := computeHost(target, c.Int("port"), conf)
	if err != nil {
		logger.Logger.Fatalf("Cannot get host '%s': %v", target, err)
	}
	w := logger.Logger.Writer()
	if err3 := host.WriteSSHConfigTo(w); err3 != nil {
		logger.Logger.Fatalf("Cannot write ssh config: %v", err3)
	}
	if err3 := w.Close(); err3 != nil {
		logger.Logger.Fatalf("Failed to close file: %v", err3)
	}

	hostJSON, err2 := json.Marshal(host)
	if err2 != nil {
		logger.Logger.Warnf("Failed to marshal host: %v", err2)
	}
	logger.Logger.Debugf("Host: %s", hostJSON)

	logger.Logger.Debugf("Proxying")
	err = proxy(host, conf, dryRun)
	if err != nil {
		logger.Logger.Fatalf("Proxy error: %v", err)
	}

	return nil
}

func computeHost(dest string, portOverride int, conf *config.Config) (*config.Host, error) {
	host := conf.GetHostSafe(dest)

	if portOverride > 0 {
		host.Port = strconv.Itoa(portOverride)
	}

	return host, nil
}

func expandSSHTokens(tokenized string, host *config.Host, gateway *config.Host) string {
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
	result = strings.Replace(result, "%d", homedir, -1)

	result = strings.Replace(result, "%%", "%", -1)
	result = strings.Replace(result, "%C", "%l%h%p%r", -1)
	result = strings.Replace(result, "%h", path.Join(host.Name(), gateway.Name()), -1)
	result = strings.Replace(result, "%i", strconv.Itoa(os.Geteuid()), -1)
	result = strings.Replace(result, "%p", host.Port, -1)

	if hostname, err := os.Hostname(); err == nil {
		result = strings.Replace(result, "%L", hostname, -1)
	} else {
		result = strings.Replace(result, "%L", "hostname", -1)
	}

	if host.User != "" {
		result = strings.Replace(result, "%r", host.User, -1)
	} else {
		if userdata, err := user.Current(); err == nil {
			result = strings.Replace(result, "%r", userdata.Username, -1)
		} else {
			result = strings.Replace(result, "%r", "username", -1)
		}
	}

	return result
}

func prepareHostControlPath(host, gateway *config.Host) error {
	if !config.BoolVal(host.ControlMasterMkdir) && ("none" == host.ControlPath || "" == host.ControlPath) {
		return nil
	}

	controlPath := expandSSHTokens(host.ControlPath, host, gateway)
	controlPathDir := path.Dir(controlPath)
	logger.Logger.Debugf("Creating control path: %s", controlPathDir)
	return os.MkdirAll(controlPathDir, 0700)
}

func proxy(host *config.Host, conf *config.Config, dryRun bool) error {

	emptygw := config.Host{}
	err := prepareHostControlPath(host.Clone(), emptygw.Clone())
	if err != nil {
		return err
	}

	if len(host.Gateways) > 0 {
		logger.Logger.Debugf("Trying gateways: %s", host.Gateways)
		for _, gateway := range host.Gateways {
			if gateway == "direct" {
				err = proxyDirect(host, dryRun)
				if err != nil {
					logger.Logger.Errorf("Failed to use 'direct' connection: %v", err)
				} else {
					return nil
				}
			} else {
				hostCopy := host.Clone()
				gatewayHost := conf.GetGatewaySafe(gateway)

				err = prepareHostControlPath(hostCopy, gatewayHost)
				if err != nil {
					return err
				}

				// FIXME: dynamically add "-v" flags

				var command string

				// FIXME: detect ssh client version and use netcat if too old
				// for now, the workaround is to configure the ProxyCommand of the host to "nc %h %p"

				if err = hostPrepare(hostCopy, gateway); err != nil {
					return err
				}

				if hostCopy.ProxyCommand != "" {
					command = "ssh %name -- " + hostCopy.ExpandString(hostCopy.ProxyCommand, gateway)
				} else {
					command = hostCopy.ExpandString("ssh -W %h:%p ", "") + "%name"
				}

				logger.Logger.Debugf("Using gateway '%s': %s", gateway, command)
				err = proxyCommand(gatewayHost, command, dryRun)
				if err == nil {
					return nil
				}
				logger.Logger.Errorf("Cannot use gateway '%s': %v", gateway, err)
			}
		}
		return fmt.Errorf("No such available gateway")
	}

	logger.Logger.Debugf("Connecting without gateway")
	return proxyDirect(host, dryRun)
}

func proxyDirect(host *config.Host, dryRun bool) error {
	if host.ProxyCommand != "" {
		return proxyCommand(host, host.ProxyCommand, dryRun)
	}
	return proxyGo(host, dryRun)
}

func proxyCommand(host *config.Host, command string, dryRun bool) error {
	command = host.ExpandString(command, "")
	logger.Logger.Debugf("ProxyCommand: %s", command)
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
		logger.Logger.Debugf("Resolving host: '%s' using nameservers %s", host.HostName, host.ResolveNameservers)
		// FIXME: resolve using custom dns server
		results, err := net.LookupAddr(host.HostName)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			host.HostName = results[0]
		}
		logger.Logger.Debugf("Resolved host is: %s", host.HostName)
	}

	if host.ResolveCommand != "" {
		command := host.ExpandString(host.ResolveCommand, gateway)
		logger.Logger.Debugf("Resolving host: %q using command: %q", host.HostName, host.ResolveCommand)

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
			logger.Logger.Errorf("ResolveCommand failed: %s", stderr.String())
			return err
		}

		host.HostName = strings.TrimSpace(stdout.String())
		logger.Logger.Debugf("Resolved host is: %s", host.HostName)
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

// ConnectHookArgs is the struture sent to the hooks and used in Go templates by the hook drivers
type ConnectHookArgs struct {
	Host  *config.Host
	Stats *ConnectionStats
	Error error
}

func proxyGo(host *config.Host, dryRun bool) error {
	stats := ConnectionStats{
		CreatedAt: time.Now(),
	}
	connectHookArgs := ConnectHookArgs{
		Host:  host,
		Stats: &stats,
	}

	logger.Logger.Debugf("Preparing host object")
	if err := hostPrepare(host, ""); err != nil {
		return err
	}

	if dryRun {
		return fmt.Errorf("dry-run: Golang native TCP connection to '%s:%s'", host.HostName, host.Port)
	}

	// BeforeConnect hook
	logger.Logger.Debugf("Calling BeforeConnect hooks")
	beforeConnectDrivers, err := host.Hooks.BeforeConnect.InvokeAll(connectHookArgs)
	if err != nil {
		logger.Logger.Errorf("BeforeConnect hook failed: %v", err)
	}
	defer beforeConnectDrivers.Close()

	logger.Logger.Debugf("Connecting to %s:%s", host.HostName, host.Port)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host.HostName, host.Port), time.Duration(host.ConnectTimeout)*time.Second)
	if err != nil {
		// OnConnectError hook
		connectHookArgs.Error = err
		logger.Logger.Debugf("Calling OnConnectError hooks")
		onConnectErrorDrivers, err2 := host.Hooks.OnConnectError.InvokeAll(connectHookArgs)
		if err2 != nil {
			logger.Logger.Errorf("OnConnectError hook failed: %v", err2)
		}
		defer onConnectErrorDrivers.Close()

		return err
	}
	logger.Logger.Debugf("Connected to %s:%s", host.HostName, host.Port)
	stats.ConnectedAt = time.Now()

	// OnConnect hook
	logger.Logger.Debugf("Calling OnConnect hooks")
	onConnectDrivers, err := host.Hooks.OnConnect.InvokeAll(connectHookArgs)
	if err != nil {
		logger.Logger.Errorf("OnConnect hook failed: %v", err)
	}
	defer onConnectDrivers.Close()

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
		bytes, err2 := humanize.ParseBytes(host.RateLimit)
		if err2 != nil {
			return err2
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

	if err2 := conn.Close(); err2 != nil {
		return err2
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
	stats.ConnectionDurationHuman = strings.Replace(connectionDurationHuman, "now", "0 sec", -1)
	stats.AverageSpeedHuman = humanize.Bytes(uint64(stats.AverageSpeed)) + "/s"

	// OnDisconnect hook
	logger.Logger.Debugf("Calling OnDisconnect hooks")
	onDisconnectDrivers, err := host.Hooks.OnDisconnect.InvokeAll(connectHookArgs)
	if err != nil {
		logger.Logger.Errorf("OnDisconnect hook failed: %v", err)
	}
	defer onDisconnectDrivers.Close()

	logger.Logger.Debugf("Byte written %v", stats.WrittenBytes)
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

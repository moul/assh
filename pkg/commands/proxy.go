package commands

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"
	shlex "github.com/moul/advanced-ssh-config/vendor/github.com/flynn/go-shlex"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdProxy(c *cli.Context) {
	if len(c.Args()) < 1 {
		Logger.Fatalf("assh: \"proxy\" requires 1 argument. See 'assh proxy --help'.")
	}

	conf, err := config.Open()
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	// FIXME: handle complete host with json

	host, err := computeHost(c.Args()[0], c.Int("port"), conf)
	if err != nil {
		Logger.Fatalf("Cannot get host '%s': %v", c.Args()[0], err)
	}

	err = conf.SaveSshConfig()
	if err != nil {
		Logger.Fatalf("Cannot save SSH config file: %v", err)
	}

	err = proxy(host, conf)
	if err != nil {
		Logger.Fatalf("Proxy error: %v", err)
	}
}

func computeHost(dest string, portOverride int, conf *config.Config) (*config.Host, error) {
	host := conf.GetHostSafe(dest)

	if portOverride > 0 {
		host.Port = uint(portOverride)
	}

	return host, nil
}

func prepareHostControlPath(host, gateway *config.Host) error {
	controlPathDir := path.Dir(os.ExpandEnv(strings.Replace(host.ControlPath, "~", "$HOME", -1)))
	gatewayControlPath := path.Join(controlPathDir, gateway.Name())
	return os.MkdirAll(gatewayControlPath, 0700)
}

func proxy(host *config.Host, conf *config.Config) error {
	if len(host.Gateways) > 0 {
		Logger.Debugf("Trying gateways: %s", host.Gateways)
		for _, gateway := range host.Gateways {
			if gateway == "direct" {
				err := proxyDirect(host)
				if err != nil {
					Logger.Errorf("Failed to use 'direct' connection")
				}
			} else {
				gatewayHost := conf.GetGatewaySafe(gateway)

				err := prepareHostControlPath(host, gatewayHost)
				if err != nil {
					return err
				}

				if host.ProxyCommand == "" {
					host.ProxyCommand = "nc %h %p"
				}
				command := "ssh %name -- " + commandApplyHost(host.ProxyCommand, host)

				Logger.Debugf("Using gateway '%s': %s", gateway, command)
				err = proxyCommand(gatewayHost, command)
				if err == nil {
					return nil
				}
				Logger.Errorf("Cannot use gateway '%s': %v", gateway, err)
			}
		}
		return fmt.Errorf("No such available gateway")
	}

	Logger.Debugf("Connecting without gateway")
	return proxyDirect(host)
}

func commandApplyHost(command string, host *config.Host) string {
	command = strings.Replace(command, "%name", host.Name(), -1)
	command = strings.Replace(command, "%h", host.HostName, -1)
	command = strings.Replace(command, "%p", fmt.Sprintf("%d", host.Port), -1)
	return command
}

func proxyDirect(host *config.Host) error {
	if host.ProxyCommand != "" {
		return proxyCommand(host, host.ProxyCommand)
	}
	return proxyGo(host)
}

func proxyCommand(host *config.Host, command string) error {
	command = commandApplyHost(command, host)
	args, err := shlex.Split(command)
	Logger.Debugf("ProxyCommand: %s", command)
	if err != nil {
		return err
	}
	spawn := exec.Command(args[0], args[1:]...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	return spawn.Run()
}

func proxyGo(host *config.Host) error {
	if host.HostName == "" {
		host.HostName = host.Name()
	}

	if len(host.ResolveNameservers) > 0 {
		Logger.Debugf("Resolving host: '%s' using nameservers %s", host.HostName, host.ResolveNameservers)
		// FIXME: resolve using custom dns server
		results, err := net.LookupAddr(host.HostName)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			host.HostName = results[0]
		}
		Logger.Debugf("Resolved host is: %s", host.HostName)
	}
	if host.ResolveCommand != "" {
		command := commandApplyHost(host.ResolveCommand, host)
		Logger.Debugf("Resolving host: '%s' using command: '%s'", host.HostName, command)

		args, err := shlex.Split(command)
		if err != nil {
			return err
		}

		out, err := exec.Command(args[0], args[1:]...).Output()
		if err != nil {
			return err
		}

		host.HostName = strings.TrimSpace(fmt.Sprintf("%s", out))
		Logger.Debugf("Resolved host is: %s", host.HostName)
	}

	Logger.Debugf("Connecting to %s:%d", host.HostName, host.Port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host.HostName, host.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	Logger.Debugf("Connected to %s:%d", host.HostName, host.Port)

	// Create Stdio pipes
	c1 := readAndWrite(conn, os.Stdout)
	c2 := readAndWrite(os.Stdin, conn)

	select {
	case err = <-c1:
	case err = <-c2:
	}
	if err != nil {
		return err
	}

	return nil
}

func readAndWrite(r io.Reader, w io.Writer) <-chan error {
	// Fixme: add an error channel
	buf := make([]byte, 1024)
	c := make(chan error)

	go func() {
		for {
			// Read
			n, err := r.Read(buf)
			if err != nil {
				if err != io.EOF {
					c <- err
				}
				break
			}

			// Write
			_, err = w.Write(buf[0:n])
			if err != nil {
				c <- err
			}
		}
		c <- nil
	}()
	return c
}

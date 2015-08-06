package commands

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/moul/advanced-ssh-config/vendor/github.com/Sirupsen/logrus"
	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"
	shlex "github.com/moul/advanced-ssh-config/vendor/github.com/flynn/go-shlex"

	"github.com/moul/advanced-ssh-config/pkg/config"
)

func cmdProxy(c *cli.Context) {
	if len(c.Args()) < 1 {
		logrus.Fatalf("assh: \"proxy\" requires 1 argument. See 'assh proxy --help'.")
	}

	host, err := computeHost(c.Args()[0], c.Int("port"))
	if err != nil {
		logrus.Fatalf("Cannot get host '%s': %v", c.Args()[0], err)
	}

	err = proxy(host)
	if err != nil {
		logrus.Fatalf("Proxy error: %v", err)
	}
}

func computeHost(dest string, portOverride int) (*config.Host, error) {
	conf, err := config.Open()
	if err != nil {
		return nil, err
	}

	host := conf.GetHostSafe(dest)
	if portOverride > 0 {
		host.Port = uint(portOverride)
	}

	return host, nil
}

func proxy(host *config.Host) error {
	if len(host.Gateways) > 0 {
		for _, gateway := range host.Gateways {
			gatewayHost, err := computeHost(gateway, 0)
			if err != nil {
				logrus.Fatalf("Cannot get host '%s': %v", gateway, err)
			}

			command := fmt.Sprintf("ssh {host} {port} nc -v -w 180 -G 5 %s %d", host.Host, host.Port)

			logrus.Debugf("Using gateway '%s': %s", gateway, command)
			err = proxyCommand(gatewayHost, command)
			if err != nil {
				logrus.Errorf("Cannot use gateway '%s': %v", gateway, err)
			}
			if err == nil {
				return nil
			}
		}
		return fmt.Errorf("No such available gateway")
	}
	// FIXME: proxyCommand(host, "nc -v -w 180 -G 5 {host} {port}")
	return proxyGo(host)
}

func proxyCommand(host *config.Host, command string) error {
	command = strings.Replace(command, "{host}", host.Host, -1)
	command = strings.Replace(command, "{port}", fmt.Sprintf("%d", host.Port), -1)
	args, err := shlex.Split(command)
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
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host.Host, host.Port))
	if err != nil {
		return err
	}

	defer conn.Close()

	logrus.Debugf("Connected to %s:%d\n", host.Host, host.Port)

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

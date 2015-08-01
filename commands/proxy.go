package commands

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/moul/advanced-ssh-config/vendor/github.com/Sirupsen/logrus"
	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/config"
)

func cmdProxy(c *cli.Context) {
	if len(c.Args()) < 1 {
		logrus.Fatalf("assh: \"proxy\" requires 1 argument. See 'assh proxy --help'.")
	}

	host, port, err := configGetHostPort(c.Args()[0], c.Int("port"))
	if err != nil {
		logrus.Fatalf("Cannot get host '%s': %v", c.Args()[0], err)
	}

	err = proxy(host, port)
	if err != nil {
		logrus.Fatalf("Proxy error: %v", err)
	}
}

func configGetHostPort(dest string, portFlag int) (string, uint, error) {
	conf, err := config.Open()
	if err != nil {
		return "", 0, err
	}

	// Get host configuration
	host := conf.GetHostSafe(dest)

	// Dial
	var port uint
	if portFlag > 0 {
		port = uint(portFlag)
	} else {
		port = host.Port
	}

	return host.Host, port, nil
}

func proxy(host string, port uint) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	defer conn.Close()

	logrus.Debugf("Connected to %s:%d\n", host, port)

	// Create Stdio pipes
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			logrus.Fatalf("Stdin pipe error: %v", err)
		}
	}()
	go func() {
		_, err := io.Copy(os.Stderr, conn)
		if err != nil {
			logrus.Fatalf("Stdout pipe error: %v", err)
		}
	}()
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		return err
	}

	return nil
}

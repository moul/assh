package commands

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/config"
)

func cmdProxy(c *cli.Context) {
	if len(c.Args()) < 1 {
		os.Exit(1)
	}

	dest := c.Args()[0]

	conf, err := config.Open()
	if err != nil {
		panic(err)
	}

	// Get host configuration
	host, err := conf.GetHost(dest)
	if err != nil {
		panic(err)
	}

	// Dial
	var port uint
	if c.Int("port") > 0 {
		port = uint(c.Int("port"))
	} else {
		port = host.Port
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host.Host, port))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Fprintf(os.Stderr, "Connected to %s:%d\n", host.Host, host.Port)

	// Create Stdio pipes
	go io.Copy(conn, os.Stdin)
	go io.Copy(os.Stdout, conn)
	_, err = io.Copy(os.Stderr, conn)
	if err != nil {
		panic(err)
	}

	return
}

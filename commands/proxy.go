package commands

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codegangsta/cli"

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

	if host, ok := conf.Hosts[dest]; ok {
		// Dial
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host.Host, host.Port))
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

	os.Exit(1)
}

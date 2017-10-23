package commands

import (
	"fmt"
	"net"
	"time"

	"github.com/urfave/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdPing(c *cli.Context) error {
	if len(c.Args()) < 1 {
		Logger.Fatalf("assh: \"ping\" requires exactly 1 argument. See 'assh ping --help'.")
	}

	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}
	if err = conf.LoadKnownHosts(); err != nil {
		Logger.Debugf("Failed to load assh known_hosts: %v", err)
	}
	target := c.Args()[0]
	host, err := computeHost(target, c.Int("port"), conf)
	if err != nil {
		Logger.Fatalf("Cannot get host '%s': %v", target, err)
	}

	if len(host.Gateways) > 0 {
		Logger.Fatalf("assh \"ping\" is not working with gateways (yet).")
	}
	if host.ProxyCommand != "" {
		Logger.Fatalf("assh \"ping\" is not working with custom ProxyCommand (yet).")
	}

	portName := "ssh" // fixme: resolve port name
	proto := "tcp"
	fmt.Printf("PING %s (%s) PORT %s (%s) PROTO %s\n", target, host.HostName, host.Port, portName, proto)
	dest := fmt.Sprintf("%s:%s", host.HostName, host.Port)
	for seq := 0; ; seq++ {
		start := time.Now()
		conn, err := net.DialTimeout(proto, dest, time.Second)
		duration := time.Now().Sub(start)
		if err == nil {
			defer conn.Close()
		}
		if err == nil {
			fmt.Printf("Connected to %s: seq=%d time=%v protocol=%s port=%s\n", host.HostName, seq, duration, proto, host.Port)
		} else {
			// FIXME: switch on error type
			fmt.Printf("Request timeout for seq %d (%v)\n", seq, err)
		}
		time.Sleep(time.Second)
	}

	return nil
}

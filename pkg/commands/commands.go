package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/moul/advanced-ssh-config/pkg/config"

	// . "github.com/moul/advanced-ssh-config/pkg/logger"
)

var SSHFlags = []cli.Flag{}

func init() {
	config.ASSHBinary = os.Args[0]

	// Populate SSHFlags
	boolFlags := []string{"1", "2", "4", "6", "A", "a", "C", "f", "G", "g", "K", "k", "M", "N", "n", "q", "s", "T", "t", "V", "v", "X", "x", "Y", "y"}
	stringFlags := []string{"b", "c", "D", "E", "e", "F", "I", "i", "L", "l", "m", "O", "o", "p", "Q", "R", "S", "W", "w"}
	for _, flag := range boolFlags {
		SSHFlags = append(SSHFlags, cli.BoolFlag{
			Name: flag,
		})
	}
	for _, flag := range stringFlags {
		SSHFlags = append(SSHFlags, cli.StringFlag{
			Name:  flag,
			Value: "",
		})
	}
	fmt.Println(SSHFlags)
}

// Commands is the list of cli commands
var Commands = []cli.Command{
	{
		Name:        "proxy",
		Usage:       "Connect to host SSH socket, used by ProxyCommand",
		Description: "Argument is a host.",
		Action:      cmdProxy,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port, p",
				Value: 0,
				Usage: "SSH destination port",
			},
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Only show how assh would connect but don't actually do it",
			},
		},
	},
	/*
		{
			Name:        "info",
			Usage:       "Print the connection config for host",
			Description: "Argument is a host.",
			Action:      cmdInfo,
		},
	*/
	/*
		{
			Name:        "init",
			Usage:       "Configure SSH to use assh",
			Description: "Build a .ssh/config.advanced file based on .ssh/config and update .ssh/config to use assh as ProxyCommand.",
			Action:      cmdInit,
		},
	*/
	{
		Name:   "build",
		Usage:  "Build .ssh/config",
		Action: cmdBuild,
	},
	/*
		{
			Name:        "etc-hosts",
			Usage:       "Generate a /etc/hosts file with assh hosts",
			Description: "Build a .ssh/config.advanced file based on .ssh/config and update .ssh/config to use assh as ProxyCommand.",
			Action:      cmdEtcHosts,
		},
	*/
	{
		Name:   "info",
		Usage:  "Display system-wide information",
		Action: cmdInfo,
	},
	{
		Name:   "wrapper",
		Usage:  "Initialize assh, then run SSH",
		Action: cmdWrapper,
		Flags:  SSHFlags,
	},
}

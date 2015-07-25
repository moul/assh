package commands

import "github.com/codegangsta/cli"

var Commands = []cli.Command{
	{
		Name:        "proxy",
		Usage:       "Connect to host SSH socket, used by ProxyCommand",
		Description: "Argument is a host.",
		Action:      cmdProxy,
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
	/*
		{
			Name:   "build",
			Usage:  "Build .ssh/config",
			Action: cmdBuild,
		},
	*/
	/*
		{
			Name:        "etc-hosts",
			Usage:       "Generate a /etc/hosts file with assh hosts",
			Description: "Build a .ssh/config.advanced file based on .ssh/config and update .ssh/config to use assh as ProxyCommand.",
			Action:      cmdEtcHosts,
		},
	*/
	{
		Name:   "stats",
		Usage:  "Print statistics",
		Action: cmdStats,
	},
}

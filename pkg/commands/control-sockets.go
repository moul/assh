package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/moul/advanced-ssh-config/pkg/control-sockets"
	// . "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdCsList(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		panic(err)
	}

	controlPath := conf.Defaults.ControlPath

	activeSockets, err := controlsockets.LookupControlPathDir(controlPath)
	if err != nil {
		panic(err)
	}

	if len(activeSockets) == 0 {
		fmt.Println("No active control sockets.")
		return nil
	}

	fmt.Printf("%d active control sockets in %q:\n\n", len(activeSockets), controlPath)
	for _, socket := range activeSockets {
		fmt.Printf("- %s\n", socket)
	}

	return nil
}

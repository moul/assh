package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/docker/go-units"

	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/moul/advanced-ssh-config/pkg/control-sockets"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
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
	now := time.Now().UTC()
	for _, socket := range activeSockets {
		createdAt, err := socket.CreatedAt()
		if err != nil {
			Logger.Warnf("%v", err)
		}

		fmt.Printf("- %s (%v)\n", socket.RelativePath(), units.HumanDuration(now.Sub(createdAt)))
	}

	return nil
}

func cmdCsFlush(c *cli.Context) error {
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
	}

	success := 0
	for _, socket := range activeSockets {
		if err := os.Remove(socket.Path()); err != nil {
			Logger.Warnf("Failed to close control socket %q: %v", socket.Path(), err)
		} else {
			success++
		}
	}

	if success > 0 {
		fmt.Printf("Closed %d control sockets.", success)
	}

	return nil
}

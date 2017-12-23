package commands

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	units "github.com/docker/go-units"
	"github.com/urfave/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/moul/advanced-ssh-config/pkg/control-sockets"
	"github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdCsList(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		logger.Logger.Errorf("%v", err)
		os.Exit(-1)
	}

	controlPath := conf.Defaults.ControlPath
	if controlPath == "" {
		logger.Logger.Errorf("Missing ControlPath in the configuration; Sockets features are disabled.")
		return nil
	}

	activeSockets, err := controlsockets.LookupControlPathDir(controlPath)
	if err != nil {
		logger.Logger.Errorf("%v", err)
		os.Exit(-1)
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
			logger.Logger.Warnf("%v", err)
		}

		fmt.Printf("- %s (%v)\n", socket.RelativePath(), units.HumanDuration(now.Sub(createdAt)))
	}

	return nil
}

func cmdCsMaster(c *cli.Context) error {
	if len(c.Args()) < 1 {
		logger.Logger.Fatalf("assh: \"sockets master\" requires 1 argument. See 'assh sockets master --help'.")
	}

	for _, target := range c.Args() {
		logger.Logger.Debugf("Opening master control socket for %q", target)

		cmd := exec.Command("ssh", target, "-M", "-N", "-f") // #nosec
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}

func cmdCsFlush(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		logger.Logger.Errorf("%v", err)
		os.Exit(-1)
	}

	controlPath := conf.Defaults.ControlPath
	if controlPath == "" {
		logger.Logger.Errorf("Missing ControlPath in the configuration; Sockets features are disabled.")
		return nil
	}

	activeSockets, err := controlsockets.LookupControlPathDir(controlPath)
	if err != nil {
		logger.Logger.Errorf("%v", err)
		os.Exit(-1)
	}

	if len(activeSockets) == 0 {
		fmt.Println("No active control sockets.")
	}

	success := 0
	for _, socket := range activeSockets {
		if err := os.Remove(socket.Path()); err != nil {
			logger.Logger.Warnf("Failed to close control socket %q: %v", socket.Path(), err)
		} else {
			success++
		}
	}

	if success > 0 {
		fmt.Printf("Closed %d control sockets.\n", success)
	}

	return nil
}

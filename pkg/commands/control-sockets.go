package commands

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	units "github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"moul.io/assh/pkg/config"
	"moul.io/assh/pkg/controlsockets"
)

func cmdCsList(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to open config")
	}

	controlPath := conf.Defaults.ControlPath
	if controlPath == "" {
		return errors.New("missing ControlPath in the configuration; Sockets features are disabled")
	}

	activeSockets, err := controlsockets.LookupControlPathDir(controlPath)
	if err != nil {
		return errors.Wrap(err, "failed to lookup control path")
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
			logger().Warn("failed to retrieve socket creation date", zap.Error(err))
		}

		fmt.Printf("- %s (%v)\n", socket.RelativePath(), units.HumanDuration(now.Sub(createdAt)))
	}

	return nil
}

func cmdCsMaster(c *cli.Context) error {
	if len(c.Args()) < 1 {
		return errors.New("assh: \"sockets master\" requires 1 argument. See 'assh sockets master --help'")
	}

	for _, target := range c.Args() {
		logger().Debug("Opening master control socket", zap.String("host", target))

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
		return errors.Wrap(err, "failed to open config")
	}

	controlPath := conf.Defaults.ControlPath
	if controlPath == "" {
		return errors.New("missing ControlPath in the configuration; Sockets features are disabled")
	}

	activeSockets, err := controlsockets.LookupControlPathDir(controlPath)
	if err != nil {
		return errors.Wrap(err, "failed to lookup control path")
	}

	if len(activeSockets) == 0 {
		fmt.Println("No active control sockets.")
		return nil
	}

	success := 0
	for _, socket := range activeSockets {
		if err := os.Remove(socket.Path()); err != nil {
			logger().Warn("Failed to close control socket", zap.String("path", socket.Path()), zap.Error(err))
		} else {
			success++
		}
	}

	if success > 0 {
		fmt.Printf("Closed %d control sockets.\n", success)
	}

	return nil
}

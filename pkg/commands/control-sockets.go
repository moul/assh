package commands

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	units "github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"moul.io/assh/pkg/config"
	"moul.io/assh/pkg/controlsockets"
)

var socketsCommand = &cobra.Command{
	Use:   "sockets",
	Short: "Manage control sockets",
}

var listSocketsCommand = &cobra.Command{
	Use:   "list",
	Short: "List active control sockets",
	RunE:  runListSocketsCommand,
}

var flushSocketsCommand = &cobra.Command{
	Use:   "flush",
	Short: "Close control sockets",
	RunE:  runFlushSocketsCommand,
}

var masterSocketCommand = &cobra.Command{
	Use:   "master",
	Short: "Open a master control socket",
	RunE:  runMasterSocketCommand,
}

func init() {
	socketsCommand.AddCommand(listSocketsCommand)
	socketsCommand.AddCommand(flushSocketsCommand)
	socketsCommand.AddCommand(masterSocketCommand)
}

func runListSocketsCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.Open(viper.GetString("config"))
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

func runMasterSocketCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("assh: \"sockets master\" requires 1 argument. See 'assh sockets master --help'")
	}

	for _, target := range args {
		logger().Debug("Opening master control socket", zap.String("host", target))

		cmd := exec.Command("ssh", target, "-M", "-N", "-f") // #nosec
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}

func runFlushSocketsCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.Open(viper.GetString("config"))
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

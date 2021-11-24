package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"moul.io/assh/v2/pkg/config"
	"moul.io/assh/v2/pkg/controlsockets"
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

// nolint:gochecknoinits
func init() {
	socketsCommand.AddCommand(listSocketsCommand)
	socketsCommand.AddCommand(flushSocketsCommand)
	socketsCommand.AddCommand(masterSocketCommand)
}

func runListSocketsCommand(_ *cobra.Command, _ []string) error {
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

func runMasterSocketCommand(_ *cobra.Command, args []string) error {
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

func runFlushSocketsCommand(_ *cobra.Command, _ []string) error {
	success := 0

	if processes, err := process.Processes(); err != nil {
		return err
	} else {
		for _, ps := range processes {
			if cmdline, err := ps.CmdlineSlice(); err == nil && len(cmdline) > 0 && path.Base(cmdline[0]) == "assh" && cmdline[1] == "connect" {
				cmd := exec.Command("ssh", "-O", "exit", cmdline[len(cmdline)-1]) // #nosec
				if err := cmd.Run(); err != nil {
					success++
				}
			}
		}
	}

	if success > 0 {
		fmt.Printf("Closed %d control sockets.\n", success)
	}

	return nil
}

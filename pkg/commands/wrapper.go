package commands

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"moul.io/assh/pkg/config"
)

var wrapperCommand = &cobra.Command{
	Use:    "wrapper",
	Short:  "Initialize assh, then run ssh/scp/rsync...",
	Hidden: true,
}

var sshWrapperCommand = &cobra.Command{
	Use:   "ssh",
	Short: "Wrap ssh",
	RunE:  runSSHWrapperCommand,
}

func init() {
	sshWrapperCommand.Flags().AddFlagSet(config.SSHFlags)
	wrapperCommand.AddCommand(sshWrapperCommand)
}

func runSSHWrapperCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing <target> argument. See usage with 'assh wrapper %s -h'", cmd.Name())
	}

	// prepare variables
	target := args[0]
	command := args[1:]
	options := []string{}
	for _, flag := range config.SSHBoolFlags {
		if viper.GetBool(flag) {
			options = append(options, fmt.Sprintf("-%s", flag))
		}
	}
	for _, flag := range config.SSHStringFlags {
		for _, val := range viper.GetStringSlice(flag) {
			if (flag == "o" || flag == "O") && val == "false" {
				logger().Debug(
					"Skip invalid option:",
					zap.String("flag", flag),
					zap.String("val", val),
				)
				continue
			}
			options = append(options, fmt.Sprintf("-%s", flag))
			options = append(options, val)
		}
	}
	sshArgs := []string{cmd.Name()}
	sshArgs = append(sshArgs, options...)
	sshArgs = append(sshArgs, target)
	sshArgs = append(sshArgs, command...)
	bin, err := exec.LookPath(cmd.Name())
	if err != nil {
		return errors.Wrapf(err, "failed to lookup %q", cmd.Name())
	}

	logger().Debug(
		"Wrapper called",
		zap.String("bin", bin),
		zap.String("target", target),
		zap.Any("command", command),
		zap.Any("options", options),
		zap.Any("sshArgs", sshArgs),
	)

	// check if config is up-to-date
	conf, err := config.Open(viper.GetString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to open config")
	}

	if err = conf.LoadKnownHosts(); err != nil {
		logger().Debug("Failed to load assh known_hosts", zap.Error(err))
	}

	// check if .ssh/config is outdated
	isOutdated, err := conf.IsConfigOutdated(target)
	if err != nil {
		logger().Error("failed to check if config is outdated", zap.Error(err))
	}
	if isOutdated {
		logger().Debug(
			"The configuration file is outdated, rebuilding it before calling command",
			zap.String("command", cmd.Name()),
		)
		if err = conf.SaveSSHConfig(); err != nil {
			logger().Error("failed to save ssh config file", zap.Error(err))
		}
	}

	// Execute Binary
	return syscall.Exec(bin, sshArgs, os.Environ()) // #nosec
}

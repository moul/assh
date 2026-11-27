package commands

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"moul.io/assh/v2/pkg/config"
)

var wrapperCommand = &cobra.Command{
	Use:    "wrapper",
	Short:  "Initialize assh, then run ssh/scp/rsync...",
	Hidden: true,
}

var sshWrapperCommand = &cobra.Command{
	Use:                "ssh [options] destination [command]",
	Short:              "Wrap ssh",
	DisableFlagParsing: true,
	RunE:               runSSHWrapperCommand,
}

// nolint:gochecknoinits
func init() {
	wrapperCommand.AddCommand(sshWrapperCommand)
}

func runSSHWrapperCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		return cmd.Help()
	}
	// Preserve the separator used by the documented shell alias.
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) == 0 {
		return fmt.Errorf("missing <target> argument. See usage with 'assh wrapper %s -h'", cmd.Name())
	}

	bin, err := exec.LookPath(cmd.Name())
	if err != nil {
		return errors.Wrapf(err, "failed to lookup %q", cmd.Name())
	}
	target := sshDestination(args)
	logger().Debug("Wrapper called", zap.String("bin", bin), zap.String("target", target), zap.Strings("sshArgs", args))

	if target != "" {
		if err := refreshWrapperConfig(target); err != nil {
			return err
		}
	}

	sshArgs := append([]string{cmd.Name()}, args...)
	return syscall.Exec(bin, sshArgs, os.Environ()) // #nosec
}

// sshDestination only identifies which host may require a config refresh.
// OpenSSH receives the original arguments and handles validation.
func sshDestination(args []string) string {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			if i+1 < len(args) {
				return args[i+1]
			}
			return ""
		}
		if len(arg) < 2 || arg[0] != '-' {
			return arg
		}
		if strings.HasPrefix(arg, "--") {
			return ""
		}

		for j := 1; j < len(arg); j++ {
			if !slices.Contains(config.SSHStringFlags, string(arg[j])) {
				continue
			}
			if j == len(arg)-1 {
				// An option with no attached value consumes the next argument.
				i++
			}
			break
		}
	}
	return ""
}

func refreshWrapperConfig(target string) error {
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
			zap.String("target", target),
		)
		if err = rewriteSSHConfig(conf); err != nil {
			logger().Error("failed to save ssh config file", zap.Error(err))
		}
	}
	return nil
}

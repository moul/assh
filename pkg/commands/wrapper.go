package commands

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"moul.io/assh/pkg/config"
)

func cmdWrapper(c *cli.Context) error {
	if len(c.Args()) < 1 {
		return fmt.Errorf("missing <target> argument. See usage with 'assh wrapper %s -h'", c.Command.Name)
	}

	// prepare variables
	target := c.Args()[0]
	command := c.Args()[1:]
	options := []string{}
	for _, flag := range config.SSHBoolFlags {
		if c.Bool(flag) {
			options = append(options, fmt.Sprintf("-%s", flag))
		}
	}
	for _, flag := range config.SSHStringFlags {
		for _, val := range c.StringSlice(flag) {
			options = append(options, fmt.Sprintf("-%s", flag))
			options = append(options, val)
		}
	}
	args := []string{c.Command.Name}
	args = append(args, options...)
	args = append(args, target)
	args = append(args, command...)
	bin, err := exec.LookPath(c.Command.Name)
	if err != nil {
		return errors.Wrapf(err, "failed to lookup %q", c.Command.Name)
	}

	logger().Debug(
		"Wrapper called",
		zap.String("bin", bin),
		zap.String("target", target),
		zap.Any("command", command),
		zap.Any("options", options),
		zap.Any("args", args),
	)

	// check if config is up-to-date
	conf, err := config.Open(c.GlobalString("config"))
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
			zap.String("command", c.Command.Name),
		)
		if err = conf.SaveSSHConfig(); err != nil {
			logger().Error("failed to save ssh config file", zap.Error(err))
		}
	}

	// Execute Binary
	return syscall.Exec(bin, args, os.Environ()) // #nosec
}

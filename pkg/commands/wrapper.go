package commands

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdWrapper(c *cli.Context) {
	if len(c.Args()) < 1 {
		Logger.Fatalf("Missing <target> argument. See usage with 'assh wrapper %s -h'.", c.Command.Name)
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
		if val := c.String(flag); val != "" {
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
		Logger.Fatalf("Cannot find %q in $PATH", c.Command.Name)
	}

	Logger.Debugf("Wrapper called with bin=%v target=%v command=%v options=%v, args=%v", bin, target, command, options, args)

	// check if config is up-to-date
	conf, err := config.Open()
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	if err = conf.LoadKnownHosts(); err != nil {
		Logger.Debugf("Failed to load assh known_hosts: %v", err)
	}

	if conf.NeedsARebuildForTarget(target) {
		Logger.Debugf("The configuration file is outdated, rebuilding it before calling %s", c.Command.Name)
		conf.SaveNewKnownHost(target)
		if err = conf.SaveSshConfig(); err != nil {
			Logger.Error(err)
		}
	}

	// Execute Binary
	syscall.Exec(bin, args, os.Environ())
}

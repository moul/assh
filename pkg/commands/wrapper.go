package commands

import (
	"fmt"
	"os"
	"syscall"

	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdWrapper(c *cli.Context) {
	if len(c.Args()) < 1 {
		Logger.Fatalf("Missing <target> argument. See usage with 'assh wrapper -h'.")
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
	args := append(options, target)
	args = append(args, command...)
	bin := "/usr/bin/ssh"
	Logger.Debugf("Wrapper called with target=%v command=%v ssh-options=%v", target, command, options)

	// check if config is up-to-date
	conf, err := config.Open()
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	if conf.NeedsARebuildForTarget(target) {
		Logger.Debugf("The configuration file is outdated, rebuilding it before calling ssh")
		//host := conf.GetHostSafe(target)
		//fmt.Println(host)
		//fmt.Println(host.name)
		//conf.WriteSshConfigTo(os.Stdout)
	}

	// Execute SSH
	syscall.Exec(bin, args, os.Environ())
}

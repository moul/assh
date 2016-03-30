package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdWrapper(c *cli.Context) {
	if len(c.Args()) < 1 {
		Logger.Fatalf("Missing <target> argument. See usage with 'assh wrapper -h'.")
	}

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

	Logger.Debugf("Wrapper called with target=%v command=%v ssh-options=%v", target, command, options)

	conf, err := config.Open()
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	fmt.Println(conf.NeedsARebuildForTarget(target))
	//host := conf.GetHostSafe(target)
	//fmt.Println(host)
	//fmt.Println(host.name)

	//conf.WriteSshConfigTo(os.Stdout)
}

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
	// FIXME: populate options

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

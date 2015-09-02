package commands

import (
	"os"

	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdBuild(c *cli.Context) {
	conf, err := config.Open()
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	conf.WriteSshConfigTo(os.Stdout)
}

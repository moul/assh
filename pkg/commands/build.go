package commands

import (
	"os"

	"github.com/moul/advanced-ssh-config/vendor/github.com/Sirupsen/logrus"
	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
)

func cmdBuild(c *cli.Context) {
	conf, err := config.Open()
	if err != nil {
		logrus.Fatalf("Cannot open configuration file: %v", err)
	}

	conf.WriteSshConfigTo(os.Stdout)
}

package commands

import (
	"fmt"

	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	// . "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdStats(c *cli.Context) {
	conf, err := config.Open()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d hosts\n", len(conf.Hosts))
	fmt.Printf("%d templates\n", len(conf.Templates))
	fmt.Printf("%d included files\n", len(conf.IncludedFiles()))
}

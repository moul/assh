package commands

import (
	"fmt"

	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/config"
)

func cmdStats(c *cli.Context) {
	conf, err := config.Open()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d hosts\n", len(conf.Hosts))
}

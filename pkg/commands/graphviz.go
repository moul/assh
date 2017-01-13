package commands

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/moul/advanced-ssh-config/pkg/config/graphviz"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdGraphviz(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		Logger.Fatalf("Cannot load configuration: %v", err)
		return nil
	}

	settings := configviz.GraphSettings{
		ShowIsolatedHosts: c.Bool("show-isolated-hosts"),
	}
	graph, err := configviz.Graph(conf, &settings)
	if err != nil {
		Logger.Fatalf("failed to build graph: %v", err)
		return nil
	}

	fmt.Println(graph)
	return nil
}

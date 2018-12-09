package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"moul.io/assh/pkg/config"
	"moul.io/assh/pkg/config/graphviz"
)

func cmdGraphviz(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	settings := graphviz.GraphSettings{
		ShowIsolatedHosts: c.Bool("show-isolated-hosts"),
		NoResolveWildcard: c.Bool("no-resolve-wildcard"),
		NoInherits:        c.Bool("no-inheritance-links"),
	}
	graph, err := graphviz.Graph(conf, &settings)
	if err != nil {
		return errors.Wrap(err, "failed to build graph")
	}

	fmt.Println(graph)
	return nil
}

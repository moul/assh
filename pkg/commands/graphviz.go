package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"moul.io/assh/pkg/config"
	"moul.io/assh/pkg/config/graphviz"
)

var graphvizConfigCommand = &cobra.Command{
	Use:   "graphviz",
	Short: "Generate a Graphviz graph of the hosts",
	RunE:  runGraphvizConfigCommand,
}

func init() {
	graphvizConfigCommand.Flags().BoolP("show-isolated-hosts", "", false, "Show isolated hosts")
	graphvizConfigCommand.Flags().BoolP("no-resolve-wildcard", "", false, "Do not resolve wildcards in Gateways")
	graphvizConfigCommand.Flags().BoolP("no-inheritance-links", "", false, "Do not show inheritance links")
	viper.BindPFlags(graphvizConfigCommand.Flags())
}

func runGraphvizConfigCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.Open(viper.GetString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	settings := graphviz.GraphSettings{
		ShowIsolatedHosts: viper.GetBool("show-isolated-hosts"),
		NoResolveWildcard: viper.GetBool("no-resolve-wildcard"),
		NoInherits:        viper.GetBool("no-inheritance-links"),
	}
	graph, err := graphviz.Graph(conf, &settings)
	if err != nil {
		return errors.Wrap(err, "failed to build graph")
	}

	fmt.Println(graph)
	return nil
}

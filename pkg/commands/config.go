package commands

import "github.com/spf13/cobra"

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Manage ssh and assh configuration",
}

// nolint:gochecknoinits
func init() {
	configCommand.AddCommand(buildConfigCommand)
	configCommand.AddCommand(buildJSONConfigCommand)
	configCommand.AddCommand(listConfigCommand)
	configCommand.AddCommand(graphvizConfigCommand)
	configCommand.AddCommand(searchConfigCommand)
}

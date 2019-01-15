package commands // import "moul.io/assh/pkg/commands"

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var genCompletionsCommand = &cobra.Command{
	Use:    "completion",
	Short:  "generate shell autocompletion scripts for assh",
	RunE:   runGenCompletionsCommand,
	Hidden: true,
}

func runGenCompletionsCommand(cmd *cobra.Command, args []string) error {
	if err := RootCmd.GenBashCompletionFile("contrib/completion/bash_autocomplete"); err != nil {
		return errors.Wrap(err, "failed to generate bash completion file")
	}
	if err := RootCmd.GenZshCompletionFile("contrib/completion/zsh_autocomplete"); err != nil {
		return errors.Wrap(err, "failed to generate zsh completion file")
	}
	return nil
}

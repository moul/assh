package main // import "moul.io/assh/contrib/completion/gen"

import (
	"log"

	"moul.io/assh/pkg/commands"
)

func main() {
	if err := commands.RootCmd.GenBashCompletionFile("../bash_autocomplete"); err != nil {
		log.Println("failed to generate bash completion file: ", err)
	}
	if err := commands.RootCmd.GenZshCompletionFile("../zsh_autocomplete"); err != nil {
		log.Println("failed to generate zsh completion file: ", err)
	}
}

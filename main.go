//go:generate sh -c "cd contrib/completion/gen && go run main.go"

package main

import (
	"fmt"
	"os"

	"moul.io/assh/v2/pkg/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

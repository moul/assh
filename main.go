//go:generate sh -c "cd contrib/completion/gen && go run main.go"

package main // import "moul.io/assh"

import (
	"fmt"
	"os"

	"moul.io/assh/pkg/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

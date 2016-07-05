package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/mgutz/ansi"
	"github.com/moul/advanced-ssh-config/pkg/config"
	// . "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdList(c *cli.Context) {
	conf, err := config.Open()
	if err != nil {
		panic(err)
	}

	// ansi coloring
	colorize := func(input string) string { return input }
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		colorize = ansi.ColorFunc("green+b+h")
	}

	for _, host := range conf.Hosts.SortedList() {
		host.ApplyDefaults(&conf.Defaults)
		fmt.Printf("    %s -> %s\n", colorize(host.Name()), host.Prototype())
		fmt.Println()
	}
}

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/mgutz/ansi"
	"github.com/moul/advanced-ssh-config/pkg/config"
	// . "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdList(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		panic(err)
	}

	// ansi coloring
	greenColorize := func(input string) string { return input }
	redColorize := func(input string) string { return input }
	yellowColorize := func(input string) string { return input }
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		greenColorize = ansi.ColorFunc("green+b+h")
		redColorize = ansi.ColorFunc("red")
		yellowColorize = ansi.ColorFunc("yellow")
	}

	fmt.Printf("Listing entries\n\n")

	for _, host := range conf.Hosts.SortedList() {
		options := host.Options()
		options.Remove("User")
		options.Remove("Port")
		host.ApplyDefaults(&conf.Defaults)
		fmt.Printf("    %s -> %s\n", greenColorize(host.Name()), host.Prototype())
		if len(options) > 0 {
			fmt.Printf("        %s %s\n", yellowColorize("[custom options]"), strings.Join(options.ToStringList(), " "))
		}
		fmt.Println()
	}

	generalOptions := conf.Defaults.Options()
	if len(generalOptions) > 0 {
		fmt.Println(greenColorize("    (*) General options:"))
		for _, opt := range conf.Defaults.Options() {
			fmt.Printf("        %s: %s\n", redColorize(opt.Name), opt.Value)
		}
		fmt.Println()
	}

	return nil
}

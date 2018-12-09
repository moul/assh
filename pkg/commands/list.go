package commands

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"moul.io/assh/pkg/config"
)

func cmdList(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	// ansi coloring
	greenColorize := func(input string) string { return input }
	redColorize := func(input string) string { return input }
	yellowColorize := func(input string) string { return input }
	cyanColorize := func(input string) string { return input }
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		greenColorize = ansi.ColorFunc("green+b+h")
		redColorize = ansi.ColorFunc("red")
		yellowColorize = ansi.ColorFunc("yellow")
		cyanColorize = ansi.ColorFunc("cyan")
	}

	fmt.Printf("Listing entries\n\n")

	if c.Bool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				return errors.Wrap(err, "failed to expand hosts")
			}
		}
	}

	generalOptions := conf.Defaults.Options()

	for _, host := range conf.Hosts.SortedList() {
		options := host.Options()
		options.Remove("User")
		options.Remove("Port")
		host.ApplyDefaults(&conf.Defaults)
		fmt.Printf("    %s -> %s\n", greenColorize(host.Name()), host.Prototype())

		for _, opt := range options {
			defaultValue := generalOptions.Get(opt.Name)
			switch {
			case defaultValue == "":
				fmt.Printf("        %s %s %s\n", yellowColorize(opt.Name), opt.Value, yellowColorize("[custom option]"))
			case defaultValue == opt.Value:
				fmt.Printf("        %s: %s\n", redColorize(opt.Name), opt.Value)
			default:
				fmt.Printf("        %s %s %s\n", cyanColorize(opt.Name), opt.Value, cyanColorize("[override]"))
			}
		}
		fmt.Println()
	}

	if len(generalOptions) > 0 {
		fmt.Println(greenColorize("    (*) General options:"))
		for _, opt := range conf.Defaults.Options() {
			fmt.Printf("        %s: %s\n", redColorize(opt.Name), opt.Value)
		}
		fmt.Println()
	}

	return nil
}

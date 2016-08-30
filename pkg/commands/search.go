package commands

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdSearch(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		Logger.Fatalf("Cannot load configuration: %v", err)
		return nil
	}

	needle := c.Args()[0]

	found := []*config.Host{}
	for _, host := range conf.Hosts.SortedList() {
		if host.Matches(needle) {
			found = append(found, host)
		}
	}

	if len(found) == 0 {
		fmt.Println("no results found.")
		return nil
	}

	fmt.Printf("Listing results for %s:\n", needle)
	for _, host := range found {
		fmt.Printf("    %s -> %s\n", host.Name(), host.Prototype())
	}

	return nil
}

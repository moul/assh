package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	// . "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdSearch(c *cli.Context) {
	conf, err := config.Open()
	if err != nil {
		panic(err)
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
		return
	}

	fmt.Printf("Listing results for %s:\n", needle)
	for _, host := range found {
		fmt.Printf("    %s -> %s\n", host.Name(), host.Prototype())
	}
}

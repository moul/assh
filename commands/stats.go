package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdStats(c *cli.Context) {
	fmt.Printf("stats: %v\n", c)
}

package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"moul.io/assh/v2/pkg/config"
)

var searchConfigCommand = &cobra.Command{
	Use:   "search",
	Short: "Search entries by given search text",
	RunE:  searchConfig,
}

func searchConfig(cmd *cobra.Command, args []string) error {
	conf, err := config.Open(viper.GetString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	if len(args) != 1 {
		return errors.New("assh config search requires 1 argument. See 'assh config search --help'")
	}

	needle := strings.ToLower(args[0]) // Make search case-insensitive

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

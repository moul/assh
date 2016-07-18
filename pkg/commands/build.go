package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdBuild(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	conf.WriteSSHConfigTo(os.Stdout)

	return nil
}

func cmdBuildJSON(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	s, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		Logger.Fatalf("JSON encoding error: %v", err)
	}
	fmt.Println(string(s))

	return nil
}

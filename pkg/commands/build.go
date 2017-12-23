package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdBuild(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		logger.Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	if c.Bool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				logger.Logger.Fatalf("Error while trying to expand hosts: %v", err)
			}
		}
	}

	if !c.Bool("ignore-known-hosts") {
		if conf.KnownHostsFileExists() == nil {
			if err := conf.LoadKnownHosts(); err != nil {
				logger.Logger.Errorf("Failed to load known-hosts file: %v", err)
			}
		}
	}

	if c.Bool("no-automatic-rewrite") {
		conf.DisableAutomaticRewrite()
	}
	return conf.WriteSSHConfigTo(os.Stdout)
}

func cmdBuildJSON(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		logger.Logger.Fatalf("Cannot open configuration file: %v", err)
	}

	if c.Bool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				logger.Logger.Fatalf("Error while trying to expand hosts: %v", err)
			}
		}
	}

	s, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		logger.Logger.Fatalf("JSON encoding error: %v", err)
	}
	fmt.Println(string(s))

	return nil
}

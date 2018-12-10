package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"moul.io/assh/pkg/config"
)

func cmdBuild(c *cli.Context) error {
	conf, err := config.Open(c.GlobalString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	if c.Bool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				return errors.Wrap(err, "failed to expand hosts")
			}
		}
	}

	if !c.Bool("ignore-known-hosts") {
		if conf.KnownHostsFileExists() == nil {
			if err := conf.LoadKnownHosts(); err != nil {
				return errors.Wrap(err, "failed to load known-hosts file")
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
		return errors.Wrap(err, "failed to open configuration file")
	}

	if c.Bool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				return errors.Wrap(err, "failed to expand hosts")
			}
		}
	}

	s, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	fmt.Println(string(s))
	return nil
}

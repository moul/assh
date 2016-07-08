package main

import (
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/commands"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
	"github.com/moul/advanced-ssh-config/pkg/version"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Manfred Touron"
	app.Email = "https://github.com/moul/advanced-ssh-config"
	app.Version = version.VERSION + " (" + version.GITCOMMIT + ")"
	app.Usage = "advanced ssh config"
	app.EnableBashCompletion = true
	app.BashComplete = BashComplete

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			EnvVar: "ASSH_CONFIG",
			Value:  "~/.ssh/assh.yml",
			Usage:  "Location of config file",
		},
		cli.BoolFlag{
			Name:   "debug, D",
			EnvVar: "ASSH_DEBUG",
			Usage:  "Enable debug mode",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "Enable verbose mode",
		},
	}

	app.Before = hookBefore

	app.Commands = commands.Commands

	app.Run(os.Args)
}

func BashComplete(c *cli.Context) {
	if len(c.Args()) == 0 {
		for _, option := range []string{"--debug", "--verbose", "--help", "--version"} {
			fmt.Println(option)
		}
		for _, command := range []string{"connect", "build", "info", "help"} {
			fmt.Println(command)
		}
	}
}

func hookBefore(c *cli.Context) error {
	if c.Bool("debug") {
		os.Setenv("ASSH_DEBUG", "1")
	}
	initLogging(c.Bool("debug"), c.Bool("verbose"))
	return nil
}

func initLogging(debug bool, verbose bool) {
	options := LoggerOptions{}

	options.InspectParent = true

	if debug {
		options.Level = logrus.DebugLevel
	} else if verbose {
		options.Level = logrus.InfoLevel
	} else {
		options.Level = logrus.WarnLevel
	}
	SetupLogging(options)
}

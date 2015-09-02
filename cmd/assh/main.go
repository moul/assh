package main

import (
	"os"
	"path"

	"github.com/moul/advanced-ssh-config/vendor/github.com/Sirupsen/logrus"
	"github.com/moul/advanced-ssh-config/vendor/github.com/codegangsta/cli"

	"github.com/moul/advanced-ssh-config/pkg/commands"
	. "github.com/moul/advanced-ssh-config/pkg/logger"
	"github.com/moul/advanced-ssh-config/pkg/version"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Manfred Touron"
	app.Email = "https://github.com/moul/advanced-ssh-config"
	app.Version = version.Version + " (" + version.GitCommit + ")"
	app.Usage = "advanced ssh config"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug mode",
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

func hookBefore(c *cli.Context) error {
	initLogging(c.Bool("debug"), c.Bool("verbose"))
	return nil
}

func initLogging(debug bool, verbose bool) {
	if debug {
		LoggerSetLevel(logrus.DebugLevel)
	} else if verbose {
		LoggerSetLevel(logrus.InfoLevel)
	} else {
		LoggerSetLevel(logrus.WarnLevel)
	}
}

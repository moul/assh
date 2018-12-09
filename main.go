package main // import "moul.io/assh"

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"moul.io/assh/pkg/commands"
	loggerpkg "moul.io/assh/pkg/logger"
	"moul.io/assh/pkg/version"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Manfred Touron"
	app.Email = "https://github.com/moul/assh"
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

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// BashComplete is used bu urfave/cli to allow dynamic bash completion.
func BashComplete(c *cli.Context) {
	if len(c.Args()) == 0 {
		for _, option := range []string{"--debug", "--verbose", "--help", "--version"} {
			fmt.Println(option)
		}
		for _, command := range []string{"connect", "config", "info", "sockets", "help"} {
			fmt.Println(command)
		}
	}
}

func hookBefore(c *cli.Context) error {
	if c.Bool("debug") {
		if err := os.Setenv("ASSH_DEBUG", "1"); err != nil {
			return err
		}
	}
	return initLogging(c.Bool("debug"), c.Bool("verbose"))
}

func initLogging(debug bool, verbose bool) error {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(loggerpkg.MustLogLevel(debug, verbose))
	if !debug {
		config.DisableStacktrace = true
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := config.Build()
	if err != nil {
		return errors.Wrap(err, "failed to initialize logger")
	}
	zap.ReplaceGlobals(l)
	return nil
}

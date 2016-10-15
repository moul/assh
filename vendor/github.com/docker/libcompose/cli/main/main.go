package main

import (
	"os"

	cliApp "github.com/docker/libcompose/cli/app"
	"github.com/docker/libcompose/cli/command"
	dockerApp "github.com/docker/libcompose/cli/docker/app"
	"github.com/docker/libcompose/version"
	"github.com/urfave/cli"
)

func main() {
	factory := &dockerApp.ProjectFactory{}

	app := cli.NewApp()
	app.Name = "libcompose-cli"
	app.Usage = "Command line interface for libcompose."
	app.Version = version.VERSION + " (" + version.GITCOMMIT + ")"
	app.Author = "Docker Compose Contributors"
	app.Email = "https://github.com/docker/libcompose"
	app.Before = cliApp.BeforeApp
	app.Flags = append(command.CommonFlags(), dockerApp.DockerClientFlags()...)
	app.Commands = []cli.Command{
		command.BuildCommand(factory),
		command.ConfigCommand(factory),
		command.CreateCommand(factory),
		command.EventsCommand(factory),
		command.DownCommand(factory),
		command.KillCommand(factory),
		command.LogsCommand(factory),
		command.PauseCommand(factory),
		command.PortCommand(factory),
		command.PsCommand(factory),
		command.PullCommand(factory),
		command.RestartCommand(factory),
		command.RmCommand(factory),
		command.RunCommand(factory),
		command.ScaleCommand(factory),
		command.StartCommand(factory),
		command.StopCommand(factory),
		command.UnpauseCommand(factory),
		command.UpCommand(factory),
		command.VersionCommand(factory),
	}

	app.Run(os.Args)

}

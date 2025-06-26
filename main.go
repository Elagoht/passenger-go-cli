package main

import (
	"os"
	"passenger-go-cli/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "passenger",
		Commands: []*cli.Command{
			cmd.ServerCommand(),
			cmd.StatusCommand(),
			cmd.LoginCommand(),
			cmd.RegisterCommand(),
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		cli.Exit(err.Error(), 1)
	}
}

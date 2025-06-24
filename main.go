package main

import (
	"fmt"
	"os"
	"passenger-go-cli/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "passenger",
		Commands: []*cli.Command{
			cmd.ServerCommand(),
			cmd.LoginCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

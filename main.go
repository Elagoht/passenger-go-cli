package main

import (
	"os"
	"passenger-go-cli/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "passenger-go",
		Commands: []*cli.Command{
			cmd.ServerCommand(),
			cmd.StatusCommand(),
			cmd.LoginCommand(),
			cmd.LogoutCommand(),
			cmd.RegisterCommand(),
			cmd.ValidateCommand(),
			cmd.ListCommand(),
			cmd.GetCommand(),
			cmd.PassphraseCommand(),
			cmd.ChangeMasterPassphraseCommand(),
			cmd.GenerateCommand(),
			cmd.AlternateCommand(),
			cmd.CreateCommand(),
			cmd.ExportCommand(),
			cmd.ImportCommand(),
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

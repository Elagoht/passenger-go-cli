package cmd

import "github.com/urfave/cli/v2"

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "show-all", "fetch-all"},
		Usage:   "Will list all accounts",
	}
}

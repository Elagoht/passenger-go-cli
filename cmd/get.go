package cmd

import "github.com/urfave/cli/v2"

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:    "get",
		Aliases: []string{"fetch", "show"},
		Usage:   "Will get the account by id",
	}
}

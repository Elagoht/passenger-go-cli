package cmd

import "github.com/urfave/cli/v2"

func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"remove", "rm", "del", "kaboom", "shred"},
		Usage:   "Will delete the account by id",
	}
}

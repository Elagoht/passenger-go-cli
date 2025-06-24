package cmd

import "github.com/urfave/cli/v2"

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"add", "new", "insert"},
		Usage:   "This will ask required fields, can pass stdin, prevents storing sensitive data in history",
	}
}

package cmd

import "github.com/urfave/cli/v2"

func UpdateCommand() *cli.Command {
	return &cli.Command{
		Name:    "update",
		Aliases: []string{"edit", "modify", "change"},
		Usage:   "This will ask required fields, can pass stdin, press ctrl-c to keep the current value",
	}
}

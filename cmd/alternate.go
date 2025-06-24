package cmd

import "github.com/urfave/cli/v2"

func AlternateCommand() *cli.Command {
	return &cli.Command{
		Name:    "alternate",
		Aliases: []string{"alt", "alternative", "manipulate", "shuffle"},
		Usage:   "Alternate characters with similar looking characters.",
	}
}

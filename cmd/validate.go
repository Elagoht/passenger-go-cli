package cmd

import "github.com/urfave/cli/v2"

func ValidateCommand() *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"verify", "check"},
		Usage:   "Validate the recovery key. Server needs to verify you have really backed up your recovery key.",
	}
}

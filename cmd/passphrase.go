package cmd

import "github.com/urfave/cli/v2"

func PassphraseCommand() *cli.Command {
	return &cli.Command{
		Name:    "passphrase",
		Aliases: []string{"pass", "passw", "password", "pw"},
		Usage:   "Will print the passphrase for the account",
	}
}

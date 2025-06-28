package cmd

import (
	"os"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func PassphraseCommand() *cli.Command {
	return &cli.Command{
		Name:    "passphrase",
		Aliases: []string{"pass", "passw", "password", "pw"},
		Usage:   "Will print the passphrase for the account",
		Action: func(c *cli.Context) error {
			passphrase, err := api.GetAccountPassphrase(c.Args().First())
			if err != nil {
				return err
			}

			os.Stdout.WriteString(passphrase + "\n")

			return nil
		},
	}
}

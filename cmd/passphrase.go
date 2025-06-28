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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "The id of the account to get the passphrase for",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			passphrase, err := api.GetAccountPassphrase(c.String("id"))
			if err != nil {
				return err
			}

			os.Stdout.WriteString(passphrase + "\n")

			return nil
		},
	}
}

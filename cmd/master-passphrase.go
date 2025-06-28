package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func ChangeMasterPassphraseCommand() *cli.Command {
	return &cli.Command{
		Name:    "master-passphrase",
		Aliases: []string{"change-passphrase", "change-master", "change-master-pass"},
		Usage:   "Will change the master passphrase.",
		Action: func(c *cli.Context) error {
			// 1. Take new passphrase from user
			passphrase, err := utilities.ReadValue("New passphrase: ", true, true)
			if err != nil {
				return cli.Exit("Failed to read passphrase: "+err.Error(), 1)
			}

			// 2. Ask API to change the master passphrase
			err = api.ChangeMasterPassphrase(passphrase)
			if err != nil {
				return err
			}
			fmt.Println("✅ Master passphrase changed")
			return nil
		},
	}
}

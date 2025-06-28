package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func ChangeMasterPassphraseCommand() *cli.Command {
	return &cli.Command{
		Name:    "master-passphrase",
		Aliases: []string{"change-passphrase", "change-master", "change-master-pass"},
		Usage:   "Will change the master passphrase.",
		Action: func(c *cli.Context) error {
			// 1. Take passphrase from user
			fmt.Print("Enter new master passphrase: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return cli.Exit("Failed to read passphrase: "+err.Error(), 1)
			}
			fmt.Println()
			passphrase := string(bytePassword)
			if passphrase == "" {
				return cli.Exit("Passphrase is required", 1)
			}

			// 2. Ask API to change the master passphrase
			err = api.ChangeMasterPassphrase(passphrase)
			if err != nil {
				return err
			}
			fmt.Println("Master passphrase changed")
			return nil
		},
	}
}

package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func RegisterCommand() *cli.Command {
	return &cli.Command{
		Name:    "register",
		Aliases: []string{"init", "initialize"},
		Usage:   "Initialize the passenger if not already initialized.",
		Action: func(context *cli.Context) error {
			// 1. Take passphrase from user
			passphrase, err := utilities.ReadValue("Passphrase: ", true, true)
			if err != nil {
				return err
			}
			// 2. Ask API to register the system
			recovery, err := api.Register(passphrase)
			if err != nil {
				return err
			}
			// 3. Print recovery key to stdin and the description to stderr
			os.Stdout.WriteString("🚨 Register flow requires you to securely store a recovery key. This key will be required if forget your master passphrase.\n This text printed to stdout, you can redirect stderr to a file to save the recovery key.\n")
			os.Stderr.WriteString(recovery)
			os.Stdout.WriteString("\n")
			return nil
		},
	}
}

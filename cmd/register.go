package cmd

import (
	"fmt"
	"os"
	"passenger-go-cli/internal/api"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func RegisterCommand() *cli.Command {
	return &cli.Command{
		Name:    "register",
		Aliases: []string{"init", "initialize"},
		Usage:   "Initialize the passenger if not already initialized.",
		Action: func(context *cli.Context) error {
			// 1. Take passphrase from user
			fmt.Print("Enter passphrase: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return cli.Exit("Failed to read passphrase: "+err.Error(), 1)
			}
			fmt.Println()
			passphrase := string(bytePassword)
			if passphrase == "" {
				return cli.Exit("Passphrase is required", 1)
			}
			// 2. Ask API to register the system
			recovery, err := api.Register(passphrase)
			if err != nil {
				return err
			}
			// 3. Print recovery key to stdin and the description to stderr
			fmt.Fprintf(os.Stderr, "🚨 Register flow requires you to securely store a recovery key. This key will be required if forget your master passphrase.\n This text printed to stderr, you can redirect stdout to a file to save the recovery key.\n\n")
			fmt.Println(recovery)
			return nil
		},
	}
}

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
			fmt.Print("Enter passphrase: ")

			// Read password from stdin without echoing
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return cli.Exit("❌ Failed to read passphrase: "+err.Error(), 1)
			}
			fmt.Println() // Print newline after password input

			passphrase := string(bytePassword)
			if passphrase == "" {
				return cli.Exit("❌ Passphrase is required", 1)
			}

			// Check if system is already initialized
			status, err := api.GetStatus()
			if err != nil {
				return cli.Exit(fmt.Sprintf("❌ Failed to check initialization status: %s", err.Error()), 1)
			}
			if status {
				return cli.Exit("❌ System is already initialized. Use 'passenger login' instead.", 1)
			}

			recovery, err := api.Register(passphrase)
			if err != nil {
				return cli.Exit(fmt.Sprintf("❌ Registration failed: %s", err.Error()), 1)
			}
			fmt.Fprintf(os.Stderr, "🚨 Register flow requires you to securely store a recovery key. This key will be required if forget your master passphrase.\n This text printed to stderr, you can redirect stdout to a file to save the recovery key.\n\n")
			fmt.Println(recovery)
			return nil
		},
	}
}

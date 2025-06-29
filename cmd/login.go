package cmd

import (
	"fmt"
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/auth"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func LoginCommand() *cli.Command {
	return &cli.Command{
		Name:    "login",
		Aliases: []string{"sign-in", "log-in"},
		Usage:   "Login to the passenger.",
		Action: func(c *cli.Context) error {
			// Check if the server is initialized first
			status, err := api.Status()
			if err != nil {
				return cli.Exit("Failed to check server status: "+err.Error(), 1)
			}

			if !status {
				return cli.Exit(`❌ Cannot login: Passenger Go server is not initialized.

To initialize the server:
1. Run 'passenger-go register' to set up the master passphrase
2. Run 'passenger-go validate' to verify your recovery key
3. Then run 'passenger-go login' to sign in

For more information, run 'passenger-go help register'`, 1)
			}

			passphrase, err := utilities.ReadValue("Passphrase", true, true)
			if err != nil {
				return cli.Exit("Failed to read passphrase: "+err.Error(), 1)
			}

			token, err := api.Login(passphrase)
			if err != nil {
				return cli.Exit("Could not login: "+err.Error(), 1)
			}

			err = auth.StoreToken(token)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to store token: %v", err), 1)
			}

			os.Stdout.WriteString("✅ Successfully logged in! Token will expire in 5 minutes.")
			return nil
		},
	}
}

package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/auth"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func LoginCommand() *cli.Command {
	return &cli.Command{
		Name:    "login",
		Aliases: []string{"sign-in", "log-in"},
		Usage:   "Login to the passenger.",
		Action: func(c *cli.Context) error {
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

			token, err := api.Login(passphrase)
			if err != nil {
				return cli.Exit("❌ Could not login: "+err.Error(), 1)
			}

			err = auth.StoreToken(token)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to store token: %v", err), 1)
			}

			fmt.Println("✅ Successfully logged in! Token will expire in 5 minutes.")
			return nil
		},
	}
}

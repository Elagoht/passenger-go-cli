package cmd

import (
	"fmt"
	"passenger-go-cli/internal/auth"

	"github.com/urfave/cli/v2"
)

func LogoutCommand() *cli.Command {
	return &cli.Command{
		Name:    "logout",
		Aliases: []string{"sign-out", "log-out"},
		Usage:   "Logout from the passenger. Token is already a short lived one. But you can remove it sooner than that.",
		Action: func(c *cli.Context) error {
			err := auth.ClearToken()
			if err != nil {
				return cli.Exit(fmt.Sprintf("❌ Failed to clear token: %v", err), 1)
			}

			fmt.Println("✅ Successfully logged out! Token has been cleared.")
			return nil
		},
	}
}

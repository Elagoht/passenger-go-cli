package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"remove", "rm", "del", "kaboom", "shred"},
		Usage:   "Will delete the account by id",
		Action: func(context *cli.Context) error {
			accountID := context.Args().First()
			if accountID == "" {
				return cli.Exit("Account ID is required", 1)
			}

			err := api.DeleteAccount(accountID)
			if err != nil {
				return cli.Exit("Failed to delete account: "+err.Error(), 1)
			}

			fmt.Println("✅ Account deleted successfully")
			return nil
		},
	}
}

package cmd

import (
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"
	"strconv"

	"github.com/urfave/cli/v2"
)

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:    "get",
		Aliases: []string{"fetch", "show"},
		Usage:   "Will get the account details by id",
		Args:    true,
		Action: func(context *cli.Context) error {
			accountID := context.Args().First()
			if accountID == "" {
				return cli.Exit("Account ID is required, use `passenger-go list` to get the account ID", 1)
			}

			account, err := api.GetAccount(accountID)
			if err != nil {
				return err
			}

			utilities.PrintTable([][]string{
				{"ID", account.ID},
				{"Platform", account.Platform},
				{"Identifier", account.Identifier},
				{"URL", account.URL},
				{"Notes", func() string {
					if account.Notes == "" {
						return "<no-notes-available>"
					} else {
						return account.Notes
					}
				}()},
				{"Strength", strconv.Itoa(account.Strength)},
			}, nil)

			return nil
		},
	}
}

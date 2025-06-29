package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "show-all", "fetch-all", "get-all"},
		Usage:   "Will list all accounts",
		Action: func(c *cli.Context) error {
			accounts, err := api.GetAccounts()
			if err != nil {
				return err
			}

			if len(accounts) == 0 {
				os.Stdout.WriteString("No accounts found, use `passenger-go create` or `passenger-go import --file=<file>` to add data.")
				return nil
			}

			// Convert accounts to string slices for table printing
			var rows [][]string
			for _, account := range accounts {
				row := []string{
					account.ID,
					account.Platform,
					account.Identifier,
					account.URL,
				}
				rows = append(rows, row)
			}

			utilities.PrintTable(rows, []string{"ID", "Platform", "Identifier", "URL"})
			return nil
		},
	}
}

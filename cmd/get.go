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
		Action: func(c *cli.Context) error {
			account, err := api.GetAccount(c.Args().First())
			if err != nil {
				return err
			}

			utilities.PrintTable([][]string{
				{"ID", account.ID},
				{"Platform", account.Platform},
				{"Identifier", account.Identifier},
				{"URL", account.URL},
				{"Notes", account.Notes},
				{"Strength", strconv.Itoa(account.Strength)},
			}, nil)

			return nil
		},
	}
}

package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
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

			os.Stdout.WriteString(account.ID + "\n")
			os.Stdout.WriteString(account.Platform + "\n")
			os.Stdout.WriteString(account.Identifier + "\n")
			os.Stdout.WriteString(account.URL + "\n")
			os.Stdout.WriteString(account.Notes + "\n")
			os.Stdout.WriteString(strconv.Itoa(account.Strength) + "\n")

			return nil
		},
	}
}

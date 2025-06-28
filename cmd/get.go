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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "The id of the account to get",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			account, err := api.GetAccount(c.String("id"))
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

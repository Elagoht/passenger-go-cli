package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func LoginCommand() *cli.Command {
	return &cli.Command{
		Name:    "login",
		Aliases: []string{"sign-in", "log-in"},
		Usage:   "Login to the passenger.",
		Action: func(c *cli.Context) error {
			passphrase := c.Args().First()
			if passphrase == "" {
				return cli.Exit("Passphrase is required", 1)
			}
			token, err := api.Login(passphrase)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			fmt.Println(token)
			return nil
		},
	}
}

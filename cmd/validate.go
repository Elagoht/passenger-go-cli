package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func ValidateCommand() *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"verify"},
		Usage:   "Validate the recovery key. Server needs to verify you have really backed up your recovery key.",
		Action: func(c *cli.Context) error {
			err := api.ValidateRecovery(c.Args().First())
			if err != nil {
				return err
			}

			fmt.Println("Recovery key validated")
			return nil
		},
	}
}

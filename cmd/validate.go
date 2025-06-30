package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func ValidateCommand() *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"verify"},
		Usage:   "Validate the recovery key. Server needs to verify you have really backed up your recovery key.",
		Action: func(c *cli.Context) error {
			recoveryKey, err := utilities.ReadValue("Recovery key", true, true)
			if err != nil {
				return err
			}

			err = api.ValidateRecovery(recoveryKey)
			if err != nil {
				return err
			}

			os.Stderr.WriteString("✅ Recovery key validated\nYou can now login with 'passenger-go login'\n")
			return nil
		},
	}
}

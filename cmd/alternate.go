package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func AlternateCommand() *cli.Command {
	return &cli.Command{
		Name:    "alternate",
		Aliases: []string{"alt", "alternative", "manipulate", "shuffle"},
		Usage:   "Alternate characters with similar looking characters.",
		Action: func(context *cli.Context) error {
			passphrase, err := utilities.ReadValue("Passphrase: ", true, true)
			if err != nil {
				return err
			}

			alternate, err := api.AlternatePassphrase(passphrase)
			if err != nil {
				return err
			}
			os.Stdout.WriteString("Alternate passphrase printed on stderr:\n")
			os.Stderr.WriteString(alternate)
			os.Stdout.WriteString("\n")
			return nil
		},
	}
}

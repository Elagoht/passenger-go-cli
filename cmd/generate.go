package cmd

import (
	"os"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"gen", "suggest", "random"},
		Usage:   "Will generate a random passphrase of the specified length. Default is 32.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "length",
				Aliases: []string{"l"},
				Usage:   "The length of the passphrase to generate. Default is 32.",
			},
		},
		Action: func(c *cli.Context) error {

			length := 32
			if c.IsSet("length") {
				length = c.Int("length")
			}

			passphrase, err := api.GeneratePassphrase(length)
			if err != nil {
				return err
			}
			os.Stdout.WriteString(passphrase)
			return nil
		},
	}
}

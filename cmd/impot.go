package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"strconv"

	"github.com/urfave/cli/v2"
)

func ImportCommand() *cli.Command {
	return &cli.Command{
		Name:    "import",
		Aliases: []string{"import-csv", "load"},
		Usage:   "Will import accounts from a CSV file, only supports Firefox and Chromium.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "file",
				Aliases:   []string{"f", "input", "i"},
				Usage:     "The file to import the CSV from.",
				Required:  true,
				TakesFile: true,
			},
		},
		Action: func(c *cli.Context) error {
			filePath := c.String("file")

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return cli.Exit("File not found: "+filePath, 1)
			}

			response, err := api.ImportCSV(filePath)
			if err != nil {
				return err
			}

			if response.Imported > 0 {
				os.Stdout.WriteString(
					"✅ Imported " + strconv.Itoa(response.Imported) +
						" accounts from " + filePath + "\n",
				)
			}

			if response.Skipped > 0 {
				os.Stdout.WriteString(
					"❌ Skipped " + strconv.Itoa(response.Skipped) +
						" accounts from " + filePath + "\n" +
						"Unimportable accounts (might be already exist) printed to stderr.\n",
				)
			}

			if len(response.Errors) > 0 {
				for _, unimportedAccount := range response.Errors {
					os.Stderr.WriteString(unimportedAccount + "\n")
				}
			}

			return nil
		},
	}
}

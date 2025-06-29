package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/utilities"
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
		Action: func(context *cli.Context) error {
			filePath := context.String("file")

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return cli.Exit("File not found: "+filePath, 1)
			}

			response, err := api.ImportCSV(filePath)
			if err != nil {
				return err
			}

			if response.SuccessCount > 0 {
				os.Stdout.WriteString(
					"✅ Imported " + strconv.Itoa(response.SuccessCount) +
						" accounts from " + filePath + "\n",
				)
			}

			if len(response.FailedOnes) > 0 {
				os.Stdout.WriteString(
					"❌ Skipped " + strconv.Itoa(len(response.FailedOnes)) +
						" accounts from " + filePath + "\n" +
						"Unimportable accounts (might be already exist) printed to stderr.\n",
				)
				var failedTable [][]string
				for _, account := range response.FailedOnes {
					failedTable = append(failedTable, []string{account.Platform, account.Identifier, account.URL})
				}
				utilities.PrintTable(failedTable, []string{"Platform", "Identifier", "URL"}, true)
			}

			return nil
		},
	}
}

package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

func ImportCommand() *cli.Command {
	return &cli.Command{
		Name:    "import",
		Aliases: []string{"import", "import-csv", "load"},
		Usage:   "Will import accounts from a CSV file, only supports Firefox and Chromium.",
		Action: func(c *cli.Context) error {
			filePath := c.Args().Get(0)
			if filePath == "" {
				return cli.Exit("❌ CSV file path is required", 1)
			}

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return cli.Exit("❌ File not found: "+filePath, 1)
			}

			// api.ImportCSV(filePath)
			return nil
		},
	}
}

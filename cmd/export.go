package cmd

import (
	"os"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func ExportCommand() *cli.Command {
	return &cli.Command{
		Name:    "export",
		Aliases: []string{"export-csv", "dump"},
		Usage:   "Will export accounts to a CSV file, exported CSV will be in Chromium format.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "output",
				Aliases:   []string{"o"},
				Usage:     "The file to export the CSV to.",
				Required:  false,
				TakesFile: true,
			},
		},
		Action: func(context *cli.Context) error {
			csvBytes, err := api.ExportCSV()
			if err != nil {
				return err
			}

			output := context.String("output")

			if output == "" {
				os.Stdout.Write(csvBytes)
				os.Stderr.WriteString("✅ Exported CSV to stdout, you can pipe it to a file.\n")
				return nil // Early return to avoid writing to file
			}

			err = os.WriteFile(output, csvBytes, 0644)
			if err != nil {
				return err
			}

			os.Stderr.WriteString("✅ Exported CSV to " + output + "\n")
			return nil
		},
	}
}

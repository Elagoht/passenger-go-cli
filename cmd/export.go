package cmd

import "github.com/urfave/cli/v2"

func ExportCommand() *cli.Command {
	return &cli.Command{
		Name:    "export",
		Aliases: []string{"export-csv", "dump"},
		Usage:   "Will export accounts to a CSV file, exported CSV will be in Chromium format.",
	}
}

package cmd

import (
	"fmt"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"add", "new", "insert"},
		Usage:   "Create a new account with interactive form",
		Action: func(c *cli.Context) error {
			form := utilities.NewInteractiveForm()

			form.AddField("platform", "Platform", false, true)
			form.AddField("identifier", "Identifier", false, true)
			form.AddField("url", "URL", false, false)
			form.AddField("notes", "Notes", false, false)
			form.AddField("passphrase", "Passphrase", true, true)

			err := form.Run()
			if err != nil {
				return cli.Exit("Failed to collect form data: "+err.Error(), 1)
			}

			fmt.Println(form.GetValues())

			return nil
		},
	}
}

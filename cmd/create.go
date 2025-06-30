package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/schemas"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"add", "new", "insert"},
		Usage:   "Create a new account with interactive form",
		Action: func(context *cli.Context) error {
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

			account, err := api.CreateAccount(schemas.UpsertAccountRequest{
				Platform:   form.GetValues()["platform"],
				Identifier: form.GetValues()["identifier"],
				URL:        form.GetValues()["url"],
				Notes:      form.GetValues()["notes"],
				Passphrase: form.GetValues()["passphrase"],
			})
			if err != nil {
				return cli.Exit("Failed to create account: "+err.Error(), 1)
			}

			os.Stdout.WriteString("Account created successfully with Id: " + account.ID + "\n")
			return nil
		},
	}
}

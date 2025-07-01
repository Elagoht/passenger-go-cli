package cmd

import (
	"os"
	"passenger-go-cli/internal/api"
	"passenger-go-cli/internal/schemas"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func UpdateCommand() *cli.Command {
	return &cli.Command{
		Name:    "update",
		Aliases: []string{"edit", "modify", "change"},
		Usage:   "Update an existing account with interactive form",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "Account ID to update",
				Required: true,
			},
		},
		Action: func(context *cli.Context) error {
			accountID := context.String("id")

			// Get the existing account
			existingAccount, err := api.GetAccount(accountID)
			if err != nil {
				return cli.Exit("Failed to get account: "+err.Error(), 1)
			}

			// Get the current passphrase
			currentPassphrase, err := api.GetAccountPassphrase(accountID)
			if err != nil {
				return cli.Exit("Failed to get account passphrase: "+err.Error(), 1)
			}

			form := utilities.NewInteractiveForm()

			form.AddFieldWithDefault("platform", "Platform", existingAccount.Platform, false, true)
			form.AddFieldWithDefault("identifier", "Identifier", existingAccount.Identifier, false, true)
			form.AddFieldWithDefault("url", "URL", existingAccount.URL, false, false)
			form.AddFieldWithDefault("notes", "Notes", existingAccount.Notes, false, false)
			form.AddFieldWithDefault("passphrase", "Passphrase", currentPassphrase, true, true)

			err = form.Run()
			if err != nil {
				return cli.Exit("Failed to collect form data: "+err.Error(), 1)
			}

			values := form.GetValues()

			// Create updated account object
			updatedAccount := schemas.UpsertAccountRequest{
				Platform:   values["platform"],
				Identifier: values["identifier"],
				URL:        values["url"],
				Notes:      values["notes"],
				Passphrase: values["passphrase"],
			}

			// Update the account
			err = api.UpdateAccount(accountID, updatedAccount)
			if err != nil {
				return cli.Exit("Failed to update account: "+err.Error(), 1)
			}

			// Check if passphrase was changed
			newPassphrase := values["passphrase"]
			if newPassphrase != currentPassphrase {
				// Update the passphrase
				err = api.UpdateAccountPassphrase(accountID, newPassphrase)
				if err != nil {
					return cli.Exit("Failed to update passphrase: "+err.Error(), 1)
				}
				os.Stdout.WriteString("Passphrase updated successfully\n")
			}

			return nil
		},
	}
}

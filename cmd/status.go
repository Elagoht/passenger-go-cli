package cmd

import (
	"fmt"
	"passenger-go-cli/internal/api"

	"github.com/urfave/cli/v2"
)

func StatusCommand() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Aliases: []string{"check", "is-initialized"},
		Usage:   "Check if the Passenger Go initialized.",
		Action: func(context *cli.Context) error {
			status, err := api.Status()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			if status {
				fmt.Println("✅ Passenger Go is initialized")
			} else {
				fmt.Println("❌ Passenger Go is not initialized")
			}
			return nil
		},
	}
}

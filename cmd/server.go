package cmd

import (
	"os"
	"passenger-go-cli/internal/config"
	"passenger-go-cli/internal/utilities"

	"github.com/urfave/cli/v2"
)

func ServerCommand() *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"set-server", "set-url", "set-server-url"},
		Usage:   "Where Passenger Go is hosting. Do not include the /api path.",
		Subcommands: []*cli.Command{
			{
				Name:    "get",
				Aliases: []string{"show", "show-url", "get-url"},
				Usage:   "Show currently set server URL.",
				Action: func(context *cli.Context) error {
					configuration, err := config.LoadConfig()
					if err != nil {
						return cli.Exit("Error loading config: "+err.Error(), 0)
					}
					if configuration.ServerURL == "" {
						return cli.Exit("Server URL is not set. Use 'server set <url>' to set it.", 0)
					}
					os.Stdout.WriteString("Server URL is set to " + configuration.ServerURL + "\n")
					return nil
				},
			},
			{
				Name:    "set",
				Aliases: []string{"set-url"},
				Usage:   "Set the server URL.",
				Args:    true,
				Action: func(context *cli.Context) error {
					serverURL, err := utilities.ReadValue("Server URL", false, true)
					if err != nil {
						return err
					}
					config.SaveConfig(&config.Config{ServerURL: serverURL})
					os.Stdout.WriteString("✅ Server URL set to " + serverURL + "\n")
					return nil
				},
			},
		},
		Action: func(context *cli.Context) error {
			return cli.Exit("Please specify either 'server get' to show the current server URL or 'server set' to set a new server URL.", 0)
		},
	}
}

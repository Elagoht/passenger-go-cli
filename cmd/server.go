package cmd

import (
	"fmt"
	"passenger-go-cli/internal/config"

	"github.com/urfave/cli/v2"
)

func ServerCommand() *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"set-server", "set-url", "set-server-url"},
		Usage:   "Where Passenger Go is hosting. Do not include the /api path.",
		Action: func(c *cli.Context) error {
			serverURL := c.Args().First()
			if serverURL == "" {
				return cli.Exit("Server URL is required", 1)
			}
			config.SaveConfig(&config.Config{
				ServerURL: serverURL,
			})
			fmt.Println("Server URL set to", serverURL)
			return nil
		},
	}
}

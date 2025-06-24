package cmd

import "github.com/urfave/cli/v2"

func LogoutCommand() *cli.Command {
	return &cli.Command{
		Name:    "logout",
		Aliases: []string{"sign-out", "log-out"},
		Usage:   "Logout from the passenger. Token is already a short lived one. But you can remove it sooner than that.",
	}
}

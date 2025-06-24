package cmd

import "github.com/urfave/cli/v2"

func RegisterCommand() *cli.Command {
	return &cli.Command{
		Name:    "register",
		Aliases: []string{"init", "initialize"},
		Usage:   "Initialize the passenger if not already initialized.",
	}
}

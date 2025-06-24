package cmd

import "github.com/urfave/cli/v2"

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"gen", "suggest", "random"},
		Usage:   "Will generate a random passphrase of the specified length. Default is 32.",
	}
}

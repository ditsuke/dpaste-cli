package create

import (
	"dpaste-cli/lib/dpaste"
	"github.com/urfave/cli/v2"
)

func GetCommand(client *dpaste.Dpaste) *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"up"},
		Usage:   "Create a new paste!",
		Action: func(c *cli.Context) error {
			return Create(c, client)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "Title of paste, defaults to _quite_ nothing.",
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "File to paste. If not provided, the paste data is expected to be piped from stdin",
				DefaultText: "Again, What does this do?",
			},
			&cli.StringFlag{
				Name:    "syntax",
				Aliases: []string{"s"},
				Usage:   "Syntax of paste. Defaults to none, or if @todo{guess extension}, @todo{auto}",
			},
			&cli.IntFlag{
				Name:        "expire_days",
				Aliases:     []string{"e"},
				Usage:       "Expiry of paste (in days)",
				DefaultText: "What does this do?",
			},
		},
	}
}

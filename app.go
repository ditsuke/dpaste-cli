package main

import (
	"dpaste-cli/lib/dpaste"
	"dpaste-cli/subcommands/create"
	"github.com/urfave/cli/v2"
	"net/http"
)

// getApp() configures and returns a cli.App instance
func getApp() *cli.App {
	var client = new(dpaste.Dpaste)

	return &cli.App{
		Name:    "dpaste cli",
		Usage:   "Do cool stuff with https://dpaste.com",
		Version: "0.1",
		Authors: []*cli.Author{
			{
				Name:  "Tushar",
				Email: "ditsuke@pm.me",
			},
		},
		Flags: getGlobalFlags(),

		// "create" as the default action, although this doesn't allow flags but perhaps that's fine
		Action: func(context *cli.Context) error {
			return create.Create(context, client)
		},
		Commands: []*cli.Command{
			create.GetCommand(client),
		},

		// Get an instance of the dpaste client before the app starts
		Before: func(c *cli.Context) error {
			*client = *(dpaste.New(c.String("token"), &http.Client{}))
			return nil
		},
	}
}

// getGlobalFlags() returns a flag array of the app's global flags
func getGlobalFlags() []cli.Flag {
	// @todo if we use "create" as the default command and support flags we'll need those here too
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Dpaste token, optional.",
			Required: false,
		},
		// config file// @todo
	}
}

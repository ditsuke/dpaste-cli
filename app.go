package main

import (
	"dpaste-cli/cmd/create"
	"dpaste-cli/lib/dpaste"
	"github.com/urfave/cli/v2"
	"net/http"
)

var (
	client *dpaste.Dpaste
)

type ExtendedActionFunc func(*cli.Context, *dpaste.Dpaste) error

// Inject the Dpaste client into an ExtendedActionFunc, allowing it to be used as a cli.ActionFunc
func injectClient(actionFunc ExtendedActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		return actionFunc(c, client)
	}
}

// getApp() configures and returns a cli.App instance
func getApp() *cli.App {

	app := cli.App{
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

		// "create" as the default action, although this doesn't yet allow flags... so
		Action: injectClient(create.Create),

		// Get an instance of the dpaste client before the app starts
		Before: func(c *cli.Context) error {
			client = dpaste.New(c.String("token"), &http.Client{})
			return nil
		},
	}

	// some configuration for the Create Command -- this is hideous 🤮
	create.Command.Action = app.Action
	app.Commands = cli.Commands{
		&create.Command,
		// @todo work on the config cmd
		{
			Name:    "config",
			Aliases: []string{"conf"},
			Usage:   "Change configuration",
		},
	}

	return &app
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

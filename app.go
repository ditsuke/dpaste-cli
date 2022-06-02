package main

import (
	cfm "github.com/ditsuke/dpaste-cli/config"
	"github.com/ditsuke/dpaste-cli/dpaste"
	"github.com/ditsuke/dpaste-cli/subcommands/config"
	"github.com/ditsuke/dpaste-cli/subcommands/create"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
)

// getApp() configures and returns a cli.App instance
func getApp() *cli.App {
	var client = new(dpaste.Dpaste)
	cfg, err := cfm.LoadConfig(cfm.DefaultConfigFile())
	if err != nil {
		panic(err)
	}

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
		Flags: globalFlags(cfg),

		// "create" as the default action, although this doesn't allow flags but perhaps that's fine
		Action: func(context *cli.Context) error {
			return create.Create(context, client)
		},
		ExitErrHandler: func(context *cli.Context, err error) {
			logger := log.New(context.App.ErrWriter, "", 0)
			logger.Println("Error:", err.Error())
		},
		Commands: []*cli.Command{
			create.GetCommand(client),
			config.GetCommand(),
		},

		// Get an instance of the dpaste client before the app starts
		Before: func(c *cli.Context) error {
			*client = *(dpaste.New(c.String(cfm.FlagDpasteToken), &http.Client{}))
			return nil
		},
	}
}

// globalFlags() returns the app's top level flags
func globalFlags(cfg cfm.Config) []cli.Flag {
	// @todo if we use "create" as the default command and support flags we'll need those here too
	return []cli.Flag{
		&cli.StringFlag{
			Name:     cfm.FlagDpasteToken,
			Aliases:  []string{"t"},
			Value:    cfg.Token,
			EnvVars:  []string{cfm.EnvDpasteToken},
			Usage:    "optional dpaste.com user token",
			Required: false,
		},
		// config file// @todo
	}
}

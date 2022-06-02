package config

import (
	"fmt"
	"github.com/ditsuke/dpaste-cli/config"
	"github.com/urfave/cli/v2"
)

const (
	FlagConfigFile = "config"
)

const (
	ExitCodeUsage             = 1
	ExitCodeReadConfigFailure = 2
)

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "configure dpaste-cli",
		Subcommands: []*cli.Command{
			{
				Name:      "get",
				Usage:     "get config",
				ArgsUsage: "KEY",
				Action: func(c *cli.Context) error {
					stdout := c.App.Writer

					args := c.Args()
					opt := args.Get(0)
					if opt == "" {
						return cli.Exit("config requires a key", ExitCodeUsage)
					}
					cfg, err := config.LoadConfig(c.String(FlagConfigFile))
					if err != nil {
						return cli.Exit(err.Error(), ExitCodeReadConfigFailure)
					}
					val, err := cfg.Get(opt)
					if err != nil {
						return cli.Exit(err.Error(), ExitCodeUsage)
					}
					_, _ = fmt.Fprintln(stdout, val)
					return nil
				},
			},
			{
				Name:      "set",
				Usage:     "set config",
				ArgsUsage: "KEY VALUE",
				Action: func(c *cli.Context) error {
					args := c.Args()
					opt := args.Get(0)
					if opt == "" {
						return cli.Exit("config requires a key", ExitCodeUsage)
					}
					val := args.Get(1)
					if val == "" {
						return cli.Exit("config requires a value", ExitCodeUsage)
					}
					cfg, err := config.LoadConfig(c.String(FlagConfigFile))
					if err != nil {
						return cli.Exit(err.Error(), ExitCodeReadConfigFailure)
					}
					err = cfg.Set(opt, val)
					if err != nil {
						return cli.Exit(err.Error(), ExitCodeUsage)
					}
					return cfg.Write(c.String(FlagConfigFile))
				},
			},
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    FlagConfigFile,
				Aliases: []string{"c"},
				Usage:   fmt.Sprintf("config file path. default: %s", config.DefaultConfigFile()),
				Value:   config.DefaultConfigFile(),
			},
		},
	}
}

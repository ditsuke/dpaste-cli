package config

import (
	"fmt"
	"github.com/ditsuke/dpaste-cli/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const (
	FlagConfigFile = "config"
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
					args := c.Args()
					opt := args.Get(0)
					if opt == "" {
						return fmt.Errorf("key is required")
					}
					cfg, err := config.LoadConfig(c.String(FlagConfigFile))
					if err != nil {
						return err
					}
					val, err := cfg.Get(opt)
					if err != nil {
						return err
					}
					fmt.Println(val)
					return nil
				},
				OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
					logger := log.New(os.Stderr, "", log.LstdFlags)
					logger.Println("Error:", err.Error())
					return fmt.Errorf("config error: %w", err)
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
						return fmt.Errorf("config requires a key")
					}
					val := args.Get(1)
					if val == "" {
						return fmt.Errorf("config requires a value")
					}
					cfg, err := config.LoadConfig(c.String(FlagConfigFile))
					if err != nil {
						return err
					}
					err = cfg.Set(opt, val)
					if err != nil {
						return err
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

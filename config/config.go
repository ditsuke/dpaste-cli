package config

import (
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	ErrUnknownKey   = "unknown config key"
	ErrInvalidValue = "invalid value for key"
)

const (
	ExitUnknownKey   = 1
	ExitInvalidValue = 2
)

const (
	Token       = "token"
	GuessSyntax = "guess-syntax"
)

type Config struct {
	Token       string `yaml:"token,omitempty"`
	GuessSyntax bool   `yaml:"guess-syntax"`
}

func (c *Config) Set(key, value string) error {
	switch key {
	case Token:
		c.Token = value
	case GuessSyntax:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return cli.Exit(ErrInvalidValue, ExitInvalidValue)
		}
		c.GuessSyntax = b
	default:
		return cli.Exit(ErrUnknownKey, ExitUnknownKey)
	}
	return nil
}

func (c Config) Get(key string) (string, error) {
	switch key {
	case Token:
		return c.Token, nil
	case GuessSyntax:
		return strconv.FormatBool(c.GuessSyntax), nil
	default:
		return "", cli.Exit(ErrUnknownKey, ExitUnknownKey)
	}
}

func (c Config) MustGet(key string) string {
	v, err := c.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c Config) Write(path string) (err error) {
	configBytes, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	err = os.WriteFile(path, configBytes, 0600)
	return
}

func DefaultConfigFile() string {
	const DefaultPath = "{HOME_DIR}/.dpaste.yml"
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return strings.Replace(DefaultPath, "{HOME_DIR}", home, -1)
}

func LoadConfig(path string) (config Config, err error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil && err.(*fs.PathError) == nil {
		return
	}
	if len(configBytes) == 0 {
		err = nil
		return
	}
	err = yaml.Unmarshal(configBytes, &config)
	return
}

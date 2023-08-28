package config

import (
	"context"
	"github.com/sabahtalateh/gic"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
}

func init() {
	gic.Add[*Config](
		gic.WithInitE(func() (*Config, error) {
			bb, err := os.ReadFile("./config.yaml")
			if err != nil {
				return nil, err
			}

			conf := new(Config)
			if err = yaml.Unmarshal(bb, conf); err != nil {
				return nil, err
			}

			return conf, nil
		}),
		gic.WithStart(func(ctx context.Context, t *Config) error {
			return nil
		}),
	)
}

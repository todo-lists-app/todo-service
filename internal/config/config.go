package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	ConfigBuilder "github.com/keloran/go-config"
)

type Config struct {
	Local
	ConfigBuilder.Config
}

func Build() (*Config, error) {
	cfg := &Config{}

	if err := BuildLocal(cfg); err != nil {
		return nil, logs.Errorf("build local: %v", err)
	}

	c, err := ConfigBuilder.Build(ConfigBuilder.Vault, ConfigBuilder.Mongo)
	if err != nil {
		return nil, logs.Errorf("build config: %v", err)
	}
	cfg.Config = *c

	return cfg, nil
}

package config

import "github.com/caarlos0/env/v8"

type Config struct {
	Local
	Vault
	Mongo
	Notifications
	Ping
}

func Build() (*Config, error) {
	cfg := &Config{}

	if err := BuildLocal(cfg); err != nil {
		return nil, err
	}

	if err := BuildVault(cfg); err != nil {
		return nil, err
	}

	if err := BuildMongo(cfg); err != nil {
		return nil, err
	}

	if err := BuildNotifications(cfg); err != nil {
		return nil, err
	}

	if err := BuildPing(cfg); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

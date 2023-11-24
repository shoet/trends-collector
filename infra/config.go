package main

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	SecretsManagerSecretId string `env:"SECRETS_MANAGER_SECRET_ID"`
	KmsKeyId               string `env:"KMS_KEY_ID"`
}

func NewConfig() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return cfg, nil
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

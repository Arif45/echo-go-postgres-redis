package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppEnv   string         `json:"app_env"`
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

var app *Config

func LoadConfig() error {
	file := os.Getenv("CONFIG_JSON")
	if file == "" {
		file = "config.development.json"
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file, err)
	}

	app = &Config{}
	if err := json.Unmarshal(data, app); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("[config] loaded: %s (env: %s)\n", file, app.AppEnv)
	return nil
}

func GetConfig() *Config {
	return app
}

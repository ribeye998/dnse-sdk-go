// Package config loads DNSE API credentials from environment variables.
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds connection parameters for the DNSE REST and WebSocket APIs.
type Config struct {
	BaseURL   string
	WSURL     string
	APIKey    string
	APISecret string
}

// FromEnv reads configuration from environment variables.
// It silently attempts to load a .env file first.
func FromEnv() (*Config, error) {
	_ = godotenv.Load()
	c := &Config{
		BaseURL:   envOrDefault("DNSE_BASE_URL", "https://openapi.dnse.com.vn"),
		WSURL:     envOrDefault("DNSE_WS_URL", "wss://ws-openapi.dnse.com.vn"),
		APIKey:    os.Getenv("DNSE_API_KEY"),
		APISecret: os.Getenv("DNSE_API_SECRET"),
	}
	return c, c.Validate()
}

// Validate returns an error if any required field is missing.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("DNSE_API_KEY is required")
	}
	if c.APISecret == "" {
		return errors.New("DNSE_API_SECRET is required")
	}
	return nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

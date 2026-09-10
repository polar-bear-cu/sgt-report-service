package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	JWTSecret        string
	SubscriptionAddr string
	SwaggerEnabled   bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Port:             env("PORT", "8080"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		SubscriptionAddr: env("SUBSCRIPTION_SERVICE_ADDR", "localhost:50051"),
		SwaggerEnabled:   env("ENABLE_SWAGGER", "false") == "true",
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET required")
	}
	return cfg, nil
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

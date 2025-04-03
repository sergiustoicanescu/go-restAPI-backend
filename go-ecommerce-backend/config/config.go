package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBProd    string
	DBTest    string
	JWTSecret string
}

func LoadConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		DBProd:    os.Getenv("DB_PROD"),
		DBTest:    os.Getenv("DB_TEST"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.DBProd == "" {
		return nil, fmt.Errorf("DB_PROD environment variable is required")
	}
	if cfg.DBTest == "" {
		return nil, fmt.Errorf("DB_TEST environment variable is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	return cfg, nil
}

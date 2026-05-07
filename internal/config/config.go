package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	APIKey      string
	LogLevel    string
}

func New() (*Config, error) {
	godotenv.Load()

	port := os.Getenv("PORT")
	environment := os.Getenv("APP_ENV")
	apiKey := os.Getenv("GEMINI_API_KEY")
	logLevel := os.Getenv("LOG_LEVEL")

	if port == "" {
		port = "8080"
	}

	if apiKey == "" {
		return nil, errors.New("Missing Gemini API Key environment variable")
	}

	return &Config{
		Port:        port,
		Environment: environment,
		APIKey:      apiKey,
		LogLevel:    logLevel,
	}, nil
}

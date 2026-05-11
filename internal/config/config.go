package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	JWTSecret   string
	JWTExpiry   time.Duration
	APIModel    string
	APIKey      string
	LogLevel    string
}

func New() (*Config, error) {
	godotenv.Load()

	port := os.Getenv("PORT")
	environment := os.Getenv("APP_ENV")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry := os.Getenv("JWT_EXPIRY_IN_SECONDS")
	apiModel := os.Getenv("GEMINI_MODEL")
	apiKey := os.Getenv("GEMINI_API_KEY")
	logLevel := os.Getenv("LOG_LEVEL")

	if port == "" {
		port = "8080"
	}

	if jwtSecret == "" {
		return nil, errors.New("Missing JWT Secret (JWT_SECRET) environment variable")
	}

	if jwtExpiry == "" {
		return nil, errors.New("Missing JWT Expiry (JWT_EXPIRY_IN_SECONDS) environment variable")
	}

	expiry, err := strconv.Atoi(jwtExpiry)
	if err != nil {
		return nil, errors.New("Invalid JWT Expiry")
	}

	if apiModel == "" {
		return nil, errors.New("Missing Gemini Model (GEMINI_MODEL) environment variable")
	}

	if apiKey == "" {
		return nil, errors.New("Missing Gemini API Key (GEMINI_API_KEY) environment variable")
	}

	return &Config{
		Port:        port,
		Environment: environment,
		JWTSecret:   jwtSecret,
		JWTExpiry:   time.Duration(expiry) * time.Second,
		APIModel:    apiModel,
		APIKey:      apiKey,
		LogLevel:    logLevel,
	}, nil
}

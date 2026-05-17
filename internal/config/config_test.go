package config

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func setValidEnv(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_ACCESS_SECRET", "secret")
	t.Setenv("JWT_ACCESS_EXPIRY_IN_SECONDS", "1")
	t.Setenv("REFRESH_EXPIRY_IN_SECONDS", "10")
	t.Setenv("GEMINI_MODEL", "tldr-v1")
	t.Setenv("GEMINI_API_KEY", "SeCrEtApIkEy")
	t.Setenv("LOG_LEVEL", "info")
}

func TestNew_ValidConfig(t *testing.T) {
	setValidEnv(t)

	expected := &Config{
		Port:          "8080",
		Environment:   "production",
		JWTSecret:     "secret",
		JWTExpiry:     1 * time.Second,
		RefreshExpiry: 10 * time.Second,
		APIModel:      "tldr-v1",
		APIKey:        "SeCrEtApIkEy",
		LogLevel:      "info",
	}

	cfg, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Fatalf("want %v, got %v", expected, cfg)
	}
}

func TestNew_DefaultPort(t *testing.T) {
	setValidEnv(t)
	t.Setenv("PORT", "")

	expected := &Config{
		Port:          "8080",
		Environment:   "production",
		JWTSecret:     "secret",
		JWTExpiry:     1 * time.Second,
		RefreshExpiry: 10 * time.Second,
		APIModel:      "tldr-v1",
		APIKey:        "SeCrEtApIkEy",
		LogLevel:      "info",
	}

	cfg, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Fatalf("want %v, got %v", expected, cfg)
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	var tests = []struct {
		name     string
		override string
		expected error
	}{
		{
			name:     "No app environment",
			override: "APP_ENV",
			expected: errors.New(
				"Missing Application Environment (APP_ENV) environment variable",
			),
		},
		{
			name:     "No JWT secret",
			override: "JWT_ACCESS_SECRET",
			expected: errors.New(
				"Missing JWT Access Secret (JWT_ACCESS_SECRET) environment variable",
			),
		},
		{
			name:     "No JWT expiry",
			override: "JWT_ACCESS_EXPIRY_IN_SECONDS",
			expected: errors.New(
				"Missing JWT Expiry (JWT_ACCESS_EXPIRY_IN_SECONDS) environment variable",
			),
		},
		{
			name:     "No refresh token expiry",
			override: "REFRESH_EXPIRY_IN_SECONDS",
			expected: errors.New(
				"Missing Refresh Token Expiry (REFRESH_EXPIRY_IN_SECONDS) environment variable",
			),
		},
		{
			name:     "No API model",
			override: "GEMINI_MODEL",
			expected: errors.New(
				"Missing Gemini Model (GEMINI_MODEL) environment variable",
			),
		},
		{
			name:     "No API key",
			override: "GEMINI_API_KEY",
			expected: errors.New(
				"Missing Gemini API Key (GEMINI_API_KEY) environment variable",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv(tt.override, "")
			_, err := New()
			if err.Error() != tt.expected.Error() {
				t.Fatalf("want %q, got %q", tt.expected.Error(), err.Error())
			}
		})
	}
}

func TestNew_InvalidExpiry(t *testing.T) {
	setValidEnv(t)
	var tests = []struct {
		name     string
		override string
		expected error
	}{
		{
			name:     "Invalid access token expiry",
			override: "JWT_ACCESS_EXPIRY_IN_SECONDS",
			expected: errors.New("Invalid JWT Expiry"),
		},
		{
			name:     "Invalid refresh token expiry",
			override: "REFRESH_EXPIRY_IN_SECONDS",
			expected: errors.New("Invalid Refresh Token Expiry"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv(tt.override, "abc")
			_, err := New()
			if err.Error() != tt.expected.Error() {
				t.Fatalf("want %q, got %q", tt.expected.Error(), err.Error())
			}
		})
	}
}

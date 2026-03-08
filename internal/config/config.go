package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tuananhlai/brevity-go/internal/telemetry"
)

type AppConfig struct {
	Port          string `env:"PORT"`
	DatabaseURL   string `env:"DATABASE_URL"`
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	// LLMAPIKey is the API key used to generate **all** articles.
	// This field might be removed once llm api key management feature
	// is developed.
	LLMAPIKey string `env:"LLM_API_KEY"`
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig

	// Read variables **only from the environment** and populate config struct fields.
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("error reading configuration: %v", err)
	}

	return &config, nil
}

func MustLoadConfig() *AppConfig {
	config, err := LoadConfig()
	if err != nil {
		telemetry.Logger("github.com/tuananhlai/brevity-go/internal/config").Error(
			"failed to read app config", "error", err)
		os.Exit(1)
	}
	return config
}

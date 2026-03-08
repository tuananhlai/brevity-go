package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Server struct {
		Port string `env:"SERVER_PORT"`
	}
	Database struct {
		URL string `env:"DATABASE_URL"`
	}
	LLM struct {
		BaseURL string `env:"LLM_BASE_URL"`
		APIKey  string `env:"LLM_API_KEY"`
		ModelID string `env:"LLM_MODEL_ID"`
	}
	Encryption struct {
		Key string `env:"ENCRYPTION_KEY"`
	}
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
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}

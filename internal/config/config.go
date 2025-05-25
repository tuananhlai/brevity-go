package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Mode string

const (
	ModeDev     Mode = "dev"
	ModeRelease Mode = "release"
)

type AppConfig struct {
	// If `mode` is `dev`, the server will run in a way that is easier for development.
	// Otherwise, it will be optimized for better performance and data safety.
	Mode   Mode `yaml:"mode" env:"MODE" env-default:"dev"`
	Server struct {
		Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	}
	Database struct {
		URL string `yaml:"url" env:"DATABASE_URL"`
	}
	LLM struct {
		BaseURL string `yaml:"base_url" env:"LLM_BASE_URL"`
		APIKey  string `yaml:"api_key" env:"LLM_API_KEY"`
		ModelID string `yaml:"model_id" env:"LLM_MODEL_ID"`
	}
	Otel struct {
		CollectorGrpcURL string `yaml:"collector_grpc_url" env:"OTEL_COLLECTOR_GRPC_URL"`
	}
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig

	err := cleanenv.ReadConfig("config.yaml", &config)
	if err != nil {
		return nil, err
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

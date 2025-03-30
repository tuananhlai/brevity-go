package config

import (
	"log"

	"github.com/spf13/viper"
)

type Mode string

const (
	ModeDev     Mode = "dev"
	ModeRelease Mode = "release"
)

type AppConfig struct {
	// If `mode` is `dev`, the server will run in a way that is easier for development.
	// Otherwise, it will be optimized for better performance and data safety.
	Mode     Mode `mapstructure:"mode"`
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`
	LLM struct {
		BaseURL string `mapstructure:"base_url"`
		APIKey  string `mapstructure:"api_key"`
		ModelID string `mapstructure:"model_id"`
	} `mapstructure:"llm"`
	Otel struct {
		CollectorGrpcURL string `mapstructure:"collector_grpc_url"`
	} `mapstructure:"otel"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigFile("./config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
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

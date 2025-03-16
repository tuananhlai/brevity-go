package config

import "github.com/spf13/viper"

type AppConfig struct {
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`
	DeepseekAPIKey string `mapstructure:"deepseek_api_key"`
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

package config

type AppConfig struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
	DeepseekAPIKey string `yaml:"deepseek_api_key"`
}

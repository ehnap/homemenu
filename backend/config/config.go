package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	LLM    LLMConfig    `yaml:"llm"`
	Upload UploadConfig `yaml:"upload"`
	Share  ShareConfig  `yaml:"share"`
}

type ServerConfig struct {
	Port      int    `yaml:"port"`
	JWTSecret string `yaml:"jwt_secret"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	Path   string `yaml:"path"`
}

type LLMConfig struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

type UploadConfig struct {
	Dir       string `yaml:"dir"`
	MaxSizeMB int    `yaml:"max_size_mb"`
}

type ShareConfig struct {
	BaseURL string `yaml:"base_url"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{Port: 8080, JWTSecret: "change-me-in-production"},
		DB:     DBConfig{Driver: "sqlite", Path: "./data/homemenu.db"},
		Upload: UploadConfig{Dir: "./data/uploads", MaxSizeMB: 10},
		Share:  ShareConfig{BaseURL: "http://localhost:8080"},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

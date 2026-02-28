package config

import (
	"errors"
	"os"
	"strconv"

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
	cfg := &Config{
		Server: ServerConfig{Port: 8080, JWTSecret: "change-me-in-production"},
		DB:     DBConfig{Driver: "sqlite", Path: "./data/homemenu.db"},
		Upload: UploadConfig{Dir: "./data/uploads", MaxSizeMB: 10},
		Share:  ShareConfig{BaseURL: "http://localhost:8080"},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		// config.yaml not found — continue with defaults
	} else {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Environment variable overrides
	if v := os.Getenv("HOMEMENU_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("HOMEMENU_JWT_SECRET"); v != "" {
		cfg.Server.JWTSecret = v
	}
	if v := os.Getenv("HOMEMENU_DB_PATH"); v != "" {
		cfg.DB.Path = v
	}
	if v := os.Getenv("HOMEMENU_LLM_BASE_URL"); v != "" {
		cfg.LLM.BaseURL = v
	}
	if v := os.Getenv("HOMEMENU_LLM_API_KEY"); v != "" {
		cfg.LLM.APIKey = v
	}
	if v := os.Getenv("HOMEMENU_LLM_MODEL"); v != "" {
		cfg.LLM.Model = v
	}
	if v := os.Getenv("HOMEMENU_UPLOAD_DIR"); v != "" {
		cfg.Upload.Dir = v
	}
	if v := os.Getenv("HOMEMENU_SHARE_BASE_URL"); v != "" {
		cfg.Share.BaseURL = v
	}

	return cfg, nil
}

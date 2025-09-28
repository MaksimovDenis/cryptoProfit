package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		WsPort       string `yaml:"ws_port"`
		RetryTimeout int    `yaml:"retry_timeout"`
		Timeout      int    `yaml:"timeout"`
		LogLevel     string `yaml:"log_level"`
	} `yaml:"server"`
	Binance struct {
		Address string `yaml:"address"`
		Limit   int    `yaml:"limit"`
		Burst   int    `yaml:"burst"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"binance_service"`
}

func LoadConfig(filename string) (*Config, error) {
	cleanPath := filepath.Clean(filename)

	f, err := os.Open(cleanPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

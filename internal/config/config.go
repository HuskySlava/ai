package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	HttpTimeoutSeconds int `yaml:"httpTimeoutSeconds"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)

	return &cfg, err
}

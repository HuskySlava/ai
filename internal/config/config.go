package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Models struct {
	Gemini string `yaml:"gemini"`
	OpenAI string `yaml:"openai"`
}

type Prompts struct {
	Rewrite   string `yaml:"rewrite"`
	Translate string `yaml:"translate"`
}

type Config struct {
	HttpTimeoutSeconds int     `yaml:"httpTimeoutSeconds"`
	Models             Models  `yaml:"models"`
	Prompts            Prompts `yaml:"prompts"`
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

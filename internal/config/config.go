package config

import (
	"gopkg.in/yaml.v3" // <--- The import
	"os"
)

type Config struct {
	HttpTimeoutSeconds string `yaml:"httpTimeoutSeconds"`
}

func Load(path string) (*Config, error) {
	// 2. Read the file from disk
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 3. Unmarshal (Parse) bytes into the struct
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)

	return &cfg, err
}

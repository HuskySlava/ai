package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Models struct {
	Gemini string `yaml:"gemini"`
	Openai string `yaml:"openai"`
	Ollama string `yaml:"ollama"`
	Claude string `yaml:"claude"`
}

type Prompts struct {
	Rewrite   string `yaml:"rewrite"`
	Translate string `yaml:"translate"`
	Summarize string `yaml:"summarize"`
}

type BaseEndpoints struct {
	Gemini string `yaml:"gemini"`
	Ollama string `yaml:"ollama"`
	Openai string `yaml:"openai"`
	Claude string `yaml:"claude"`
}

type Config struct {
	HttpTimeoutSeconds int           `yaml:"httpTimeoutSeconds"`
	Models             Models        `yaml:"models"`
	Prompts            Prompts       `yaml:"prompts"`
	BaseEndpoints      BaseEndpoints `yaml:"baseEndpoint"`
	InputFileLimitKB   int           `yaml:"inputFileLimitKB"`
}

func Load() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		return nil, err
	}

	exeDir := filepath.Dir(exePath)
	path := filepath.Join(exeDir, "config.yaml")

	// Try to load config.yaml relative to the binary
	res, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading config from binary dir:", err)
		log.Println("Trying to read the config from current directory (pwd)")

		// Fallback to read config from pwd (for dev with go run ...)
		res, err = os.ReadFile("config.yaml")
		if err != nil {
			log.Println("Error reading config from pwd:", err)
			return nil, err
		}
	}

	var cfg *Config
	if err := yaml.Unmarshal(res, &cfg); err != nil {
		log.Println("Error reading YAML:", err)
		return nil, err
	}

	return cfg, nil
}

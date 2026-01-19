package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(path string, sizeLimitKB int) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(file) > sizeLimitKB*1024 {
		return "", fmt.Errorf("file too large (%d bytes)", len(file))
	}

	// Check nul bytes, verify accidental executable not passed
	for _, b := range file {
		if b == 0 {
			return "", fmt.Errorf("incorrect file type")
		}
	}

	return string(file), nil
}

func WriteFile(name string, data []byte) error {

	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	if err := os.WriteFile(name, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

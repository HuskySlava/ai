package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(path string, sizeLimitKB int) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if info.Size() > int64(sizeLimitKB*1024) {
		return "", fmt.Errorf("file too large (%d bytes)", info.Size())
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
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

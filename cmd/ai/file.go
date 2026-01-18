package main

import (
	"fmt"
	"os"
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

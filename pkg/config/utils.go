package config

import (
	"errors"
	"os"
	"strings"
)

func expandUser(path string) (string, error) {
	// Expand variables
	path = os.ExpandEnv(path)

	if path[:2] == "~/" {
		homeDir := os.Getenv("HOME") // *nix
		if homeDir == "" {           // Windows
			homeDir = os.Getenv("USERPROFILE")
		}
		if homeDir == "" {
			return "", errors.New("user home directory not found")
		}

		return strings.Replace(path, "~", homeDir, 1), nil
	}
	return path, nil
}

// expandField expands environment variables in field
func expandField(input string) string {
	if input == "" {
		return ""
	}
	return os.ExpandEnv(input)
}

package config

import (
	"errors"
	"os"
	"strings"
)

// GetHomeDir returns '~' as a path
func GetHomeDir() string {
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		return homeDir
	}
	if homeDir := os.Getenv("USERPROFILE"); homeDir != "" {
		return homeDir
	}
	return ""
}

func expandUser(path string) (string, error) {
	// Expand variables
	path = os.ExpandEnv(path)

	if path[:2] == "~/" {
		homeDir := GetHomeDir()
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

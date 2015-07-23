package utils

import (
	"errors"
	"os"
	"strings"
)

func ExpandUser(path string) (string, error) {
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

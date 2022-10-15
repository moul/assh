package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// EscapeSpaces escapes all space characters with a backslash (\)
func EscapeSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "\\ ")
}

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

// ExpandEnvSafe replaces ${var} or $var in the string according to the values
// of the current environment variables.
// As opposed to os.ExpandEnv, ExpandEnvSafe won't remove the dollar in '$(...)'
// See https://golang.org/src/os/env.go?s=963:994#L22 for the original function
func ExpandEnvSafe(s string) string {
	buf := make([]byte, 0, 2*len(s))
	i := 0
	for j := 0; j < len(s); j++ {
		// the following line is the only one changing
		if s[j] == '$' && j+1 < len(s) && s[j+1] != '(' {
			buf = append(buf, s[i:j]...)
			name, w := getShellName(s[j+1:])
			buf = append(buf, os.Getenv(name)...)
			j += w
			i = j + 1
		}
	}
	return string(buf) + s[i:]
}

// ExpandUser expands tild and env vars in unix paths
func ExpandUser(path string) (string, error) {
	// Expand variables
	path = ExpandEnvSafe(path)

	if strings.HasPrefix(path, "~/") {
		homeDir := GetHomeDir()
		if homeDir == "" {
			return "", errors.New("user home directory not found")
		}

		path = strings.Replace(path, "~", homeDir, 1)
	}

	path = filepath.FromSlash(path) // OS-agnostic slashes
	path = EscapeSpaces(path)

	return path, nil
}

// ExpandField expands environment variables in field
func ExpandField(input string) string {
	if input == "" {
		return ""
	}
	return ExpandEnvSafe(input)
}

// ExpandSliceField expands environment variables in every entries of a slice field
func ExpandSliceField(input []string) []string {
	ret := []string{}
	for _, entry := range input {
		ret = append(ret, ExpandField(entry))
	}
	return ret
}

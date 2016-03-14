package config

import "strings"

// BoolVal returns a boolean matching a configuration string
func BoolVal(input string) bool {
	input = strings.ToLower(input)
	trueValues := []string{"yes", "ok", "true", "1", "enabled"}
	for _, val := range trueValues {
		if val == input {
			return true
		}
	}
	return false
}

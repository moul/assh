package config

import (
	"bytes"
	"fmt"
	"strings"
)

func isDynamicHostname(hostname string) bool {
	return strings.Contains(hostname, `*`) ||
		strings.Contains(hostname, `[`) ||
		strings.Contains(hostname, `]`)
}

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

// stringComment splits comment strings into <1024 char lines
func stringComment(name, value string) string {
	maxLength := 1024 - len(name) - 9
	ret := []string{}
	for _, line := range splitSubN(value, maxLength) {
		ret = append(ret, fmt.Sprintf("  # %s: %s", name, line))
	}
	return strings.Join(ret, "\n") + "\n"
}

// sliceComment splits comment strings into <1024 char lines
func sliceComment(name string, slice []string) string {
	var (
		bundles   [][]string
		bundleIdx = 0
		curLen    = 0
		maxLength = 1024 - len(name) - 12
	)
	bundles = append(bundles, []string{})

	for _, item := range slice {
		for _, line := range strings.Split(item, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if curLen+len(line) >= maxLength {
				bundleIdx++
				bundles = append(bundles, []string{})
				curLen = 0
			}
			bundles[bundleIdx] = append(bundles[bundleIdx], line)
			curLen += len(line) + 2
		}
	}

	ret := []string{}
	for _, bundle := range bundles {
		ret = append(ret, fmt.Sprintf("  # %s: [%s]", name, strings.Join(bundle, ", ")))
	}
	return strings.Join(ret, "\n") + "\n"
}

// splitSubN splits a string by length
// from: http://stackoverflow.com/questions/25686109/split-string-by-length-in-golang
func splitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}
	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}

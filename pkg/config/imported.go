// This file contains imported functions
// The license and copyright is reported for each functions in the comments.

package config

// Imported and unmodified from https://golang.org/src/os/env.go
// Function under the BSD-License - Copyrighted by the Go Authors
// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// Imported and unmodified from https://golang.org/src/os/env.go
// Function under the BSD-License - Copyrighted by the Go Authors
// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// Imported and unmodified from https://golang.org/src/os/env.go
// Function under the BSD-License - Copyrighted by the Go Authors
// getShellName returns the name that begins the string and the number of bytes
// consumed to extract it.  If the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name.
func getShellName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; just eat the brace.
	case isShellSpecialVar(s[0]):
		return s[0:1], 1
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}
	return s[:i], i
}

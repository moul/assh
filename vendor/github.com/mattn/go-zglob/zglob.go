package zglob

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	envre = regexp.MustCompile(`^(\$[a-zA-Z][a-zA-Z0-9_]+|\$\([a-zA-Z][a-zA-Z0-9_]+\))$`)
)

func Glob(pattern string) ([]string, error) {
	globmask := ""
	root := ""
	matches := []string{}
	relative := !filepath.IsAbs(pattern)
	for n, i := range strings.Split(filepath.ToSlash(pattern), "/") {
		if root == "" && strings.Index(i, "*") != -1 {
			if globmask == "" {
				root = "."
			} else {
				root = filepath.ToSlash(globmask)
			}
		}
		if n == 0 && i == "~" {
			if runtime.GOOS == "windows" {
				i = os.Getenv("USERPROFILE")
			} else {
				i = os.Getenv("HOME")
			}
		}
		if envre.MatchString(i) {
			i = strings.Trim(strings.Trim(os.Getenv(i[1:]), "()"), `"`)
		}

		globmask = filepath.Join(globmask, i)
		if n == 0 {
			if runtime.GOOS == "windows" && filepath.VolumeName(i) != "" {
				globmask = i + "/"
			} else if len(globmask) == 0 {
				globmask = "/"
			}
		}
	}
	if root == "" {
		_, err := os.Stat(pattern)
		if err != nil {
			return nil, os.ErrNotExist
		}
		return []string{pattern}, nil
	}
	if globmask == "" {
		globmask = "."
	}
	globmask = filepath.ToSlash(filepath.Clean(globmask))

	cc := []rune(globmask)
	dirmask := ""
	filemask := ""
	for i := 0; i < len(cc); i++ {
		if cc[i] == '*' {
			if i < len(cc)-2 && cc[i+1] == '*' && cc[i+2] == '/' {
				filemask += "(.*/)?"
				dirmask = filemask
				i += 2
			} else {
				filemask += "[^/]*"
			}
		} else {
			c := cc[i]
			if c == '/' || ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || 255 < c {
				filemask += string(c)
			} else {
				filemask += fmt.Sprintf("[\\x%02X]", c)
			}
			if c == '/' && dirmask == "" && strings.Index(filemask, "*") != -1 {
				dirmask = filemask
			}
		}
	}
	if dirmask == "" {
		dirmask = filemask
	}
	if len(filemask) > 0 && filemask[len(filemask)-1] == '/' {
		if root == "" {
			root = filemask
		}
		filemask += "[^/]*"
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		dirmask = "(?i:" + dirmask + ")"
		filemask = "(?i:" + filemask + ")"
	}
	dre := regexp.MustCompile("^" + dirmask)
	fre := regexp.MustCompile("^" + filemask + "$")

	root = filepath.Clean(root)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}

		path = filepath.ToSlash(path)

		if info.IsDir() {
			if path == "." || len(path) <= len(root) {
				return nil
			}
			if !dre.MatchString(path + "/") {
				return filepath.SkipDir
			}
		}

		if fre.MatchString(path) {
			if relative && filepath.IsAbs(path) {
				path = path[len(root)+1:]
			}
			matches = append(matches, path)
		}
		return nil
	})
	return matches, nil
}

package controlsockets

import (
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/moul/advanced-ssh-config/pkg/utils"
)

func translateControlPath(input string) string {
	controlPath, err := utils.ExpandUser(input)
	if err != nil {
		return input
	}

	controlPath = strings.Replace(controlPath, "%h", "**/*", -1)

	for _, component := range []string{"%L", "%p", "%n", "%C", "%l", "%r"} {
		controlPath = strings.Replace(controlPath, component, "*", -1)
	}
	return controlPath
}

func LookupControlPathDir(controlPath string) ([]string, error) {
	return zglob.Glob(translateControlPath(controlPath))
}

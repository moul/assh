package controlsockets

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-zglob"
	"github.com/moul/advanced-ssh-config/pkg/utils"
)

type ControlSocket struct {
	path        string
	controlPath string
}

type ControlSockets []ControlSocket

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

func LookupControlPathDir(controlPath string) (ControlSockets, error) {
	controlPath = translateControlPath(controlPath)

	matches, err := zglob.Glob(controlPath)
	if err != nil {
		return nil, err
	}

	list := ControlSockets{}
	for _, socketPath := range matches {
		list = append(list, ControlSocket{
			path:        socketPath,
			controlPath: controlPath,
		})
	}
	return list, nil
}

func (s *ControlSocket) Path() string {
	return s.path
}

func (s *ControlSocket) RelativePath() string {
	idx := strings.Index(s.controlPath, "*")
	return s.path[idx:]
}

func (s *ControlSocket) CreatedAt() (time.Time, error) {
	return time.Now(), fmt.Errorf("not implemented")
}

func (s *ControlSocket) ActiveConnections() (int, error) {
	return -1, fmt.Errorf("not implemented")
}

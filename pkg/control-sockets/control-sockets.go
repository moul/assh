package controlsockets

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-zglob"
	"github.com/moul/advanced-ssh-config/pkg/utils"
)

// ControlSocket defines a unix-domain socket controlled by a master SSH process
type ControlSocket struct {
	path        string
	controlPath string
}

// ControlSockets is a list of ControlSocket
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

// LookupControlPathDir returns the ControlSockets in the ControlPath directory
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

// Path returns the absolute path of the socket
func (s *ControlSocket) Path() string {
	return s.path
}

// RelativePath returns a path relative to the configured ControlPath
func (s *ControlSocket) RelativePath() string {
	idx := strings.Index(s.controlPath, "*")
	return s.path[idx:]
}

// CreatedAt returns the modification time of the sock file
func (s *ControlSocket) CreatedAt() (time.Time, error) {
	stat, err := os.Stat(s.path)
	if err != nil {
		return time.Now(), err
	}

	return stat.ModTime(), nil
}

// ActiveConnections returns the amount of active connections using a control socket
func (s *ControlSocket) ActiveConnections() (int, error) {
	return -1, fmt.Errorf("not implemented")
}

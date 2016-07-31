package hooks

import (
	"bytes"
	"os"
	"os/exec"

	. "github.com/moul/advanced-ssh-config/pkg/logger"
	"github.com/moul/advanced-ssh-config/pkg/templates"
)

// DaemonDriver is a driver that daemons some texts to the terminal
type DaemonDriver struct {
	line string
}

// NewDaemonDriver returns a DaemonDriver instance
func NewDaemonDriver(line string) (DaemonDriver, error) {
	return DaemonDriver{
		line: line,
	}, nil
}

// Run daemons a line to the terminal
func (d DaemonDriver) Run(args RunArgs) error {
	var buff bytes.Buffer
	tmpl, err := templates.New(d.line + "\n")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buff, args); err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", buff.String())
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		cmd.Wait()
		Logger.Infof("daemon %q exited", d.line)
	}()

	return nil
}

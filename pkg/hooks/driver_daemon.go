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
	cmd  *exec.Cmd
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

	d.cmd = exec.Command("/bin/sh", "-c", buff.String())
	d.cmd.Stdout = os.Stderr
	d.cmd.Stderr = os.Stderr
	d.cmd.Stdin = os.Stdin
	if err := d.cmd.Start(); err != nil {
		return err
	}

	go func() {
		d.cmd.Wait()
		Logger.Infof("daemon %q exited", d.line)
	}()

	return nil
}

// Close closes a running command
func (d DaemonDriver) Close() error {
	if d.cmd == nil || d.cmd.Process == nil {
		return nil
	}

	err := d.cmd.Process.Kill()
	if err != nil {
		Logger.Warnf("daemon %q failed to stop: %v", d.line, err)
	}
	return err
}

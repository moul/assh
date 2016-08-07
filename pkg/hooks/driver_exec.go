package hooks

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/moul/advanced-ssh-config/pkg/templates"
)

// ExecDriver is a driver that execs some texts to the terminal
type ExecDriver struct {
	line string
}

// NewExecDriver returns a ExecDriver instance
func NewExecDriver(line string) (ExecDriver, error) {
	return ExecDriver{
		line: line,
	}, nil
}

// Run execs a line to the terminal
func (d ExecDriver) Run(args RunArgs) error {
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
	return cmd.Wait()
}

// Close is mandatory for the interface, here it does nothing
func (d ExecDriver) Close() error { return nil }

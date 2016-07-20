package hooks

import (
	"bytes"
	"os"
	"os/exec"
)
import "text/template"

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
	tmpl, err := template.New("exec").Parse(d.line + "\n")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buff, args); err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", buff.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

package hooks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

	var (
		availableShells = []string{"/bin/sh", "/bin/bash", "/bin/zsh"}
		selectedShell   = ""
	)
	for _, shell := range availableShells {
		info, err := os.Stat(shell)
		if err != nil {
			continue
		}
		if info.Mode()&0111 != 0 {
			selectedShell = shell
			break
		}
	}
	if selectedShell == "" {
		return fmt.Errorf("No available shell found. (tried %s)", strings.Join(availableShells, ", "))
	}

	cmd := exec.Command(selectedShell, "-c", buff.String()) // #nosec
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

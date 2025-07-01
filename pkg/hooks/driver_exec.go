package hooks

import (
	"bytes"
	"fmt"
	"moul.io/assh/v2/pkg/templates"
	"os"
	"os/exec"
	"runtime"
)

var (
	possibleShells = []string{
		"/bin/sh", "/bin/bash", "/bin/zsh",
		"/usr/bin/sh", "/usr/bin/bash", "/usr/bin/zsh",
		"/usr/local/bin/sh", "/usr/local/bin/bash", "/usr/local/bin/zsh",

		"C:\\Program Files\\Git\\bin\\bash.exe",
		"C:\\Program Files\\Git\\bin\\sh.exe",
		"C:\\Windows\\System32\\bash.exe",
		"C:\\Users\\paul.schroeder\\scoop\\shims\\sh.exe",
	}
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
	selectedShell := findAvailableShell(possibleShells)
	if selectedShell == "" {
		return fmt.Errorf("no available shell found. (tried %s)", possibleShells)
	}

	command, err := renderCommand(d.line, args)
	if err != nil {
		return fmt.Errorf("failed to render command: %w", err)
	}

	proc := exec.Command(selectedShell, "-c", command) // #nosec

	proc.Stdout = os.Stderr
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin

	if err = proc.Start(); err != nil {
		return err
	}

	return proc.Wait()
}

// Close is mandatory for the interface, here it does nothing
func (d ExecDriver) Close() error { return nil }

func renderCommand(line string, tmplArgs RunArgs) (string, error) {
	tmpl, err := templates.New(line + "\n")
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, tmplArgs); err != nil {
		return "", err
	}

	return buff.String(), nil
}

func findAvailableShell(shells []string) string {
	for _, shell := range shells {
		if isExecutable(shell) {
			return shell
		}
	}

	return ""
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// on Windows, the executable bit is not necessary
	if runtime.GOOS == "windows" {
		return true
	}

	// on Linux actually check if the file is executable
	return info.Mode()&0111 != 0
}

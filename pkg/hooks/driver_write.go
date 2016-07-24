package hooks

import (
	"os"

	"github.com/moul/advanced-ssh-config/pkg/templates"
)

// WriteDriver is a driver that writes some texts to the terminal
type WriteDriver struct {
	line string
}

// NewWriteDriver returns a WriteDriver instance
func NewWriteDriver(line string) (WriteDriver, error) {
	return WriteDriver{
		line: line,
	}, nil
}

// Run writes a line to the terminal
func (d WriteDriver) Run(args RunArgs) error {
	tmpl, err := templates.New(d.line + "\n")
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stderr, args)
}

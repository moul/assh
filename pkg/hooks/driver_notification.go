package hooks

import (
	"bytes"

	"github.com/haklop/gnotifier"
)
import "text/template"

// NotificationDriver is a driver that notifications some texts to the terminal
type NotificationDriver struct {
	line string
}

// NewNotificationDriver returns a NotificationDriver instance
func NewNotificationDriver(line string) (NotificationDriver, error) {
	return NotificationDriver{
		line: line,
	}, nil
}

// Run notifications a line to the terminal
func (d NotificationDriver) Run(args RunArgs) error {
	var buff bytes.Buffer
	tmpl, err := template.New("notification").Parse(d.line + "\n")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buff, args); err != nil {
		return err
	}

	notification := gnotifier.Notification("ASSH", buff.String())
	notification.GetConfig().Expiration = 3000
	notification.GetConfig().ApplicationName = "assh"

	return notification.Push()
}

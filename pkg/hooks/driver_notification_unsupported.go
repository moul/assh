// +build openbsd freebsd netbsd windows

package hooks

type NotificationDriver struct{}

func NewNotificationDriver(_ string) (NotificationDriver, error) { return NotificationDriver{}, nil }
func (NotificationDriver) Run(_ RunArgs) error                   { return nil }
func (d NotificationDriver) Close() error                        { return nil }

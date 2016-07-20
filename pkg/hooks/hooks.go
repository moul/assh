package hooks

import (
	"fmt"
	"strings"
)

// Hooks represents a slice of Hook
type Hooks []Hook

// Hook is a string
type Hook string

// HookDriver represents a hook driver
type HookDriver interface {
	Run(RunArgs) error
}

// RunArgs is a map of interface{}
type RunArgs map[string]interface{}

// InvokeAll calls all hooks
func (h *Hooks) InvokeAll(args RunArgs) error {
	drivers := []HookDriver{}

	for _, expr := range *h {
		driver, err := New(expr)
		if err != nil {
			return err
		}
		drivers = append(drivers, driver)
	}

	for _, driver := range drivers {
		if err := driver.Run(args); err != nil {
			return err
		}
	}
	return nil
}

// New returns an HookDriver instance
func New(expr Hook) (HookDriver, error) {
	driverName := strings.Split(string(expr), " ")[0]
	param := strings.Join(strings.Split(string(expr), " ")[1:], " ")
	switch driverName {
	case "write":
		driver, err := NewWriteDriver(param)
		return driver, err
	case "notify":
		driver, err := NewNotificationDriver(param)
		return driver, err
	case "exec":
		driver, err := NewExecDriver(param)
		return driver, err
	default:
		return nil, fmt.Errorf("No such driver %q", driverName)
	}
}

package gnotifier

import (
	"github.com/deckarep/gosx-notifier"
)

func (n *notifier) Push() error {
	err := n.IsValid()
	if err != nil {
		return err
	}

	notification := gosxnotifier.NewNotification(n.Config.Message)
	notification.Title = n.Config.Title
	notification.Sound = gosxnotifier.Default

	err = notification.Push()
	return err
}

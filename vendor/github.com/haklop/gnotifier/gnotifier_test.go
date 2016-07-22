package gnotifier

import (
	"testing"
)

func Test_Notification(t *testing.T) {
	n := Notification("Hey", "Hello")

	if n.GetConfig().Title != "Hey" {
		t.Error("NewNotification doesn't have a Title specified")
	}

	if n.GetConfig().Message != "Hello" {
		t.Error("NewNotification doesn't have a Message specified")
	}
}

func Test_Notification_Title_Validity(t *testing.T) {
	n := Notification("", "Hello")

	err := n.Push()
	if err == nil {
		t.Error("Notification should trigger an error, title is mandatory")
	}
}

func Test_Notification_Message_Validity(t *testing.T) {
	n := Notification("Title", "")

	err := n.Push()
	if err == nil {
		t.Error("Notification should trigger an error, message is mandatory")
	}
}

func Test_Builder_Types(t *testing.T) {
	var _ Builder = Notification
	var _ Builder = NullNotification
}

func Test_Records_A_Push(t *testing.T) {
	r := NewTestRecorder()
	var p Builder
	p = r.Notification

	n := p("title", "message")
	n.GetConfig().Expiration = 1000
	n.Push()

	if len(r.Pushed) != 1 {
		t.Fatal("Expected one message")
	}
}

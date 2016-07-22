package gnotifier

import (
	"errors"
)

// GNotifier interface
type GNotifier interface {
	Push() error
	GetConfig() *Config
	IsValid() error
}

// Config define the notification options
type Config struct {
	Title           string
	Message         string
	Expiration      int32
	ApplicationName string
}

func (c *Config) IsValid() error {
	if c.Title == "" {
		return errors.New("A Title is mandatory")
	}
	if c.Message == "" {
		return errors.New("A Message is mandatory")
	}
	return nil
}

// Builder abstracts the concrete function Notification
type Builder func(title, message string) GNotifier

type notifier struct {
	Config *Config
}

func (n *notifier) GetConfig() *Config {
	return n.Config
}

func (n *notifier) IsValid() error {
	return n.GetConfig().IsValid()
}

// Notification is the builder
func Notification(title, message string) GNotifier {
	config := &Config{title, message, 5000, ""}
	n := &notifier{Config: config}
	return n
}

type nullNotifier struct {
	Config *Config
}

func (n *nullNotifier) GetConfig() *Config {
	return n.Config
}

func (n *nullNotifier) IsValid() error {
	return n.GetConfig().IsValid()
}

func (n *nullNotifier) Push() error {
	return nil
}

// NullNotification is the builder for tests where no side effects are desired
func NullNotification(title, message string) GNotifier {
	config := &Config{title, message, 5000, ""}
	n := &nullNotifier{Config: config}
	return n
}

type recordingNotifier struct {
	Config   *Config
	Recorder *TestRecorder
}

func (n *recordingNotifier) GetConfig() *Config {
	return n.Config
}

func (n *recordingNotifier) IsValid() error {
	return n.GetConfig().IsValid()
}

func (n *recordingNotifier) Push() error {
	n.Recorder.Pushed = append(n.Recorder.Pushed, n.GetConfig())
	return nil
}

// TestRecorder provides a way to verify the GNotifier api use
// (not intended for production code)
type TestRecorder struct {
	Pushed []*Config
}

// NewTestRecorder constructs a new TestRecorder. Use its
// Notification method as test Builder.
func NewTestRecorder() *TestRecorder {
	return &TestRecorder{[]*Config{}}
}

func (r *TestRecorder) Notification(title, message string) GNotifier {
	config := &Config{title, message, 5000, ""}
	n := &recordingNotifier{Config: config, Recorder: r}
	return n
}

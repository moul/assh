package config

import "encoding/json"

// HostHooks represents a static list of Hooks
type HostHooks struct {
	OnConnect    Hooks `yaml:"onconnect,omitempty,flow" json:"OnConnect,omitempty"`
	OnDisconnect Hooks `yaml:"ondisconnect,omitempty,flow" json:"OnDisconnect,omitempty"`
}

// Length returns the quantity of hooks of any type
func (hh *HostHooks) Length() int {
	if hh == nil {
		return 0
	}
	return len(hh.OnConnect) + len(hh.OnDisconnect)
}

// String returns the JSON output
func (hh *HostHooks) String() string {
	s, _ := json.Marshal(hh)
	return string(s)
}

// Hooks represents a slice of Hook
type Hooks []Hook

// Hook is a string
type Hook string

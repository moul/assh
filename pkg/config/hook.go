package config

import "encoding/json"

type HostHooks struct {
	OnConnect    Hooks `yaml:"onconnect,omitempty,flow" json:"OnConnect,omitempty"`
	OnDisconnect Hooks `yaml:"ondisconnect,omitempty,flow" json:"OnDisconnect,omitempty"`
}

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

type Hooks []Hook

type Hook string

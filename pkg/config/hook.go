package config

import (
	"encoding/json"

	"github.com/moul/advanced-ssh-config/pkg/hooks"
)

// HostHooks represents a static list of Hooks
type HostHooks struct {
	BeforeConnect  hooks.Hooks `yaml:"beforeconnect,omitempty,flow" json:"BeforeConnect,omitempty"`
	OnConnect      hooks.Hooks `yaml:"onconnect,omitempty,flow" json:"OnConnect,omitempty"`
	OnDisconnect   hooks.Hooks `yaml:"ondisconnect,omitempty,flow" json:"OnDisconnect,omitempty"`
	OnConnectError hooks.Hooks `yaml:"onconnecterror,omitempty,flow" json:"OnConnectError,omitempty"`
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

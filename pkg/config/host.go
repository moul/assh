package config

import (
	"fmt"
	"io"
)

// Host defines the configuration flags of a host
type Host struct {
	// ssh-config fields
	Hostname    string `yaml:"Hostname,omitempty,flow" json:"Hostname,omitempty"`
	User        string `yaml:"User,omitempty,flow" json:"User,omitempty"`
	Port        uint   `yaml:"Port,omitempty,flow" json:"Port,omitempty"`
	ControlPath string `yaml:"Controlpath,omitempty,flow" json:"Controlpath,omitempty"`

	// ssh-config fields with a different behavior
	ProxyCommand string `yaml:"ProxyCommand,omitempty,flow" json:"ProxyCommand,omitempty"`

	// assh fields
	name               string   `yaml:- json:-`
	Gateways           []string `yaml:"Gateways,omitempty,flow" json:"Gateways,omitempty"`
	ResolveNameservers []string `yaml:"ResolveNameservers,omitempty,flow" json:"ResolveNameservers,omitempty"`
	ResolveCommand     string   `yaml:"ResolveCommand,omitempty,flow" json:"ResolveCommand,omitempty"`
}

func (h *Host) Name() string {
	return h.name
}

// ApplyDefaults ensures a Host is valid by filling the missing fields with defaults
func (h *Host) ApplyDefaults(defaults *Host) {
	if h.Hostname == "" {
		h.Hostname = defaults.Hostname
	}
	if h.User == "" {
		h.User = defaults.User
	}
	if h.Port == 0 {
		h.Port = defaults.Port
	}
	if len(h.Gateways) == 0 {
		h.Gateways = defaults.Gateways
	}
	if h.ProxyCommand == "" {
		h.ProxyCommand = defaults.ProxyCommand
	}
	if len(h.ResolveNameservers) == 0 {
		h.ResolveNameservers = defaults.ResolveNameservers
	}
	if h.ResolveCommand == "" {
		h.ResolveCommand = defaults.ResolveCommand
	}
	if h.ControlPath == "" {
		h.ControlPath = defaults.ControlPath
	}

	// Extra defaults
	if h.Port == 0 {
		h.Port = 22
	}
}

func (h *Host) WriteSshConfigTo(w io.Writer) error {
	fmt.Fprintf(w, "Host %s\n", h.Name())
	if h.Hostname != "" {
		fmt.Fprintf(w, "  HostName %s\n", h.Hostname)
	}
	if h.Port != 0 {
		fmt.Fprintf(w, "  Port %d\n", h.Port)
	}
	if h.User != "" {
		fmt.Fprintf(w, "  User %s\n", h.User)
	}
	if h.ControlPath != "" {
		fmt.Fprintf(w, "  ControlPath %s\n", h.ControlPath)
	}
	fmt.Fprintln(w)
	return nil
}

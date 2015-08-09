package config

// Host defines the configuration flags of a host
type Host struct {
	Name     string   `yaml:-`
	Host     string   `yaml:"host,omitempty,flow" json:"host,omitempty"`
	User     string   `yaml:"user,omitempty,flow" json:"user,omitempty"`
	Port     uint     `yaml:"port,omitempty,flow" json:"port,omitempty"`
	Gateways []string `yaml:"gateways,omitempty,flow" json:"gateways,omitempty"`
}

// ApplyDefaults ensures a Host is valid by filling the missing fields with defaults
func (h *Host) ApplyDefaults(defaults Host) {
	if h.Host == "" {
		h.Host = defaults.Host
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

	// Extra defaults
	if h.Port == 0 {
		h.Port = 22
	}
}

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/moul/advanced-ssh-config/vendor/gopkg.in/yaml.v2"
)

// Config contains a list of Hosts sections and a Defaults section representing a configuration file
type Config struct {
	Hosts    map[string]Host `json:"hosts"`
	Defaults Host            `json:"defaults,omitempty"`
}

// JsonString returns a string representing the JSON of a Config object
func (c *Config) JsonString() error {
	output, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s\n", output)
	return nil
}

func (c *Config) getHostByName(name string, safe bool) (*Host, error) {
	if host, ok := c.Hosts[name]; ok {
		var computedHost Host = host
		computedHost.ApplyDefaults(c.Defaults)
		computedHost.Name = name
		return &computedHost, nil
	}

	for pattern, host := range c.Hosts {
		matched, err := path.Match(pattern, name)
		if err != nil {
			return nil, err
		}
		if matched {
			var computedHost Host = host
			computedHost.ApplyDefaults(c.Defaults)
			computedHost.Name = name
			return &computedHost, nil
		}
	}

	if safe {
		host := &Host{
			Host: name,
			Name: name,
		}
		host.ApplyDefaults(c.Defaults)
		return host, nil
	}

	return nil, fmt.Errorf("no such host: %s", name)
}

func (c *Config) getHostByPath(path string, safe bool) (*Host, error) {
	parts := strings.SplitN(path, "/", 2)

	host, err := c.getHostByName(parts[0], safe)
	if err != nil {
		return nil, err
	}

	if len(parts) > 1 {
		host.Gateways = []string{parts[1]}
	}

	return host, nil
}

// GetGatewaySafe returns gateway Host configuration, a gateway is like a Host, except, the host path is not resolved
func (c *Config) GetGatewaySafe(name string) *Host {
	host, err := c.getHostByName(name, true)
	if err != nil {
		panic(err)
	}
	return host
}

// GetHost returns a matching host form Config hosts list
func (c *Config) GetHost(name string) (*Host, error) {
	return c.getHostByPath(name, false)
}

// GetHostSafe won't fail, in case the host is not found, it will returns a virtual host matching the pattern
func (c *Config) GetHostSafe(name string) *Host {
	host, err := c.getHostByPath(name, true)
	if err != nil {
		panic(err)
	}
	return host
}

// LoadFile loads the content of a configuration file in the Config object
func (c *Config) LoadFile(filename string) error {
	filepath, err := expandUser(filename)
	if err != nil {
		return err
	}

	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return err
	}
	return nil
}

// New returns an instantiated Config object
func New() *Config {
	var config Config
	config.Hosts = make(map[string]Host)
	return &config
}

// Open returns a Config object loaded with standard configuration file paths
func Open() (*Config, error) {
	config := New()
	err := config.LoadFile("~/.ssh/assh.yml")
	if err != nil {
		return nil, err
	}
	return config, nil
}

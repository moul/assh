package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/moul/advanced-ssh-config/vendor/gopkg.in/yaml.v2"
)

// Config contains a list of Hosts sections and a Defaults section representing a configuration file
type Config struct {
	Hosts    map[string]Host `json:"hosts"`
	Defaults Host            `json:"defaults",omitempty`
	Includes []string        `json:"includes",omitempty`

	includedFiles map[string]bool
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
		computedHost.ApplyDefaults(&c.Defaults)
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
			computedHost.ApplyDefaults(&c.Defaults)
			computedHost.Name = name
			return &computedHost, nil
		}
	}

	if safe {
		host := &Host{
			Host: name,
			Name: name,
		}
		host.ApplyDefaults(&c.Defaults)
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

// LoadConfig loads the content of an io.Reader source
func (c *Config) LoadConfig(source io.Reader) error {
	buf, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, &c)
}

// LoadFile loads the content of a configuration file in the Config object
func (c *Config) LoadFile(filename string) error {
	// Resolve '~' and '$HOME'
	filepath, err := expandUser(filename)
	if err != nil {
		return err
	}
	logrus.Debugf("Loading config file '%s'", filepath)

	// Anti-loop protection
	if _, ok := c.includedFiles[filepath]; ok {
		logrus.Debugf("File %s already loaded", filepath)
		return nil
	}
	c.includedFiles[filepath] = false

	// Read file
	source, err := os.Open(filepath)
	if err != nil {
		return err
	}

	// Load config stream
	err = c.LoadConfig(source)
	if err != nil {
		return err
	}

	// Successful loading
	c.includedFiles[filepath] = true

	// Handling includes
	for _, include := range c.Includes {
		c.LoadFiles(include)
	}

	return nil
}

// Loadfiles will try to glob the pattern and load each maching entries
func (c *Config) LoadFiles(pattern string) error {
	// Resolve '~' and '$HOME'
	expandedPattern, err := expandUser(pattern)
	if err != nil {
		return err
	}

	// Globbing
	filepaths, err := filepath.Glob(expandedPattern)
	if err != nil {
		return err
	}

	// Load files iteratively
	for _, filepath := range filepaths {
		err := c.LoadFile(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

// New returns an instantiated Config object
func New() *Config {
	var config Config
	config.Hosts = make(map[string]Host)
	config.includedFiles = make(map[string]bool)
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

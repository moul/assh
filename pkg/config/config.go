package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/moul/advanced-ssh-config/vendor/gopkg.in/yaml.v2"

	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

// Config contains a list of Hosts sections and a Defaults section representing a configuration file
type Config struct {
	Hosts    map[string]Host `json:"hosts"`
	Defaults Host            `json:"defaults",omitempty`
	Includes []string        `json:"includes",omitempty`

	includedFiles map[string]bool
}

// JsonString returns a string representing the JSON of a Config object
func (c *Config) JsonString() ([]byte, error) {
	output, err := json.MarshalIndent(c, "", "  ")
	return output, err
}

// computeHost applies defaults, inherits hosts and configure internal fields
func computeHost(host *Host, config *Config, name string, fullCompute bool) (*Host, error) {
	var computedHost Host
	if host != nil {
		computedHost = *host
	}

	// name internal field
	computedHost.name = name
	computedHost.inherited = make(map[string]bool, 0)
	// self is already inherited
	computedHost.inherited[name] = true

	if fullCompute {
		// apply defaults based on "Host *"
		computedHost.ApplyDefaults(&config.Defaults)

		// Inheritance
		// FIXME: allow deeper inheritance:
		//     currently not resolving inherited hosts
		//     we should resolve all inherited hosts and pass the
		//     currently resolved hosts to avoid computing an host twice
		for _, name := range host.Inherits {
			_, found := computedHost.inherited[name]
			if found {
				Logger.Debugf("Detected circular loop inheritance, skiping...")
				continue
			}
			computedHost.inherited[name] = true

			target, err := config.getHostByPath(name, false, false)
			if err != nil {
				Logger.Warnf("Cannot inherits from %q: %v", name, err)
				continue
			}
			computedHost.ApplyDefaults(target)
		}
	}

	return &computedHost, nil
}

func (c *Config) getHostByName(name string, safe bool, compute bool) (*Host, error) {
	if host, ok := c.Hosts[name]; ok {
		return computeHost(&host, c, name, compute)
	}

	for pattern, host := range c.Hosts {
		matched, err := path.Match(pattern, name)
		if err != nil {
			return nil, err
		}
		if matched {
			return computeHost(&host, c, name, compute)
		}
	}

	if safe {
		host := Host{
			HostName: name,
		}
		return computeHost(&host, c, name, compute)
	}

	return nil, fmt.Errorf("no such host: %s", name)
}

func (c *Config) getHostByPath(path string, safe bool, compute bool) (*Host, error) {
	parts := strings.SplitN(path, "/", 2)

	host, err := c.getHostByName(parts[0], safe, compute)
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
	host, err := c.getHostByName(name, true, true) // FIXME: fullCompute for gateway ?
	if err != nil {
		panic(err)
	}
	return host
}

// GetHost returns a matching host form Config hosts list
func (c *Config) GetHost(name string) (*Host, error) {
	return c.getHostByPath(name, false, true)
}

// GetHostSafe won't fail, in case the host is not found, it will returns a virtual host matching the pattern
func (c *Config) GetHostSafe(name string) *Host {
	host, err := c.getHostByPath(name, true, true)
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
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return err
	}
	c.applyMissingNames()
	return nil
}

func (c *Config) applyMissingNames() {
	for key, _ := range c.Hosts {
		host := c.Hosts[key]
		host.name = key
		c.Hosts[key] = host
	}
	c.Defaults.isDefault = true
}

// LoadFile loads the content of a configuration file in the Config object
func (c *Config) LoadFile(filename string) error {
	// Resolve '~' and '$HOME'
	filepath, err := expandUser(filename)
	if err != nil {
		return err
	}
	Logger.Debugf("Loading config file '%s'", filepath)

	// Anti-loop protection
	if _, ok := c.includedFiles[filepath]; ok {
		Logger.Debugf("File %s already loaded", filepath)
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
		err = c.LoadFiles(include)
		if err != nil {
			return err
		}
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

// sortedNames returns the host names sorted alphabetically
func (c *Config) sortedNames() []string {
	names := sort.StringSlice{}
	for key, _ := range c.Hosts {
		names = append(names, key)
	}
	sort.Sort(names)
	return names
}

// ExportSshConfig returns a .ssh/config valid file containing assh configuration
func (c *Config) WriteSshConfigTo(w io.Writer) error {
	fmt.Fprintln(w, "# ssh config generated by advanced-ssh-config")
	// FIXME: add version
	fmt.Fprintln(w)

	fmt.Fprintln(w, "# host-based configuration")
	for _, name := range c.sortedNames() {
		host := c.Hosts[name]
		host.WriteSshConfigTo(w)
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "# global configuration")
	c.Defaults.name = "*"
	c.Defaults.WriteSshConfigTo(w)

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

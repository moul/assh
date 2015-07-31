package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/moul/advanced-ssh-config/vendor/gopkg.in/yaml.v2"

	"github.com/moul/advanced-ssh-config/utils"
)

type Config struct {
	Hosts    map[string]Host `json:"hosts"`
	Defaults Host            `json:"defaults,omitempty"`
}

func (c *Config) PrettyPrint() error {
	output, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s\n", output)
	return nil
}

func LoadFile(filename string) (*Config, error) {
	var config Config

	filepath, err := utils.ExpandUser(filename)
	if err != nil {
		return nil, err
	}

	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return nil, err
	}

	// config.PrettyPrint()
	// fmt.Printf("Config: %v\n", config)

	return &config, nil
}

func Open() (*Config, error) {
	return LoadFile("~/.ssh/assh.yml")
}

func (c *Config) GetHost(name string) (*Host, error) {
	if host, ok := c.Hosts[name]; ok {
		var computedHost Host = host
		computedHost.ApplyDefaults(c.Defaults)
		return &computedHost, nil
	}
	return nil, fmt.Errorf("no such host: %s", name)
}

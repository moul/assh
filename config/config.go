package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/moul/advanced-ssh-config/utils"
)

type Host struct {
	Host string `yaml:"host,omitempty,flow" json:"host,omitempty"`
	User string `yaml:"user,omitempty,flow" json:"user,omitempty"`
	Port string `yaml:"port,omitempty,flow" json:"port,omitempty"`
}

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

func Open() (*Config, error) {
	filename := "~/.ssh/assh.yml"
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

	config.PrettyPrint()
	fmt.Printf("Config: %v\n", config)

	return &config, nil
}

package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/moul/advanced-ssh-config/utils"
)

type Host struct {
	Host string `yaml:"host,omitempty,flow"`
	User string `yaml:"user,omitempty,flow"`
	Port string `yaml:"port,omitempty,flow"`
}

type Config struct {
	Hosts    map[string]Host
	Defaults Host
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

	fmt.Printf("Config: %v\n", config)

	return &config, nil
}

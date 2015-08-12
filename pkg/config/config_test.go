package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummyConfig() *Config {
	config := New()
	config.Hosts["toto"] = Host{
		Host: "1.2.3.4",
	}
	config.Hosts["titi"] = Host{
		Host: "tata",
		Port: 23,
		User: "moul",
	}
	config.Defaults = Host{
		Port: 22,
		User: "root",
	}
	return config
}

func TestNew(t *testing.T) {
	config := New()

	assert.Equal(t, len(config.Hosts), 0)
	assert.Equal(t, config.Defaults.Port, uint(0))
	assert.Equal(t, config.Defaults.Host, "")
	assert.Equal(t, config.Defaults.User, "")
}

func TestConfig(t *testing.T) {
	config := dummyConfig()

	assert.Equal(t, len(config.Hosts), 2)
	assert.Equal(t, config.Hosts["toto"].Host, "1.2.3.4")
	assert.Equal(t, config.Defaults.Port, uint(22))
}

func TestConfigLoadFile(t *testing.T) {
	config := New()
	err := config.LoadConfig(strings.NewReader(`
hosts:
  aaa:
    host: 1.2.3.4
  bbb:
    port: 21
  ccc:
    host: 5.6.7.8
    port: 24
    user: toor
defaults:
  port: 22
  user: root
`))
	assert.Nil(t, err)
	assert.Equal(t, len(config.Hosts), 3)
	assert.Equal(t, config.Hosts["aaa"].Host, "1.2.3.4")
	assert.Equal(t, config.Hosts["aaa"].Port, uint(0))
	assert.Equal(t, config.Hosts["aaa"].User, "")
	assert.Equal(t, config.Hosts["bbb"].Host, "")
	assert.Equal(t, config.Hosts["bbb"].Port, uint(21))
	assert.Equal(t, config.Hosts["bbb"].User, "")
	assert.Equal(t, config.Hosts["ccc"].Host, "5.6.7.8")
	assert.Equal(t, config.Hosts["ccc"].Port, uint(24))
	assert.Equal(t, config.Hosts["ccc"].User, "toor")
	assert.Equal(t, config.Defaults.Port, uint(22))
	assert.Equal(t, config.Defaults.User, "root")
}

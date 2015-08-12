package config

import (
	"testing"

	"github.com/moul/advanced-ssh-config/vendor/github.com/stretchr/testify/assert"
)

func TestHostApplyDefaults(t *testing.T) {
	host := &Host{
		Name: "example",
		Host: "example.com",
		User: "root",
	}
	defaults := &Host{
		User: "bobby",
		Port: 42,
	}
	host.ApplyDefaults(defaults)
	assert.Equal(t, host.Port, uint(42))
	assert.Equal(t, host.Name, "example")
	assert.Equal(t, host.Host, "example.com")
	assert.Equal(t, host.User, "root")
	assert.Equal(t, len(host.Gateways), 0)
	assert.Equal(t, host.ProxyCommand, "")
	assert.Equal(t, len(host.ResolveNameservers), 0)
	assert.Equal(t, host.ResolveCommand, "")
	assert.Equal(t, host.ControlPath, "")
}

func TestHostApplyDefaults_empty(t *testing.T) {
	host := &Host{}
	defaults := &Host{}
	host.ApplyDefaults(defaults)
	assert.Equal(t, host.Port, uint(22))
	assert.Equal(t, host.Name, "")
	assert.Equal(t, host.Host, "")
	assert.Equal(t, host.User, "")
	assert.Equal(t, len(host.Gateways), 0)
	assert.Equal(t, host.ProxyCommand, "")
	assert.Equal(t, len(host.ResolveNameservers), 0)
	assert.Equal(t, host.ResolveCommand, "")
	assert.Equal(t, host.ControlPath, "")
}

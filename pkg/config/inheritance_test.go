package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCyclicInheritanceKeepsOtherDefaults(t *testing.T) {
	const source = `
hosts:
  target:
    aliases: target-alias
    inherits: [shared, transport]
    user: alice
templates:
  shared:
    inherits: target-alias/gateway
    port: 2200
  transport:
    inherits: common
  common:
    compression: "yes"
`
	cfg := New()
	require.NoError(t, cfg.LoadConfig(strings.NewReader(source)))

	host, err := cfg.GetHost("target")
	require.NoError(t, err)
	require.Equal(t, "alice", host.User)
	require.Equal(t, "2200", host.Port)
	require.Equal(t, "yes", host.Compression)
	require.Empty(t, host.Gateways)

	var output strings.Builder
	require.NoError(t, cfg.WriteSSHConfigTo(&output))
	require.Contains(t, output.String(), "Host target\n")
	require.Contains(t, output.String(), "  Port 2200\n")
}

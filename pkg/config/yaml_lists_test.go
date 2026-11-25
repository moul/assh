package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/assh/v2/pkg/hooks"
	"moul.io/assh/v2/pkg/utils"
)

func TestYAMLScalarsAndScalarOrStringSlices(t *testing.T) {
	const source = `
defaults:
  hooks:
    beforeconnect: exec echo ready
    ondisconnect:
      - write closed
hosts:
  example:
    aliases: [first, second]
    ciphers: value1,value2
    identityfile: "path with spaces"
    sendenv: []
`
	cfg := New()
	require.NoError(t, cfg.LoadConfig(strings.NewReader(source)))
	host := cfg.Hosts["example"]
	require.Equal(t, utils.ScalarOrStringSlice{"first", "second"}, host.Aliases)
	require.Equal(t, utils.ScalarOrStringSlice{"value1,value2"}, host.Ciphers)
	require.Equal(t, utils.ScalarOrStringSlice{"path with spaces"}, host.IdentityFile)
	require.Nil(t, host.SendEnv)
	require.Equal(t, hooks.Hooks{"exec echo ready"}, cfg.Defaults.Hooks.BeforeConnect)
	require.Equal(t, hooks.Hooks{"write closed"}, cfg.Defaults.Hooks.OnDisconnect)

	invalid := "hosts:\n  invalid:\n    aliases: [valid, 1]\n"
	require.ErrorContains(t, New().LoadConfig(strings.NewReader(invalid)), "expected a string")
}

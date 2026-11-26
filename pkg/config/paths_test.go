package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigInHomeDirectoryWithSpaces(t *testing.T) {
	home := filepath.Join(t.TempDir(), "home with spaces")
	t.Setenv("HOME", home)
	sshDir := filepath.Join(home, ".ssh")
	require.NoError(t, os.MkdirAll(sshDir, 0o700))

	input := []byte("includes:\n  - ~/.ssh/*.inc\nhosts:\n  root-host: {}\n")
	require.NoError(t, os.WriteFile(filepath.Join(sshDir, "assh.yml"), input, 0o600))
	include := []byte("hosts:\n  included: {}\n")
	require.NoError(t, os.WriteFile(filepath.Join(sshDir, "hosts with spaces.inc"), include, 0o600))

	cfg, err := Open("~/.ssh/assh.yml")
	require.NoError(t, err)
	require.Contains(t, cfg.Hosts, "root-host")
	require.Contains(t, cfg.Hosts, "included")
	require.NoError(t, cfg.SaveSSHConfig())
	require.FileExists(t, filepath.Join(sshDir, "config"))

	cfg.SaveNewKnownHost("included")
	require.NoError(t, cfg.KnownHostsFileExists())
	require.NoError(t, cfg.LoadKnownHosts())
	knownHosts, err := os.ReadFile(filepath.Join(sshDir, "assh_known_hosts"))
	require.NoError(t, err)
	require.Equal(t, "included\n", string(knownHosts))
}

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyIncludesKeepHostsAndTemplates(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ASSH_TEST_INCLUDES", dir)
	for name, content := range map[string]string{
		"01-first.yml": "hosts:\n  first:\n    user: alice\n",
		"02-empty.yml": "hosts:\ntemplates: null\n",
		"03-null.yml":  "null\n",
		"04-last.yml":  "hosts:\n  last:\n    user: bob\n",
	} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	cfg := New()
	require.NoError(t, cfg.LoadConfig(strings.NewReader("templates:\n  common:\n    port: 2200\n")))
	require.NoError(t, cfg.LoadFiles("$ASSH_TEST_INCLUDES/*.yml"))

	require.Len(t, cfg.Hosts, 2)
	require.Equal(t, "alice", cfg.Hosts["first"].User)
	require.Equal(t, "bob", cfg.Hosts["last"].User)
	require.Equal(t, "2200", cfg.Templates["common"].Port)
	require.Len(t, cfg.IncludedFiles(), 4)
}

func TestLoadFileCanRetryAfterReadFailure(t *testing.T) {
	file := filepath.Join(t.TempDir(), "retry.yml")
	cfg := New()
	require.ErrorIs(t, cfg.LoadFile(file), os.ErrNotExist)
	require.Empty(t, cfg.IncludedFiles())

	require.NoError(t, os.WriteFile(file, []byte("hosts:\n  recovered: {}\n"), 0o600))
	require.NoError(t, cfg.LoadFile(file))
	require.Contains(t, cfg.Hosts, "recovered")
	require.Equal(t, []string{file}, cfg.IncludedFiles())
}

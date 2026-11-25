package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/assh/v2/pkg/config"
	"moul.io/assh/v2/pkg/hooks"
)

func TestRewriteSSHConfigRunsHooksInOrder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test requires a POSIX shell")
	}

	dir := t.TempDir()
	t.Setenv("HOME", dir)
	backup := filepath.Join(dir, "backup")
	t.Setenv("ASSH_TEST_BACKUP", backup)
	sshDir := filepath.Join(dir, ".ssh")
	require.NoError(t, os.Mkdir(sshDir, 0o700))
	destination := filepath.Join(sshDir, "config")
	const original = "# previous config\n"
	require.NoError(t, os.WriteFile(destination, []byte(original), 0o600))

	conf := config.New()
	require.NoError(t, conf.LoadConfig(strings.NewReader("hosts:\n  example:\n    user: alice\n")))
	conf.Defaults.Hooks.BeforeConfigWrite = hooks.Hooks{"exec cat {{.SSHConfigPath}} > \"$ASSH_TEST_BACKUP\""}
	conf.Defaults.Hooks.AfterConfigWrite = hooks.Hooks{"exec echo '# after' >> {{.SSHConfigPath}}"}

	t.Run("failed write", func(t *testing.T) {
		conf.Hosts["example"].ControlMaster = "invalid"
		require.Error(t, rewriteSSHConfig(conf))
		contents, err := os.ReadFile(destination)
		require.NoError(t, err)
		require.Equal(t, original, string(contents))
	})
	t.Run("successful write", func(t *testing.T) {
		conf.Hosts["example"].ControlMaster = ""
		require.NoError(t, rewriteSSHConfig(conf))
		saved, err := os.ReadFile(backup)
		require.NoError(t, err)
		require.Equal(t, original, string(saved))
		contents, err := os.ReadFile(destination)
		require.NoError(t, err)
		require.Contains(t, string(contents), "Host example")
		require.True(t, strings.HasSuffix(string(contents), "# after\n"))
	})
}

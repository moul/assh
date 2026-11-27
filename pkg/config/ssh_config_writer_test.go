package config

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingConfigOutput struct {
	remaining int
	err       error
}

func (w *failingConfigOutput) Write(p []byte) (int, error) {
	if len(p) > w.remaining {
		return w.remaining, w.err
	}
	w.remaining -= len(p)
	return len(p), nil
}

func TestWriteSSHConfigReportsOutputErrors(t *testing.T) {
	cfg := New()
	require.NoError(t, cfg.LoadConfig(strings.NewReader("hosts:\n  example:\n    user: test\n")))

	t.Run("configuration header", func(t *testing.T) {
		w := &failingConfigOutput{err: io.ErrClosedPipe}
		require.ErrorIs(t, cfg.WriteSSHConfigTo(w), io.ErrClosedPipe)
	})
	t.Run("host entry", func(t *testing.T) {
		w := &failingConfigOutput{remaining: len("Host example\n"), err: io.ErrClosedPipe}
		require.ErrorIs(t, cfg.Hosts["example"].WriteSSHConfigTo(w), io.ErrClosedPipe)
	})
	t.Run("short write", func(t *testing.T) {
		w := &failingConfigOutput{remaining: 1}
		require.ErrorIs(t, cfg.WriteSSHConfigTo(w), io.ErrShortWrite)
	})
}

func TestSaveSSHConfigRemovesTemporaryFileAfterRenameFailure(t *testing.T) {
	dir := t.TempDir()
	cfg := New()
	cfg.sshConfigPath = filepath.Join(dir, "destination")
	require.NoError(t, os.Mkdir(cfg.sshConfigPath, 0o700))

	require.Error(t, cfg.SaveSSHConfig())

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "destination", entries[0].Name())
	require.True(t, entries[0].IsDir())
}

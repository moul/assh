package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWildcardDefaultsUseTheWholeHostName(t *testing.T) {
	const source = `
hosts:
  "prod-*":
    user: prefix
  "*-backup":
    port: 2022
  "pre-*mid*end":
    user: multiple
  "*.example":
    user: broad
  "prod-*.example":
    user: narrow
  prod-db-backup: {}
  other-prod-db: {}
  main-backup-dev: {}
  pre-first-mid-last-end: {}
  mid-pre-first-end: {}
  prod-db.example: {}
  prod-explicit.example:
    user: explicit
`
	cfg := New()
	require.NoError(t, cfg.LoadConfig(strings.NewReader(source)))
	require.Equal(t, "prefix", cfg.Hosts["prod-db-backup"].User)
	require.Equal(t, "2022", cfg.Hosts["prod-db-backup"].Port)
	require.Empty(t, cfg.Hosts["other-prod-db"].User)
	require.Empty(t, cfg.Hosts["main-backup-dev"].Port)
	require.Equal(t, "multiple", cfg.Hosts["pre-first-mid-last-end"].User)
	require.Empty(t, cfg.Hosts["mid-pre-first-end"].User)
	require.Equal(t, "broad", cfg.Hosts["prod-db.example"].User)
	require.Equal(t, "explicit", cfg.Hosts["prod-explicit.example"].User)
}

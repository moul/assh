package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSSHDestination(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{name: "plain", args: []string{"host"}, want: "host"},
		{name: "grouped flags", args: []string{"-vvp", "2200", "host"}, want: "host"},
		{name: "attached values", args: []string{"-p2200", "-Fcustom config", "host", "ls", "-l"}, want: "host"},
		{name: "jump host", args: []string{"-J", "jump", "user@host"}, want: "user@host"},
		{name: "native separator", args: []string{"-vv", "--", "host", "-V"}, want: "host"},
		{name: "version", args: []string{"-V"}},
		{name: "query", args: []string{"-Q", "cipher"}},
		{name: "missing value", args: []string{"-F"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, sshDestination(tc.args))
		})
	}
}

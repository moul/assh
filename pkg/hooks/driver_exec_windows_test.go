//go:build windows
// +build windows

package hooks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_isExecutable(t *testing.T) {
	type args struct {
		path string
		goos string
	}
	type want struct {
		executable bool
	}

	tt := map[string]struct {
		args args
		want want
	}{
		"ErrFileDoesNotExist": {
			args: args{
				path: "testdata/not_exit",
			},
			want: want{
				executable: false,
			},
		},
		"Ok": {
			args: args{
				path: "testdata/win.shell",
			},
			want: want{
				executable: true,
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := isExecutable(tc.args.path)
			require.Equal(t, tc.want.executable, got)
		})
	}
}

func Test_findAvailableShell(t *testing.T) {
	type args struct {
		shells []string
	}
	type want struct {
		out string
	}

	tt := map[string]struct {
		args args
		want want
	}{
		"Nil": {
			args: args{
				shells: nil,
			},
			want: want{
				out: "",
			},
		},
		"Empty": {
			args: args{
				shells: []string{},
			},
			want: want{
				out: "",
			},
		},
		"NoShells": {
			args: args{
				shells: []string{
					"not existing shell",
					"another not existing shell",
				},
			},
			want: want{
				out: "",
			},
		},
		"Ok": {
			args: args{
				shells: []string{
					"not existing shell",
					"another not existing shell",
					"testdata/win.shell",
				},
			},
			want: want{
				out: "testdata/win.shell",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := findAvailableShell(tc.args.shells)
			require.Equal(t, tc.want.out, got)
		})
	}
}

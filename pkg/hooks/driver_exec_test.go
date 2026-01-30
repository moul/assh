package hooks

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func EqualError(tb testing.TB, want, got error) {
	tb.Helper()

	if want == nil || got == nil {
		require.Equal(tb, want, got)
		return
	}

	require.EqualError(tb, got, want.Error())
}

func Test_renderCommand(t *testing.T) {
	type args struct {
		line string
		args RunArgs
	}
	type want struct {
		out string
		err error
	}

	tt := map[string]struct {
		args args
		want want
	}{
		"Ok": {
			args: args{
				line: "echo {{.foo}}",
				args: map[string]interface{}{
					"foo": "bar",
				},
			},
			want: want{
				out: "echo bar\n",
				err: nil,
			},
		},
		"ErrNewTemplate": {
			args: args{
				line: "{{ not_a_function",
				args: nil,
			},
			want: want{
				out: "",
				err: errors.Errorf("template: :1: function %q not defined", "not_a_function"),
			},
		},
		"ErrExecute": {
			args: args{
				line: "echo {{.foo}}",
				args: map[string]interface{}{
					"foo": make(chan int), // not templatable
				},
			},
			want: want{
				out: "",
				err: errors.New("template: :1:7: executing \"\" at <{{.foo}}>: can't print {{.foo}} of type chan int"),
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := renderCommand(tc.args.line, tc.args.args)

			EqualError(t, tc.want.err, err)
			require.Equal(t, tc.want.out, got)
		})
	}
}

package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/assh/v2/pkg/config"
)

func TestTCPDestination(t *testing.T) {
	for _, test := range []struct {
		host string
		want string
	}{
		{host: "example.test", want: "example.test:2222"},
		{host: "192.0.2.1", want: "192.0.2.1:2222"},
		{host: "::1", want: "[::1]:2222"},
		{host: "[::1]", want: "[::1]:2222"},
		{host: "fe80::1%eth0", want: "[fe80::1%eth0]:2222"},
	} {
		t.Run(test.host, func(t *testing.T) {
			require.Equal(t, test.want, tcpDestination(test.host, "2222"))
		})
	}
}

func TestProxyDryRunFormatsIPv6Destination(t *testing.T) {
	host := config.NewHost("ipv6")
	host.HostName = "::1"
	host.Port = "2222"
	require.ErrorContains(t, proxyGo(host, true), "'[::1]:2222'")
}

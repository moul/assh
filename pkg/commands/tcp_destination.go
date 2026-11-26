package commands

import (
	"net"
	"strings"
)

func tcpDestination(hostname, port string) string {
	// The port is stored separately, so either spelling of an IPv6 host is valid.
	if strings.HasPrefix(hostname, "[") && strings.HasSuffix(hostname, "]") {
		hostname = hostname[1 : len(hostname)-1]
	}
	return net.JoinHostPort(hostname, port)
}

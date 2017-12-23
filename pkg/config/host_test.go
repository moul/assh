package config

import (
	"fmt"
	"os/user"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHost_ApplyDefaults(t *testing.T) {
	Convey("Testing Host.ApplyDefaults", t, func() {
		Convey("Standard configuration", func() {
			host := &Host{
				name:     "example",
				HostName: "example.com",
				User:     "root",
			}
			defaults := &Host{
				User: "bobby",
				Port: "42",
			}
			host.ApplyDefaults(defaults)
			So(host.Port, ShouldEqual, "42")
			So(host.Name(), ShouldEqual, "example")
			So(host.HostName, ShouldEqual, "example.com")
			So(host.User, ShouldEqual, "root")
			So(len(host.Gateways), ShouldEqual, 0)
			So(host.ProxyCommand, ShouldEqual, "")
			So(len(host.ResolveNameservers), ShouldEqual, 0)
			So(host.ResolveCommand, ShouldEqual, "")
			So(host.ControlPath, ShouldEqual, "")
		})
		Convey("Empty configuration", func() {
			host := &Host{}
			defaults := &Host{}
			host.ApplyDefaults(defaults)
			So(host.Port, ShouldEqual, "22")
			So(host.Name(), ShouldEqual, "")
			So(host.HostName, ShouldEqual, "")
			So(host.User, ShouldEqual, "")
			So(len(host.Gateways), ShouldEqual, 0)
			So(host.ProxyCommand, ShouldEqual, "")
			So(len(host.ResolveNameservers), ShouldEqual, 0)
			So(host.ResolveCommand, ShouldEqual, "")
			So(host.ControlPath, ShouldEqual, "")
		})
	})
}

func TestHost_ExpandString(t *testing.T) {
	Convey("Testing Host.ExpandString()", t, func() {
		host := NewHost("abc")
		host.HostName = "1.2.3.4"
		host.Port = "42"

		var input, output, expected string

		input = "ls -la"
		output = host.ExpandString(input, "")
		expected = "ls -la"
		So(output, ShouldEqual, expected)

		input = "nc %h %p"
		output = host.ExpandString(input, "")
		expected = "nc 1.2.3.4 42"
		So(output, ShouldEqual, expected)

		input = "ssh %name"
		output = host.ExpandString(input, "")
		expected = "ssh abc"
		So(output, ShouldEqual, expected)

		input = "echo %h %p %name %h %p %name"
		output = host.ExpandString(input, "")
		expected = "echo 1.2.3.4 42 abc 1.2.3.4 42 abc"
		So(output, ShouldEqual, expected)

		input = "echo %g"
		output = host.ExpandString(input, "")
		expected = "echo "
		So(output, ShouldEqual, expected)

		input = "echo %g"
		output = host.ExpandString(input, "def")
		expected = "echo def"
		So(output, ShouldEqual, expected)
	})
}

func TestHost_Clone(t *testing.T) {
	Convey("Testing Host.Clone()", t, func() {
		a := NewHost("abc")
		a.HostName = "1.2.3.4"
		a.Port = "42"

		b := a.Clone()

		So(a, ShouldNotEqual, b)
		So(a.HostName, ShouldEqual, "1.2.3.4")
		So(b.HostName, ShouldEqual, "1.2.3.4")
		So(a.Port, ShouldEqual, "42")
		So(b.Port, ShouldEqual, "42")
	})
}

func TestHost_Prototype(t *testing.T) {
	Convey("Testing Host.Prototype()", t, func() {
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}

		host := NewHost("abc")
		So(host.Prototype(), ShouldEqual, fmt.Sprintf("%s@abc:22", currentUser.Username))

		host = NewHost("abc-*")
		So(host.Prototype(), ShouldEqual, fmt.Sprintf("%s@[dynamic]:22", currentUser.Username))

		host = NewHost("abc")
		host.User = "toto"
		host.HostName = "1.2.3.4"
		host.Port = "42"
		So(host.Prototype(), ShouldEqual, "toto@1.2.3.4:42")
	})
}

func TestHost_Matches(t *testing.T) {
	Convey("Testing Host.Matches()", t, func() {
		host := NewHost("abc")
		So(host.Matches("a"), ShouldBeTrue)
		So(host.Matches("ab"), ShouldBeTrue)
		So(host.Matches("abc"), ShouldBeTrue)
		So(host.Matches("bcd"), ShouldBeFalse)
		So(host.Matches("b"), ShouldBeTrue)
		So(host.Matches("bc"), ShouldBeTrue)
		So(host.Matches("c"), ShouldBeTrue)

		host.User = "bcd"
		So(host.Matches("a"), ShouldBeTrue)
		So(host.Matches("ab"), ShouldBeTrue)
		So(host.Matches("abc"), ShouldBeTrue)
		So(host.Matches("bcd"), ShouldBeTrue)
		So(host.Matches("b"), ShouldBeTrue)
		So(host.Matches("bc"), ShouldBeTrue)
		So(host.Matches("c"), ShouldBeTrue)
	})
}

func TestHost_Validate(t *testing.T) {
	Convey("Testing Host.Validate()", t, FailureContinues, func() {
		host := NewHost("abc")

		errs := host.Validate()
		So(len(errs), ShouldEqual, 0)

		for _, value := range []string{"yes", "no", "ask", "auto", "autoask", "", "Yes", "YES", "yEs", " yes "} {
			host.ControlMaster = value
			errs = host.Validate()
			So(len(errs), ShouldEqual, 0)
		}

		for _, value := range []string{"blah blah", "invalid"} {
			host.ControlMaster = value
			errs = host.Validate()
			So(len(errs), ShouldEqual, 1)
		}
	})
}

func TestHost_Options(t *testing.T) {
	Convey("Testing Host.Options()", t, func() {
		host := NewHost("abc")
		options := host.Options()
		So(len(options), ShouldEqual, 0)
		So(options, ShouldResemble, OptionsList{})

		host = dummyHost()
		options = host.Options()
		So(len(options), ShouldEqual, 95)
		So(options, ShouldResemble, OptionsList{{Name: "AddKeysToAgent", Value: "yes"}, {Name: "AddressFamily", Value: "any"}, {Name: "AskPassGUI", Value: "yes"}, {Name: "BatchMode", Value: "no"}, {Name: "CanonicalDomains", Value: "42.am"}, {Name: "CanonicalizeFallbackLocal", Value: "no"}, {Name: "CanonicalizeHostname", Value: "yes"}, {Name: "CanonicalizeMaxDots", Value: "1"}, {Name: "CanonicalizePermittedCNAMEs", Value: "*.a.example.com:*.b.example.com:*.c.example.com"}, {Name: "ChallengeResponseAuthentication", Value: "yes"}, {Name: "CheckHostIP", Value: "yes"}, {Name: "Cipher", Value: "blowfish"}, {Name: "Ciphers", Value: "aes128-ctr,aes192-ctr,aes256-ctr,test"}, {Name: "ClearAllForwardings", Value: "yes"}, {Name: "Compression", Value: "yes"}, {Name: "CompressionLevel", Value: "6"}, {Name: "ConnectionAttempts", Value: "1"}, {Name: "ConnectTimeout", Value: "10"}, {Name: "ControlMaster", Value: "yes"}, {Name: "ControlPath", Value: "/tmp/%L-%l-%n-%p-%u-%r-%C-%h"}, {Name: "ControlPersist", Value: "yes"}, {Name: "DynamicForward", Value: "0.0.0.0:4242"}, {Name: "DynamicForward", Value: "0.0.0.0:4343"}, {Name: "EnableSSHKeysign", Value: "yes"}, {Name: "EscapeChar", Value: "~"}, {Name: "ExitOnForwardFailure", Value: "yes"}, {Name: "FingerprintHash", Value: "sha256"}, {Name: "ForwardAgent", Value: "yes"}, {Name: "ForwardX11", Value: "yes"}, {Name: "ForwardX11Timeout", Value: "42"}, {Name: "ForwardX11Trusted", Value: "yes"}, {Name: "GatewayPorts", Value: "yes"}, {Name: "GlobalKnownHostsFile", Value: "/etc/ssh/ssh_known_hosts /tmp/ssh_known_hosts"}, {Name: "GSSAPIAuthentication", Value: "no"}, {Name: "GSSAPIClientIdentity", Value: "moul"}, {Name: "GSSAPIDelegateCredentials", Value: "no"}, {Name: "GSSAPIKeyExchange", Value: "no"}, {Name: "GSSAPIRenewalForcesRekey", Value: "no"}, {Name: "GSSAPIServerIdentity", Value: "gssapi.example.com"}, {Name: "GSSAPITrustDNS", Value: "no"}, {Name: "HashKnownHosts", Value: "no"}, {Name: "HostbasedAuthentication", Value: "no"}, {Name: "HostbasedKeyTypes", Value: "*"}, {Name: "HostKeyAlgorithms", Value: "ecdsa-sha2-nistp256-cert-v01@openssh.com"}, {Name: "HostKeyAlias", Value: "z"}, {Name: "IdentitiesOnly", Value: "yes"}, {Name: "IdentityFile", Value: "~/.ssh/identity"}, {Name: "IdentityFile", Value: "~/.ssh/identity2"}, {Name: "IgnoreUnknown", Value: "testtest"}, {Name: "IPQoS", Value: "lowdelay highdelay"}, {Name: "KbdInteractiveAuthentication", Value: "yes"}, {Name: "KbdInteractiveDevices", Value: "bsdauth,test"}, {Name: "KexAlgorithms", Value: "curve25519-sha256@libssh.org,test"}, {Name: "KeychainIntegration", Value: "yes"}, {Name: "LocalCommand", Value: "echo %h > /tmp/logs"}, {Name: "LocalForward", Value: "0.0.0.0:1234"}, {Name: "LocalForward", Value: "0.0.0.0:1235"}, {Name: "LogLevel", Value: "DEBUG3"}, {Name: "MACs", Value: "umac-64-etm@openssh.com,umac-128-etm@openssh.com,test"}, {Name: "Match", Value: "all"}, {Name: "NoHostAuthenticationForLocalhost", Value: "yes"}, {Name: "NumberOfPasswordPrompts", Value: "3"}, {Name: "PasswordAuthentication", Value: "yes"}, {Name: "PermitLocalCommand", Value: "yes"}, {Name: "PKCS11Provider", Value: "/a/b/c/pkcs11.so"}, {Name: "Port", Value: "22"}, {Name: "PreferredAuthentications", Value: "gssapi-with-mic,hostbased,publickey"}, {Name: "Protocol", Value: "2,3"}, {Name: "ProxyUseFdpass", Value: "no"}, {Name: "PubkeyAcceptedKeyTypes", Value: "+ssh-dss"}, {Name: "PubkeyAuthentication", Value: "yes"}, {Name: "RekeyLimit", Value: "default none"}, {Name: "RemoteForward", Value: "0.0.0.0:1234"}, {Name: "RemoteForward", Value: "0.0.0.0:1235"}, {Name: "RequestTTY", Value: "yes"}, {Name: "RevokedHostKeys", Value: "/a/revoked-keys"}, {Name: "RhostsRSAAuthentication", Value: "no"}, {Name: "RSAAuthentication", Value: "yes"}, {Name: "SendEnv", Value: "CUSTOM_*,TEST"}, {Name: "SendEnv", Value: "TEST2"}, {Name: "ServerAliveCountMax", Value: "3"}, {Name: "StreamLocalBindMask", Value: "0177"}, {Name: "StreamLocalBindUnlink", Value: "no"}, {Name: "StrictHostKeyChecking", Value: "ask"}, {Name: "TCPKeepAlive", Value: "yes"}, {Name: "Tunnel", Value: "yes"}, {Name: "TunnelDevice", Value: "any:any"}, {Name: "UpdateHostKeys", Value: "ask"}, {Name: "UseKeychain", Value: "no"}, {Name: "UsePrivilegedPort", Value: "no"}, {Name: "User", Value: "moul"}, {Name: "UserKnownHostsFile", Value: "~/.ssh/known_hosts ~/.ssh/known_hosts2 /tmp/known_hosts"}, {Name: "VerifyHostKeyDNS", Value: "no"}, {Name: "VisualHostKey", Value: "yes"}, {Name: "XAuthLocation", Value: "xauth"}})
	})
}

func dummyHost() *Host {
	return &Host{
		// ssh-config fields
		AddKeysToAgent:                  "yes",
		AddressFamily:                   "any",
		AskPassGUI:                      "yes",
		BatchMode:                       "no",
		BindAddress:                     "",
		CanonicalDomains:                "42.am",
		CanonicalizeFallbackLocal:       "no",
		CanonicalizeHostname:            "yes",
		CanonicalizeMaxDots:             "1",
		CanonicalizePermittedCNAMEs:     "*.a.example.com:*.b.example.com:*.c.example.com",
		ChallengeResponseAuthentication: "yes",
		CheckHostIP:                     "yes",
		Cipher:                          "blowfish",
		Ciphers:                         []string{"aes128-ctr,aes192-ctr,aes256-ctr", "test"},
		ClearAllForwardings:             "yes",
		Compression:                     "yes",
		CompressionLevel:                6,
		ConnectionAttempts:              "1",
		ConnectTimeout:                  10,
		ControlMaster:                   "yes",
		ControlPath:                     "/tmp/%L-%l-%n-%p-%u-%r-%C-%h",
		ControlPersist:                  "yes",
		DynamicForward:                  []string{"0.0.0.0:4242", "0.0.0.0:4343"},
		EnableSSHKeysign:                "yes",
		EscapeChar:                      "~",
		ExitOnForwardFailure:            "yes",
		FingerprintHash:                 "sha256",
		ForwardAgent:                    "yes",
		ForwardX11:                      "yes",
		ForwardX11Timeout:               42,
		ForwardX11Trusted:               "yes",
		GatewayPorts:                    "yes",
		GlobalKnownHostsFile:            []string{"/etc/ssh/ssh_known_hosts", "/tmp/ssh_known_hosts"},
		GSSAPIAuthentication:            "no",
		GSSAPIKeyExchange:               "no",
		GSSAPIClientIdentity:            "moul",
		GSSAPIServerIdentity:            "gssapi.example.com",
		GSSAPIDelegateCredentials:       "no",
		GSSAPIRenewalForcesRekey:        "no",
		GSSAPITrustDNS:                  "no",
		HashKnownHosts:                  "no",
		HostbasedAuthentication:         "no",
		HostbasedKeyTypes:               "*",
		HostKeyAlgorithms:               "ecdsa-sha2-nistp256-cert-v01@openssh.com",
		HostKeyAlias:                    "z",
		IdentitiesOnly:                  "yes",
		IdentityFile:                    []string{"~/.ssh/identity", "~/.ssh/identity2"},
		IgnoreUnknown:                   "testtest", // FIXME: looks very interesting to generate .ssh/config without comments !
		IPQoS:                           []string{"lowdelay", "highdelay"},
		KbdInteractiveAuthentication: "yes",
		KbdInteractiveDevices:        []string{"bsdauth", "test"},
		KeychainIntegration:          "yes",
		KexAlgorithms:                []string{"curve25519-sha256@libssh.org", "test"}, // for all algorithms/ciphers, we could have an "assh diagnose" that warns about unsafe connections
		LocalCommand:                 "echo %h > /tmp/logs",
		LocalForward:                 []string{"0.0.0.0:1234", "0.0.0.0:1235"}, // FIXME: may be a list
		LogLevel:                     "DEBUG3",
		MACs:                         []string{"umac-64-etm@openssh.com,umac-128-etm@openssh.com", "test"},
		Match:                        "all",
		NoHostAuthenticationForLocalhost: "yes",
		NumberOfPasswordPrompts:          "3",
		PasswordAuthentication:           "yes",
		PermitLocalCommand:               "yes",
		PKCS11Provider:                   "/a/b/c/pkcs11.so",
		Port:                             "22",
		PreferredAuthentications: "gssapi-with-mic,hostbased,publickey",
		Protocol:                 []string{"2", "3"},
		ProxyUseFdpass:           "no",
		PubkeyAcceptedKeyTypes:   "+ssh-dss",
		PubkeyAuthentication:     "yes",
		RekeyLimit:               "default none",
		RemoteForward:            []string{"0.0.0.0:1234", "0.0.0.0:1235"},
		RequestTTY:               "yes",
		RevokedHostKeys:          "/a/revoked-keys",
		RhostsRSAAuthentication:  "no",
		RSAAuthentication:        "yes",
		SendEnv:                  []string{"CUSTOM_*,TEST", "TEST2"},
		ServerAliveCountMax:      3,
		ServerAliveInterval:      0,
		StreamLocalBindMask:      "0177",
		StreamLocalBindUnlink:    "no",
		StrictHostKeyChecking:    "ask",
		TCPKeepAlive:             "yes",
		Tunnel:                   "yes",
		TunnelDevice:             "any:any",
		UpdateHostKeys:           "ask",
		UseKeychain:              "no",
		UsePrivilegedPort:        "no",
		User:                     "moul",
		UserKnownHostsFile:       []string{"~/.ssh/known_hosts ~/.ssh/known_hosts2", "/tmp/known_hosts"},
		VerifyHostKeyDNS:         "no",
		VisualHostKey:            "yes",
		XAuthLocation:            "xauth",

		// ssh-config fields with a different behavior
		ProxyCommand: "nc %h %p",
		HostName:     "zzz.com",

		// assh fields
		isDefault:          false,
		Inherits:           []string{},
		Gateways:           []string{},
		Aliases:            []string{},
		ResolveNameservers: []string{},
		ResolveCommand:     "",
	}
}

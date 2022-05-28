package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	composeyaml "github.com/docker/libcompose/yaml"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	yamlConfig = `
hosts:

  aaa:
    HostName: 1.2.3.4

  BBB:
    Port: ${ENV_VAR_PORT}
    HostName: $ENV_VAR_HOSTNAME
    User: user-$ENV_VAR_USER-user
    LocalCommand: ${ENV_VAR_LOCALCOMMAND:-hello}
    IdentityFile: ${NON_EXISTING_ENV_VAR}
    Comment: Hello World !

  ccc:
    HostName: 5.6.7.8
    Port: 24
    User: toor
  "*.ddd":
    HostName: 1.3.5.7

  eee:
    Inherits:
    - aaa
    - BBB
    - aaa
    ControlMasterMkdir: true
    Comment:
    - AAA
    - BBB

  FFF:
    Inherits:
    - BBB
    - eee
    - "*.ddd"
    Comment: >
      First line
      Second line
      Third line
    RemoteCommand: date >> /tmp/logs

  ggg:
    Gateways:
    - direct
    - FFF

  hhh:
    Gateways:
    - ggg
    - direct

  iii:
    Gateways: test.ddd

  jjj:
    HostName: "%h.jjjjj"

  "*.kkk":
    HostName: "%h.kkkkk"

  "lll-*":
    HostName: "%h.lll"

  "toto[1-5]toto":
    User: toto1

  "toto[7-9]toto":
    User: toto2

  nnn:
    Inherits: mmm
    User: nnnn

  ooo1:
    Port: 23
    Aliases:
    - ooo11
    - ooo12

  ooo2:
    Port: 24
    Aliases:
    - ooo21
    - ooo22

templates:

  kkk:
    Port: 25
    User: kkkk

  lll:
    HostName: 5.5.5.5

  mmm:
    Inherits: iii

defaults:
  Port: 22
  User: root

includes:
  - /path/to/dir/*.yml
  - /path/to/file.yml
`
)

func TestNew(t *testing.T) {
	Convey("Testing New()", t, func() {
		config := New()

		So(len(config.Hosts), ShouldEqual, 0)
		So(config.Defaults.Port, ShouldEqual, "")
		So(config.Defaults.HostName, ShouldEqual, "")
		So(config.Defaults.User, ShouldEqual, "")
	})
}

func dummyConfig() *Config {
	config := New()
	config.Hosts["toto"] = &Host{ // TOTO in yaml becomes toto
		HostName: "1.2.3.4",
	}
	config.Hosts["titi"] = &Host{
		HostName:           "tata",
		Port:               "23",
		User:               "moul",
		ProxyCommand:       "nc -v 4242",
		ControlMasterMkdir: "true",
		Comment:            []string{"Hello World"},
	}
	config.Hosts["tonton"] = &Host{
		ResolveNameservers: []string{"a.com", "1.2.3.4"},
		Comment:            []string{"AAA", "BBB"},
	}
	config.Hosts["toutou"] = &Host{
		ResolveCommand: "dig -t %h",
		Comment:        []string{"First line Second line Third line\n"},
		RemoteCommand:  "date >> /tmp/logs",
	}
	config.Hosts["tutu"] = &Host{
		Gateways: []string{"titi", "direct", "1.2.3.4"},
		Inherits: []string{"TOTO", "tutu", "*.ddd"},
	}
	config.Hosts["empty"] = &Host{}
	config.Hosts["tata"] = &Host{
		Inherits: []string{"tutu", "titi", "TOTO", "tutu"},
	}
	config.Hosts["*.ddd"] = &Host{
		HostName:               "1.3.5.7",
		PasswordAuthentication: "yes",
	}
	config.Defaults = Host{
		Port: "22",
		User: "root",
	}
	config.Templates["mmm"] = &Host{
		Port:     "25",
		User:     "mmmm",
		HostName: "5.5.5.5",
		Inherits: []string{"tata"},
	}
	config.Hosts["nnn"] = &Host{
		Port:     "26",
		Inherits: []string{"mmm"},
	}
	config.Hosts["ooo1"] = &Host{
		Port:    "23",
		Aliases: []string{"ooo11", "ooo12"},
	}
	config.Hosts["ooo2"] = &Host{
		Port:    "24",
		Aliases: []string{"ooo21", "ooo22"},
	}
	config.Hosts["toto[1-5]toto"] = &Host{
		User: "toto1",
	}
	config.Hosts["toto[7-9]toto"] = &Host{
		User: "toto2",
	}
	config.Hosts["zzz"] = &Host{
		// ssh-config fields
		AddressFamily:                    "any",
		AskPassGUI:                       "yes",
		BatchMode:                        "no",
		BindAddress:                      "",
		CanonicalDomains:                 "42.am",
		CanonicalizeFallbackLocal:        "no",
		CanonicalizeHostname:             "yes",
		CanonicalizeMaxDots:              "1",
		CanonicalizePermittedCNAMEs:      "*.a.example.com:*.b.example.com:*.c.example.com",
		ChallengeResponseAuthentication:  "yes",
		CheckHostIP:                      "yes",
		Cipher:                           "blowfish",
		Ciphers:                          []string{"aes128-ctr,aes192-ctr", "aes256-ctr"},
		ClearAllForwardings:              "yes",
		Compression:                      "yes",
		CompressionLevel:                 6,
		ConnectionAttempts:               "1",
		ConnectTimeout:                   10,
		ControlMaster:                    "yes",
		ControlPath:                      "/tmp/%L-%l-%n-%p-%u-%r-%C-%h",
		ControlPersist:                   "yes",
		DynamicForward:                   []string{"0.0.0.0:4242", "0.0.0.0:4343"},
		EnableSSHKeysign:                 "yes",
		EscapeChar:                       "~",
		ExitOnForwardFailure:             "yes",
		FingerprintHash:                  "sha256",
		ForwardAgent:                     "yes",
		ForwardX11:                       "yes",
		ForwardX11Timeout:                42,
		ForwardX11Trusted:                "yes",
		GatewayPorts:                     "yes",
		GlobalKnownHostsFile:             []string{"/etc/ssh/ssh_known_hosts", "/tmp/ssh_known_hosts"},
		GSSAPIAuthentication:             "no",
		GSSAPIKeyExchange:                "no",
		GSSAPIClientIdentity:             "moul",
		GSSAPIServerIdentity:             "gssapi.example.com",
		GSSAPIDelegateCredentials:        "no",
		GSSAPIRenewalForcesRekey:         "no",
		GSSAPITrustDNS:                   "no",
		HashKnownHosts:                   "no",
		HostbasedAuthentication:          "no",
		HostbasedKeyTypes:                "*",
		HostKeyAlgorithms:                []string{"ecdsa-sha2-nistp256-cert-v01@openssh.com", "test"},
		HostKeyAlias:                     "z",
		IdentitiesOnly:                   "yes",
		IdentityFile:                     []string{"~/.ssh/identity", "~/.ssh/identity2"},
		IgnoreUnknown:                    "testtest", // FIXME: looks very interesting to generate .ssh/config without comments !
		IPQoS:                            []string{"lowdelay", "highdelay"},
		KbdInteractiveAuthentication:     "yes",
		KbdInteractiveDevices:            []string{"bsdauth", "test"},
		KeychainIntegration:              "yes",
		KexAlgorithms:                    []string{"curve25519-sha256@libssh.org", "test"}, // for all algorithms/ciphers, we could have an "assh diagnose" that warns about unsafe connections
		LocalCommand:                     "echo %h > /tmp/logs",
		LocalForward:                     []string{"0.0.0.0:1234", "0.0.0.0:1235"},
		LogLevel:                         "DEBUG3",
		MACs:                             []string{"umac-64-etm@openssh.com,umac-128-etm@openssh.com", "test"},
		Match:                            "all",
		NoHostAuthenticationForLocalhost: "yes",
		NumberOfPasswordPrompts:          "3",
		PasswordAuthentication:           "yes",
		PermitLocalCommand:               "yes",
		PKCS11Provider:                   "/a/b/c/pkcs11.so",
		Port:                             "22",
		PreferredAuthentications:         "gssapi-with-mic,hostbased,publickey",
		Protocol:                         []string{"2", "3"},
		ProxyUseFdpass:                   "no",
		PubkeyAuthentication:             "yes",
		RekeyLimit:                       "default none",
		RemoteForward:                    []string{"0.0.0.0:1234", "0.0.0.0:1255"},
		RequestTTY:                       "yes",
		RevokedHostKeys:                  "/a/revoked-keys",
		RhostsRSAAuthentication:          "no",
		RSAAuthentication:                "yes",
		SendEnv:                          []string{"CUSTOM_*,TEST", "TEST2"},
		ServerAliveCountMax:              3,
		ServerAliveInterval:              0,
		StreamLocalBindMask:              "0177",
		StreamLocalBindUnlink:            "no",
		StrictHostKeyChecking:            "ask",
		TCPKeepAlive:                     "yes",
		Tunnel:                           "yes",
		TunnelDevice:                     "any:any",
		UpdateHostKeys:                   "ask",
		UseKeychain:                      "no",
		UsePrivilegedPort:                "no",
		User:                             "moul",
		UserKnownHostsFile:               []string{"~/.ssh/known_hosts ~/.ssh/known_hosts2", "/tmp/known_hosts"},
		VerifyHostKeyDNS:                 "no",
		VisualHostKey:                    "yes",
		XAuthLocation:                    "xauth",

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
	config.applyMissingNames()
	return config
}

func TestConfig(t *testing.T) {
	Convey("Testing dummyConfig", t, func() {
		config := dummyConfig()

		So(len(config.Hosts), ShouldEqual, 14)

		So(config.Hosts["toto"].HostName, ShouldEqual, "1.2.3.4")
		So(config.Hosts["toto"].Port, ShouldEqual, "")
		So(config.Hosts["toto"].name, ShouldEqual, "toto")
		So(config.Hosts["toto"].isDefault, ShouldEqual, false)

		So(config.Hosts["titi"].HostName, ShouldEqual, "tata")
		So(config.Hosts["titi"].User, ShouldEqual, "moul")
		So(config.Hosts["titi"].ProxyCommand, ShouldEqual, "nc -v 4242")
		So(BoolVal(config.Hosts["titi"].ControlMasterMkdir), ShouldBeTrue)
		So(config.Hosts["titi"].Port, ShouldEqual, "23")
		So(config.Hosts["titi"].isDefault, ShouldEqual, false)

		So(config.Hosts["tonton"].isDefault, ShouldEqual, false)
		So(config.Hosts["tonton"].Port, ShouldEqual, "")
		So(config.Hosts["tonton"].ResolveNameservers, ShouldResemble, composeyaml.Stringorslice{"a.com", "1.2.3.4"})

		So(config.Hosts["toutou"].isDefault, ShouldEqual, false)
		So(config.Hosts["toutou"].Port, ShouldEqual, "")
		So(config.Hosts["toutou"].ResolveCommand, ShouldEqual, "dig -t %h")

		So(config.Hosts["tutu"].isDefault, ShouldEqual, false)
		So(config.Hosts["tutu"].Port, ShouldEqual, "")
		So(config.Hosts["tutu"].Gateways, ShouldResemble, composeyaml.Stringorslice{"titi", "direct", "1.2.3.4"})

		So(config.Hosts["*.ddd"].isDefault, ShouldEqual, false)
		So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")

		So(config.Hosts["empty"].isDefault, ShouldEqual, false)
		So(config.Hosts["empty"].Port, ShouldEqual, "")

		So(len(config.Templates), ShouldEqual, 1)

		So(config.Defaults.User, ShouldEqual, "root")
		So(config.Defaults.Port, ShouldEqual, "22")
		So(config.Defaults.isDefault, ShouldEqual, true)
	})
}

func TestConfig_LoadConfig(t *testing.T) {
	Convey("Testing Config.LoadConfig", t, func() {
		Convey("standard", func() {
			config := New()
			err := config.LoadConfig(strings.NewReader(yamlConfig))
			So(err, ShouldBeNil)
			So(len(config.Hosts), ShouldEqual, 17)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, "")
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(config.Hosts["bbb"].Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(config.Hosts["bbb"].User, ShouldEqual, "user-$ENV_VAR_USER-user")
			So(config.Hosts["bbb"].IdentityFile[0], ShouldEqual, "${NON_EXISTING_ENV_VAR}")
			So(config.Hosts["bbb"].LocalCommand, ShouldEqual, "${ENV_VAR_LOCALCOMMAND:-hello}")
			So(config.Hosts["ccc"].HostName, ShouldEqual, "5.6.7.8")
			So(config.Hosts["ccc"].Port, ShouldEqual, "24")
			So(config.Hosts["ccc"].User, ShouldEqual, "toor")
			So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")
			So(config.Hosts["*.ddd"].Port, ShouldEqual, "")
			So(config.Hosts["*.ddd"].User, ShouldEqual, "")
			So(config.Defaults.Port, ShouldEqual, "22")
			So(config.Defaults.User, ShouldEqual, "root")
			So(len(config.Templates), ShouldEqual, 3)
		})
	})
}

func TestConfig_JSONString(t *testing.T) {
	Convey("Testing Config.JSONString", t, func() {
		Convey("dummyConfig", func() {
			config := dummyConfig()
			expected := `{
  "hosts": {
    "*.ddd": {
      "PasswordAuthentication": "yes",
      "HostName": "1.3.5.7"
    },
    "empty": {},
    "nnn": {
      "Port": "26",
      "Inherits": [
        "mmm"
      ]
    },
    "ooo1": {
      "Port": "23",
      "Aliases": [
        "ooo11",
        "ooo12"
      ]
    },
    "ooo2": {
      "Port": "24",
      "Aliases": [
        "ooo21",
        "ooo22"
      ]
    },
    "tata": {
      "Inherits": [
        "tutu",
        "titi",
        "toto",
        "tutu"
      ]
    },
    "titi": {
      "Port": "23",
      "User": "moul",
      "HostName": "tata",
      "ProxyCommand": "nc -v 4242",
      "ControlMasterMkdir": "true",
      "Comment": [
        "Hello World"
      ]
    },
    "tonton": {
      "ResolveNameservers": [
        "a.com",
        "1.2.3.4"
      ],
      "Comment": [
        "AAA",
        "BBB"
      ]
    },
    "toto": {
      "HostName": "1.2.3.4"
    },
    "toto[1-5]toto": {
      "User": "toto1"
    },
    "toto[7-9]toto": {
      "User": "toto2"
    },
    "toutou": {
      "RemoteCommand": "date \u003e\u003e /tmp/logs",
      "ResolveCommand": "dig -t %h",
      "Comment": [
        "First line Second line Third line\n"
      ]
    },
    "tutu": {
      "Inherits": [
        "toto",
        "tutu",
        "*.ddd"
      ],
      "Gateways": [
        "titi",
        "direct",
        "1.2.3.4"
      ]
    },
    "zzz": {
      "AddressFamily": "any",
      "AskPassGUI": "yes",
      "BatchMode": "no",
      "CanonicalDomains": "42.am",
      "CanonicalizeFallbackLocal": "no",
      "CanonicalizeHostname": "yes",
      "CanonicalizeMaxDots": "1",
      "CanonicalizePermittedCNAMEs": "*.a.example.com:*.b.example.com:*.c.example.com",
      "ChallengeResponseAuthentication": "yes",
      "CheckHostIP": "yes",
      "Cipher": "blowfish",
      "Ciphers": [
        "aes128-ctr,aes192-ctr",
        "aes256-ctr"
      ],
      "ClearAllForwardings": "yes",
      "Compression": "yes",
      "CompressionLevel": 6,
      "ConnectionAttempts": "1",
      "ConnectTimeout": 10,
      "ControlMaster": "yes",
      "ControlPath": "/tmp/%L-%l-%n-%p-%u-%r-%C-%h",
      "ControlPersist": "yes",
      "DynamicForward": [
        "0.0.0.0:4242",
        "0.0.0.0:4343"
      ],
      "EnableSSHKeysign": "yes",
      "EscapeChar": "~",
      "ExitOnForwardFailure": "yes",
      "FingerprintHash": "sha256",
      "ForwardAgent": "yes",
      "ForwardX11": "yes",
      "ForwardX11Timeout": 42,
      "ForwardX11Trusted": "yes",
      "GatewayPorts": "yes",
      "GlobalKnownHostsFile": [
        "/etc/ssh/ssh_known_hosts",
        "/tmp/ssh_known_hosts"
      ],
      "GSSAPIAuthentication": "no",
      "GSSAPIClientIdentity": "moul",
      "GSSAPIDelegateCredentials": "no",
      "GSSAPIKeyExchange": "no",
      "GSSAPIRenewalForcesRekey": "no",
      "GSSAPIServerIdentity": "gssapi.example.com",
      "GSSAPITrustDNS": "no",
      "HashKnownHosts": "no",
      "HostbasedAuthentication": "no",
      "HostbasedKeyTypes": "*",
      "HostKeyAlgorithms": [
        "ecdsa-sha2-nistp256-cert-v01@openssh.com",
        "test"
      ],
      "HostKeyAlias": "z",
      "IdentitiesOnly": "yes",
      "IdentityFile": [
        "~/.ssh/identity",
        "~/.ssh/identity2"
      ],
      "IgnoreUnknown": "testtest",
      "IPQoS": [
        "lowdelay",
        "highdelay"
      ],
      "KbdInteractiveAuthentication": "yes",
      "KbdInteractiveDevices": [
        "bsdauth",
        "test"
      ],
      "KexAlgorithms": [
        "curve25519-sha256@libssh.org",
        "test"
      ],
      "KeychainIntegration": "yes",
      "LocalCommand": "echo %h \u003e /tmp/logs",
      "LocalForward": [
        "0.0.0.0:1234",
        "0.0.0.0:1235"
      ],
      "LogLevel": "DEBUG3",
      "MACs": [
        "umac-64-etm@openssh.com,umac-128-etm@openssh.com",
        "test"
      ],
      "Match": "all",
      "NoHostAuthenticationForLocalhost": "yes",
      "NumberOfPasswordPrompts": "3",
      "PasswordAuthentication": "yes",
      "PermitLocalCommand": "yes",
      "PKCS11Provider": "/a/b/c/pkcs11.so",
      "Port": "22",
      "PreferredAuthentications": "gssapi-with-mic,hostbased,publickey",
      "Protocol": [
        "2",
        "3"
      ],
      "ProxyUseFdpass": "no",
      "PubkeyAuthentication": "yes",
      "RekeyLimit": "default none",
      "RemoteForward": [
        "0.0.0.0:1234",
        "0.0.0.0:1255"
      ],
      "RequestTTY": "yes",
      "RevokedHostKeys": "/a/revoked-keys",
      "RhostsRSAAuthentication": "no",
      "RSAAuthentication": "yes",
      "SendEnv": [
        "CUSTOM_*,TEST",
        "TEST2"
      ],
      "ServerAliveCountMax": 3,
      "StreamLocalBindMask": "0177",
      "StreamLocalBindUnlink": "no",
      "StrictHostKeyChecking": "ask",
      "TCPKeepAlive": "yes",
      "Tunnel": "yes",
      "TunnelDevice": "any:any",
      "UpdateHostKeys": "ask",
      "UseKeychain": "no",
      "UsePrivilegedPort": "no",
      "User": "moul",
      "UserKnownHostsFile": [
        "~/.ssh/known_hosts ~/.ssh/known_hosts2",
        "/tmp/known_hosts"
      ],
      "VerifyHostKeyDNS": "no",
      "VisualHostKey": "yes",
      "XAuthLocation": "xauth",
      "HostName": "zzz.com",
      "ProxyCommand": "nc %h %p"
    }
  },
  "templates": {
    "mmm": {
      "Port": "25",
      "User": "mmmm",
      "HostName": "5.5.5.5",
      "Inherits": [
        "tata"
      ]
    }
  },
  "defaults": {
    "Port": "22",
    "User": "root",
    "Hooks": {}
  },
  "asshknownhostfile": "~/.ssh/assh_known_hosts"
}`
			json, err := config.JSONString()
			So(err, ShouldBeNil)
			So(string(json), ShouldEqual, expected)
		})
		Convey("yamlConfig", func() {
			config := New()
			err := config.LoadConfig(strings.NewReader(yamlConfig))
			So(err, ShouldBeNil)
			expected := `{
  "hosts": {
    "*.ddd": {
      "HostName": "1.3.5.7"
    },
    "*.kkk": {
      "HostName": "%h.kkkkk"
    },
    "aaa": {
      "HostName": "1.2.3.4"
    },
    "bbb": {
      "IdentityFile": [
        "${NON_EXISTING_ENV_VAR}"
      ],
      "LocalCommand": "${ENV_VAR_LOCALCOMMAND:-hello}",
      "Port": "${ENV_VAR_PORT}",
      "User": "user-$ENV_VAR_USER-user",
      "HostName": "$ENV_VAR_HOSTNAME",
      "Comment": [
        "Hello World !"
      ]
    },
    "ccc": {
      "Port": "24",
      "User": "toor",
      "HostName": "5.6.7.8"
    },
    "eee": {
      "Inherits": [
        "aaa",
        "bbb",
        "aaa"
      ],
      "ControlMasterMkdir": "true",
      "Comment": [
        "AAA",
        "BBB"
      ]
    },
    "fff": {
      "RemoteCommand": "date \u003e\u003e /tmp/logs",
      "Inherits": [
        "bbb",
        "eee",
        "*.ddd"
      ],
      "Comment": [
        "First line Second line Third line\n"
      ]
    },
    "ggg": {
      "Gateways": [
        "direct",
        "fff"
      ]
    },
    "hhh": {
      "Gateways": [
        "ggg",
        "direct"
      ]
    },
    "iii": {
      "Gateways": [
        "test.ddd"
      ]
    },
    "jjj": {
      "HostName": "%h.jjjjj"
    },
    "lll-*": {
      "HostName": "%h.lll"
    },
    "nnn": {
      "User": "nnnn",
      "Inherits": [
        "mmm"
      ]
    },
    "ooo1": {
      "Port": "23",
      "Aliases": [
        "ooo11",
        "ooo12"
      ]
    },
    "ooo2": {
      "Port": "24",
      "Aliases": [
        "ooo21",
        "ooo22"
      ]
    },
    "toto[1-5]toto": {
      "User": "toto1"
    },
    "toto[7-9]toto": {
      "User": "toto2"
    }
  },
  "templates": {
    "kkk": {
      "Port": "25",
      "User": "kkkk"
    },
    "lll": {
      "HostName": "5.5.5.5"
    },
    "mmm": {
      "Inherits": [
        "iii"
      ]
    }
  },
  "defaults": {
    "Port": "22",
    "User": "root",
    "Hooks": {}
  },
  "includes": [
    "/path/to/dir/*.yml",
    "/path/to/file.yml"
  ],
  "asshknownhostfile": "~/.ssh/assh_known_hosts"
}`
			json, err := config.JSONString()
			So(err, ShouldBeNil)
			So(string(json), ShouldEqual, expected)
		})
	})
}

func TestComputeHost(t *testing.T) {
	Convey("Testing computeHost()", t, func() {
		config := New()
		err := config.LoadConfig(strings.NewReader(yamlConfig))
		So(err, ShouldBeNil)

		Convey("Standard", func() {

		})
		Convey("Expand variables in HostName", func() {
			host := config.Hosts["jjj"]
			computed, err := computeHost(host, config, "jjj", false)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "%h.jjjjj")

			computed, err = computeHost(host, config, "jjj", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "jjj.jjjjj")

			host = config.Hosts["*.kkk"]
			computed, err = computeHost(host, config, "test.kkk", false)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "%h.kkkkk")

			computed, err = computeHost(host, config, "test.kkk", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "test.kkk.kkkkk")
		})
		Convey("Do not expand variables twice", func() {
			host := config.Hosts["lll-*"]
			computed, err := computeHost(host, config, "lll-42", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "lll-42.lll")

			computed, err = computeHost(host, config, "lll-43.lll", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "lll-43.lll")
		})
		Convey("Expand variables using environment", func() {
			host := config.Hosts["bbb"]
			So(host.HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(host.Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(host.IdentityFile[0], ShouldEqual, "${NON_EXISTING_ENV_VAR}")
			So(host.LocalCommand, ShouldEqual, "${ENV_VAR_LOCALCOMMAND:-hello}")
			So(host.User, ShouldEqual, "user-$ENV_VAR_USER-user")

			So(os.Setenv("ENV_VAR_HOSTNAME", "aaa"), ShouldBeNil)
			So(os.Setenv("ENV_VAR_PORT", "42"), ShouldBeNil)
			So(os.Unsetenv("NON_EXISTING_ENV_VAR"), ShouldBeNil)
			//os.Setenv("ENV_VAR_LOCALCOMMAND", "bbb")
			So(os.Setenv("ENV_VAR_USER", "ccc"), ShouldBeNil)

			computed, err := computeHost(host, config, "bbb", true)
			So(err, ShouldBeNil)

			So(computed.HostName, ShouldEqual, "aaa")
			So(computed.Port, ShouldEqual, "42")
			So(len(computed.IdentityFile), ShouldEqual, 1)
			So(computed.IdentityFile[0], ShouldEqual, "")
			So(computed.LocalCommand, ShouldEqual, "") // FIXME: it should be "hello"
			So(computed.User, ShouldEqual, "user-ccc-user")
		})
	})
}

func TestConfig_getHostByName(t *testing.T) {
	Convey("Testing Config.getHostByName", t, func() {
		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.getHostByName("titi", false, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")

			host, err = config.getHostByName("titi", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")

			host, err = config.getHostByName("dontexists", false, true, false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("dontexists", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists")

			host, err = config.getHostByName("regex.ddd", false, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")

			host, err = config.getHostByName("regex.ddd", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
		})

		Convey("With gateway", func() {
			host, err = config.getHostByName("titi/gateway", false, true, false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("titi/gateway", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByName("dontexists/gateway", false, true, false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("dontexists/gateway", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByName("regex.ddd/gateway", false, true, false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByName("regex.ddd/gateway", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd/gateway")
			So(host.HostName, ShouldNotEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})
	})
}

func TestConfig_GetGatewaySafe(t *testing.T) {
	Convey("Testing Config.GetGatewaySafe", t, func() {

		config := dummyConfig()
		var host *Host

		Convey("Without gateway", func() {
			host = config.GetGatewaySafe("titi")
			So(host.Name(), ShouldEqual, "titi")

			host = config.GetGatewaySafe("dontexists")
			So(host.Name(), ShouldEqual, "dontexists")

			host = config.GetGatewaySafe("regex.ddd")
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
		})

		Convey("With gateway", func() {
			host = config.GetGatewaySafe("titi/gateway")
			So(host.Name(), ShouldEqual, "titi/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetGatewaySafe("dontexists/gateway")
			So(host.Name(), ShouldEqual, "dontexists/gateway")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetGatewaySafe("regex.ddd/gateway")
			So(host.Name(), ShouldEqual, "regex.ddd/gateway")
			So(host.HostName, ShouldNotEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})
	})
}

func TestConfig_needsARebuildForTarget(t *testing.T) {
	Convey("Testing Config.needsARebuildForTarget", t, func() {
		config := dummyConfig()

		So(config.needsARebuildForTarget("totototo"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto0toto"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto1toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto2toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto3toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto4toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto5toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto6toto"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto7toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto8toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto9toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto10toto"), ShouldBeFalse)

		config.addKnownHost("toto1toto")
		config.addKnownHost("toto2toto")

		So(config.needsARebuildForTarget("totototo"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto0toto"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto1toto"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto2toto"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto3toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto4toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto5toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto6toto"), ShouldBeFalse)
		So(config.needsARebuildForTarget("toto7toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto8toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto9toto"), ShouldBeTrue)
		So(config.needsARebuildForTarget("toto10toto"), ShouldBeFalse)
	})
}

func TestConfig_LoadFiles(t *testing.T) {
	Convey("Testing Config.LoadFiles", t, func() {
		config := New()
		file, err := ioutil.TempFile(os.TempDir(), "assh-tests")
		So(err, ShouldBeNil)
		defer func() {
			So(os.Remove(file.Name()), ShouldBeNil)
		}()
		_, err = file.Write([]byte(yamlConfig))
		So(err, ShouldBeNil)

		Convey("Loading a simple file", func() {
			err = config.LoadFiles(file.Name())
			So(err, ShouldBeNil)
			So(config.IncludedFiles(), ShouldResemble, []string{file.Name()})
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 1)
			So(len(config.Hosts), ShouldEqual, 17)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, "")
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(config.Hosts["bbb"].Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(config.Hosts["bbb"].User, ShouldEqual, "user-$ENV_VAR_USER-user")
			So(config.Hosts["bbb"].IdentityFile[0], ShouldEqual, "${NON_EXISTING_ENV_VAR}")
			So(config.Hosts["bbb"].LocalCommand, ShouldEqual, "${ENV_VAR_LOCALCOMMAND:-hello}")
			So(config.Hosts["ccc"].HostName, ShouldEqual, "5.6.7.8")
			So(config.Hosts["ccc"].Port, ShouldEqual, "24")
			So(config.Hosts["ccc"].User, ShouldEqual, "toor")
			So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")
			So(config.Hosts["*.ddd"].Port, ShouldEqual, "")
			So(config.Hosts["*.ddd"].User, ShouldEqual, "")
			So(config.Defaults.Port, ShouldEqual, "22")
			So(config.Defaults.User, ShouldEqual, "root")
			So(len(config.Templates), ShouldEqual, 3)
			So(config.Templates["kkk"].Port, ShouldEqual, "25")
			So(config.Templates["kkk"].User, ShouldEqual, "kkkk")
		})
		Convey("Loading the same file again", func() {
			So(config.LoadFiles(file.Name()), ShouldBeNil)
			err = config.LoadFiles(file.Name())
			So(err, ShouldBeNil)
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 1)
			So(len(config.Hosts), ShouldEqual, 17)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, "")
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(config.Hosts["bbb"].Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(config.Hosts["bbb"].User, ShouldEqual, "user-$ENV_VAR_USER-user")
			So(config.Hosts["bbb"].IdentityFile[0], ShouldEqual, "${NON_EXISTING_ENV_VAR}")
			So(config.Hosts["bbb"].LocalCommand, ShouldEqual, "${ENV_VAR_LOCALCOMMAND:-hello}")
			So(config.Hosts["ccc"].HostName, ShouldEqual, "5.6.7.8")
			So(config.Hosts["ccc"].Port, ShouldEqual, "24")
			So(config.Hosts["ccc"].User, ShouldEqual, "toor")
			So(config.Hosts["*.ddd"].HostName, ShouldEqual, "1.3.5.7")
			So(config.Hosts["*.ddd"].Port, ShouldEqual, "")
			So(config.Hosts["*.ddd"].User, ShouldEqual, "")
			So(config.Defaults.Port, ShouldEqual, "22")
			So(config.Defaults.User, ShouldEqual, "root")
			So(len(config.Templates), ShouldEqual, 3)
			So(config.Templates["kkk"].Port, ShouldEqual, "25")
			So(config.Templates["kkk"].User, ShouldEqual, "kkkk")
		})
		Convey("Expand includes environment", func() {
			config := New()
			file, err := ioutil.TempFile(os.TempDir(), "assh-tests")
			So(err, ShouldBeNil)
			defer func() { So(os.Remove(file.Name()), ShouldBeNil) }()
			_, err = file.Write([]byte(`
includes:
- $DUMMY_ENV_VAR/assh-tests*
`))
			So(err, ShouldBeNil)
			tempDir, err := ioutil.TempDir(os.TempDir(), "assh-tests")
			So(err, ShouldBeNil)
			defer func() { So(os.RemoveAll(tempDir), ShouldBeNil) }()

			file2, err := ioutil.TempFile(tempDir, "assh-tests")
			So(err, ShouldBeNil)
			defer func() { So(os.Remove(file2.Name()), ShouldBeNil) }()
			So(os.Setenv("DUMMY_ENV_VAR", tempDir), ShouldBeNil)

			So(config.LoadFiles(file.Name()), ShouldBeNil)

			So(err, ShouldBeNil)
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(config.includedFiles[file2.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 2)

		})

	})
	// FIXME: test globbing
}

func TestConfig_getHostByPath(t *testing.T) {
	Convey("Testing Config.getHostByPath", t, func() {

		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.getHostByPath("titi", false, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("titi", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("dontexists", false, true, false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByPath("dontexists", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("regex.ddd", false, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.getHostByPath("regex.ddd", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host, err = config.getHostByPath("titi/gateway", false, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("titi/gateway", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("dontexists/gateway", false, true, false)
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.getHostByPath("dontexists/gateway", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("regex.ddd/gateway", false, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.getHostByPath("regex.ddd/gateway", true, true, false)
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})
	})
}

func TestConfig_GetHost(t *testing.T) {
	Convey("Testing Config.GetHost", t, func() {
		config := dummyConfig()
		var host *Host
		var err error

		Convey("Without gateway", func() {
			host, err = config.GetHost("titi")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host, err = config.GetHost("dontexists")
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			host, err = config.GetHost("regex.ddd")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host, err = config.GetHost("titi/gateway")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host, err = config.GetHost("dontexists/gateway")
			So(err, ShouldNotBeNil)
			So(host, ShouldBeNil)

			// FIXME: check if this is a normal behavior
			host, err = config.GetHost("regex.ddd/gateway")
			So(err, ShouldBeNil)
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})

		Convey("Inheritance", FailureContinues, func() {
			host, err = config.GetHost("tata")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"tata": true,
				"tutu": true,
				"titi": true,
				"toto": true,
			})
			So(host.ProxyCommand, ShouldEqual, "nc -v 4242")
			So(host.User, ShouldEqual, "moul")
			So(host.Gateways, ShouldResemble, composeyaml.Stringorslice{"titi", "direct", "1.2.3.4"})
			So(host.PasswordAuthentication, ShouldEqual, "yes")

			host, err = config.GetHost("tutu")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"tutu":  true,
				"toto":  true,
				"*.ddd": true,
			})
			So(host.User, ShouldEqual, "root")
			So(host.Gateways, ShouldResemble, composeyaml.Stringorslice{"titi", "direct", "1.2.3.4"})
			So(host.PasswordAuthentication, ShouldEqual, "yes")

			host, err = config.GetHost("nnn")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"nnn": true,
				"mmm": true,
			})
			So(host.User, ShouldEqual, "mmmm")
			So(host.Port, ShouldEqual, "26")
			So(host.Gateways, ShouldResemble, composeyaml.Stringorslice{"titi", "direct", "1.2.3.4"})
		})

		Convey("Aliases", FailureContinues, func() {
			host, err = config.GetHost("ooo1")
			So(err, ShouldBeNil)
			So(host.name, ShouldEqual, "ooo1")
			So(host.Aliases, ShouldResemble, composeyaml.Stringorslice{
				"ooo11",
				"ooo12",
			})
			So(host.Port, ShouldEqual, "23")

			host, err = config.GetHost("ooo2")
			So(err, ShouldBeNil)
			So(host.name, ShouldEqual, "ooo2")
			So(host.Aliases, ShouldResemble, composeyaml.Stringorslice{
				"ooo21",
				"ooo22",
			})
			So(host.Port, ShouldEqual, "24")

			host, err = config.GetHost("ooo11")
			So(err, ShouldBeNil)
			So(host.name, ShouldEqual, "ooo11")
			So(host.Aliases, ShouldResemble, composeyaml.Stringorslice{
				"ooo11",
				"ooo12",
			})
			So(host.Port, ShouldEqual, "23")

			host, err = config.GetHost("ooo22")
			So(err, ShouldBeNil)
			So(host.name, ShouldEqual, "ooo22")
			So(host.Aliases, ShouldResemble, composeyaml.Stringorslice{
				"ooo21",
				"ooo22",
			})
			So(host.Port, ShouldEqual, "24")

		})
	})
}

func TestConfig_GetHostSafe(t *testing.T) {
	Convey("Testing Config.GetHostSafe", t, func() {

		config := dummyConfig()
		var host *Host

		Convey("Without gateway", func() {
			host = config.GetHostSafe("titi")
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetHostSafe("dontexists")
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 0)

			host = config.GetHostSafe("regex.ddd")
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 0)
		})

		Convey("With gateway", func() {
			host = config.GetHostSafe("titi/gateway")
			So(host.Name(), ShouldEqual, "titi")
			So(len(host.Gateways), ShouldEqual, 1)

			host = config.GetHostSafe("dontexists/gateway")
			So(host.Name(), ShouldEqual, "dontexists")
			So(len(host.Gateways), ShouldEqual, 1)

			host = config.GetHostSafe("regex.ddd/gateway")
			So(host.Name(), ShouldEqual, "regex.ddd")
			So(host.HostName, ShouldEqual, "1.3.5.7")
			So(len(host.Gateways), ShouldEqual, 1)
		})
	})
}

func TestConfig_String(t *testing.T) {
	Convey("Testing Config.String", t, func() {
		config := dummyConfig()
		So(config.String(), ShouldEqual, `{"hosts":{"*.ddd":{"PasswordAuthentication":"yes","HostName":"1.3.5.7"},"empty":{},"nnn":{"Port":"26","Inherits":["mmm"]},"ooo1":{"Port":"23","Aliases":["ooo11","ooo12"]},"ooo2":{"Port":"24","Aliases":["ooo21","ooo22"]},"tata":{"Inherits":["tutu","titi","toto","tutu"]},"titi":{"Port":"23","User":"moul","HostName":"tata","ProxyCommand":"nc -v 4242","ControlMasterMkdir":"true","Comment":["Hello World"]},"tonton":{"ResolveNameservers":["a.com","1.2.3.4"],"Comment":["AAA","BBB"]},"toto":{"HostName":"1.2.3.4"},"toto[1-5]toto":{"User":"toto1"},"toto[7-9]toto":{"User":"toto2"},"toutou":{"RemoteCommand":"date \u003e\u003e /tmp/logs","ResolveCommand":"dig -t %h","Comment":["First line Second line Third line\n"]},"tutu":{"Inherits":["toto","tutu","*.ddd"],"Gateways":["titi","direct","1.2.3.4"]},"zzz":{"AddressFamily":"any","AskPassGUI":"yes","BatchMode":"no","CanonicalDomains":"42.am","CanonicalizeFallbackLocal":"no","CanonicalizeHostname":"yes","CanonicalizeMaxDots":"1","CanonicalizePermittedCNAMEs":"*.a.example.com:*.b.example.com:*.c.example.com","ChallengeResponseAuthentication":"yes","CheckHostIP":"yes","Cipher":"blowfish","Ciphers":["aes128-ctr,aes192-ctr","aes256-ctr"],"ClearAllForwardings":"yes","Compression":"yes","CompressionLevel":6,"ConnectionAttempts":"1","ConnectTimeout":10,"ControlMaster":"yes","ControlPath":"/tmp/%L-%l-%n-%p-%u-%r-%C-%h","ControlPersist":"yes","DynamicForward":["0.0.0.0:4242","0.0.0.0:4343"],"EnableSSHKeysign":"yes","EscapeChar":"~","ExitOnForwardFailure":"yes","FingerprintHash":"sha256","ForwardAgent":"yes","ForwardX11":"yes","ForwardX11Timeout":42,"ForwardX11Trusted":"yes","GatewayPorts":"yes","GlobalKnownHostsFile":["/etc/ssh/ssh_known_hosts","/tmp/ssh_known_hosts"],"GSSAPIAuthentication":"no","GSSAPIClientIdentity":"moul","GSSAPIDelegateCredentials":"no","GSSAPIKeyExchange":"no","GSSAPIRenewalForcesRekey":"no","GSSAPIServerIdentity":"gssapi.example.com","GSSAPITrustDNS":"no","HashKnownHosts":"no","HostbasedAuthentication":"no","HostbasedKeyTypes":"*","HostKeyAlgorithms":["ecdsa-sha2-nistp256-cert-v01@openssh.com","test"],"HostKeyAlias":"z","IdentitiesOnly":"yes","IdentityFile":["~/.ssh/identity","~/.ssh/identity2"],"IgnoreUnknown":"testtest","IPQoS":["lowdelay","highdelay"],"KbdInteractiveAuthentication":"yes","KbdInteractiveDevices":["bsdauth","test"],"KexAlgorithms":["curve25519-sha256@libssh.org","test"],"KeychainIntegration":"yes","LocalCommand":"echo %h \u003e /tmp/logs","LocalForward":["0.0.0.0:1234","0.0.0.0:1235"],"LogLevel":"DEBUG3","MACs":["umac-64-etm@openssh.com,umac-128-etm@openssh.com","test"],"Match":"all","NoHostAuthenticationForLocalhost":"yes","NumberOfPasswordPrompts":"3","PasswordAuthentication":"yes","PermitLocalCommand":"yes","PKCS11Provider":"/a/b/c/pkcs11.so","Port":"22","PreferredAuthentications":"gssapi-with-mic,hostbased,publickey","Protocol":["2","3"],"ProxyUseFdpass":"no","PubkeyAuthentication":"yes","RekeyLimit":"default none","RemoteForward":["0.0.0.0:1234","0.0.0.0:1255"],"RequestTTY":"yes","RevokedHostKeys":"/a/revoked-keys","RhostsRSAAuthentication":"no","RSAAuthentication":"yes","SendEnv":["CUSTOM_*,TEST","TEST2"],"ServerAliveCountMax":3,"StreamLocalBindMask":"0177","StreamLocalBindUnlink":"no","StrictHostKeyChecking":"ask","TCPKeepAlive":"yes","Tunnel":"yes","TunnelDevice":"any:any","UpdateHostKeys":"ask","UseKeychain":"no","UsePrivilegedPort":"no","User":"moul","UserKnownHostsFile":["~/.ssh/known_hosts ~/.ssh/known_hosts2","/tmp/known_hosts"],"VerifyHostKeyDNS":"no","VisualHostKey":"yes","XAuthLocation":"xauth","HostName":"zzz.com","ProxyCommand":"nc %h %p"}},"templates":{"mmm":{"Port":"25","User":"mmmm","HostName":"5.5.5.5","Inherits":["tata"]}},"defaults":{"Port":"22","User":"root","Hooks":{}},"asshknownhostfile":"~/.ssh/assh_known_hosts"}`)
	})
}

func TestConfig_WriteSSHConfig(t *testing.T) {
	Convey("Testing Config.WriteSSHConfig", t, func() {
		config := dummyConfig()

		config.addKnownHost("toto1toto")
		config.addKnownHost("toto2toto")
		config.addKnownHost("toto7toto")

		var buffer bytes.Buffer

		config.ASSHBinaryPath = "assh"
		asshBinaryPath = "assh"

		err := config.WriteSSHConfigTo(&buffer)
		So(err, ShouldBeNil)

		expected := `# more info: https://github.com/moul/assh

# host-based configuration
Host *.ddd
  PasswordAuthentication yes
  # HostName: 1.3.5.7

Host empty

Host nnn
  PasswordAuthentication yes
  Port 26
  User mmmm
  # ProxyCommand nc -v 4242
  # HostName: 5.5.5.5
  # ControlMasterMkdir: true
  # Inherits: [mmm]
  # Gateways: [titi, direct, 1.2.3.4]
  # Comment: [Hello World]

Host ooo1
  Port 23
  # Aliases: [ooo11, ooo12]

Host ooo11
  Port 23
  # AliasOf: ooo1

Host ooo12
  Port 23
  # AliasOf: ooo1

Host ooo2
  Port 24
  # Aliases: [ooo21, ooo22]

Host ooo21
  Port 24
  # AliasOf: ooo2

Host ooo22
  Port 24
  # AliasOf: ooo2

Host tata
  PasswordAuthentication yes
  Port 22
  User moul
  # ProxyCommand nc -v 4242
  # HostName: 1.2.3.4
  # ControlMasterMkdir: true
  # Inherits: [tutu, titi, toto, tutu]
  # Gateways: [titi, direct, 1.2.3.4]
  # Comment: [Hello World]

Host titi
  Port 23
  User moul
  # ProxyCommand nc -v 4242
  # HostName: tata
  # ControlMasterMkdir: true
  # Comment: [Hello World]

Host tonton
  # Comment: [AAA, BBB]
  # ResolveNameservers: [a.com, 1.2.3.4]

Host toto
  # HostName: 1.2.3.4

Host toto[1-5]toto
  User toto1
  # KnownHosts: [toto1toto, toto2toto]

Host toto1toto
  User toto1
  # KnownHostOf: toto[1-5]toto

Host toto2toto
  User toto1
  # KnownHostOf: toto[1-5]toto

Host toto[7-9]toto
  User toto2
  # KnownHosts: [toto7toto]

Host toto7toto
  User toto2
  # KnownHostOf: toto[7-9]toto

Host toutou
  RemoteCommand date >> /tmp/logs
  # Comment: [First line Second line Third line]
  # ResolveCommand: dig -t %h

Host tutu
  PasswordAuthentication yes
  Port 22
  # HostName: 1.2.3.4
  # Inherits: [toto, tutu, *.ddd]
  # Gateways: [titi, direct, 1.2.3.4]

Host zzz
  AddressFamily any
  AskPassGUI yes
  BatchMode no
  CanonicalDomains 42.am
  CanonicalizeFallbackLocal no
  CanonicalizeHostname yes
  CanonicalizeMaxDots 1
  CanonicalizePermittedCNAMEs *.a.example.com:*.b.example.com:*.c.example.com
  ChallengeResponseAuthentication yes
  CheckHostIP yes
  Cipher blowfish
  Ciphers aes128-ctr,aes192-ctr,aes256-ctr
  ClearAllForwardings yes
  Compression yes
  CompressionLevel 6
  ConnectionAttempts 1
  ConnectTimeout 10
  ControlMaster yes
  ControlPath /tmp/%L-%l-%n-%p-%u-%r-%C-%h
  ControlPersist yes
  DynamicForward 0.0.0.0:4242
  DynamicForward 0.0.0.0:4343
  EnableSSHKeysign yes
  EscapeChar ~
  ExitOnForwardFailure yes
  FingerprintHash sha256
  ForwardAgent yes
  ForwardX11 yes
  ForwardX11Timeout 42
  ForwardX11Trusted yes
  GatewayPorts yes
  GlobalKnownHostsFile /etc/ssh/ssh_known_hosts /tmp/ssh_known_hosts
  GSSAPIAuthentication no
  GSSAPIClientIdentity moul
  GSSAPIDelegateCredentials no
  GSSAPIKeyExchange no
  GSSAPIRenewalForcesRekey no
  GSSAPIServerIdentity gssapi.example.com
  GSSAPITrustDNS no
  HashKnownHosts no
  HostbasedAuthentication no
  HostbasedKeyTypes *
  HostKeyAlgorithms ecdsa-sha2-nistp256-cert-v01@openssh.com,test
  HostKeyAlias z
  IdentitiesOnly yes
  IdentityFile ~/.ssh/identity
  IdentityFile ~/.ssh/identity2
  IgnoreUnknown testtest
  IPQoS lowdelay highdelay
  KbdInteractiveAuthentication yes
  KbdInteractiveDevices bsdauth,test
  KexAlgorithms curve25519-sha256@libssh.org,test
  KeychainIntegration yes
  LocalCommand echo %h > /tmp/logs
  LocalForward 0.0.0.0:1234
  LocalForward 0.0.0.0:1235
  LogLevel DEBUG3
  MACs umac-64-etm@openssh.com,umac-128-etm@openssh.com,test
  Match all
  NoHostAuthenticationForLocalhost yes
  NumberOfPasswordPrompts 3
  PasswordAuthentication yes
  PermitLocalCommand yes
  PKCS11Provider /a/b/c/pkcs11.so
  Port 22
  PreferredAuthentications gssapi-with-mic,hostbased,publickey
  Protocol 2,3
  ProxyUseFdpass no
  PubkeyAuthentication yes
  RekeyLimit default none
  RemoteForward 0.0.0.0:1234
  RemoteForward 0.0.0.0:1255
  RequestTTY yes
  RevokedHostKeys /a/revoked-keys
  RhostsRSAAuthentication no
  RSAAuthentication yes
  SendEnv CUSTOM_*,TEST
  SendEnv TEST2
  ServerAliveCountMax 3
  StreamLocalBindMask 0177
  StreamLocalBindUnlink no
  StrictHostKeyChecking ask
  TCPKeepAlive yes
  Tunnel yes
  TunnelDevice any:any
  UpdateHostKeys ask
  UseKeychain no
  UsePrivilegedPort no
  User moul
  UserKnownHostsFile ~/.ssh/known_hosts ~/.ssh/known_hosts2 /tmp/known_hosts
  VerifyHostKeyDNS no
  VisualHostKey yes
  XAuthLocation xauth
  # ProxyCommand nc %h %p
  # HostName: zzz.com

# global configuration
Host *
  Port 22
  User root
  ProxyCommand assh connect --port=%p %h
`
		output := strings.Join(strings.Split(buffer.String(), "\n")[3:], "\n")
		So(output, ShouldEqual, expected)
	})
	Convey("Testing very long string comment", t, func() {
		config := New()
		config.Hosts["toto"] = &Host{
			ResolveCommand: `0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222333333333344444444445555555555666666666677777777778888888888999999999900000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222333333333344444444445555555555666666666677777777778888888888999999999900000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222333333333344444444445555555555666666666677777777778888888888999999999900000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222`,
		}
		var buffer bytes.Buffer
		config.ASSHBinaryPath = "assh"
		asshBinaryPath = "assh"
		err := config.WriteSSHConfigTo(&buffer)
		So(err, ShouldBeNil)
		expected := `# more info: https://github.com/moul/assh

# host-based configuration
Host toto
  # ResolveCommand: 00000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222333333333344444444445555555555666666666677777777778888888888999999999900000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222333333333344444444445555555555666666666677777777778888888888999999999900000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999000000000011111111112222222222333333333344444444445555555555666666666677777777778888888888999999999900000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990
  # ResolveCommand: 00000000011111111112222222222

# global configuration
Host *
`
		output := strings.Join(strings.Split(buffer.String(), "\n")[3:], "\n")
		So(output, ShouldEqual, expected)
	})
	Convey("Testing very long slice comment", t, func() {
		config := New()
		config.Hosts["toto"] = &Host{
			Gateways: []string{},
		}
		for i := 0; i < 10; i++ {
			config.Hosts["toto"].Gateways = append(config.Hosts["toto"].Gateways, []string{
				"0000000000",
				"1111111111",
				"2222222222",
				"3333333333",
				"4444444444",
				"5555555555",
				"6666666666",
				"7777777777",
				"8888888888",
				"9999999999",
			}...)
		}
		var buffer bytes.Buffer
		config.ASSHBinaryPath = "assh"
		asshBinaryPath = "assh"
		err := config.WriteSSHConfigTo(&buffer)
		So(err, ShouldBeNil)
		expected := `# more info: https://github.com/moul/assh

# host-based configuration
Host toto
  # Gateways: [0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222]
  # Gateways: [3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999, 0000000000, 1111111111, 2222222222, 3333333333, 4444444444, 5555555555, 6666666666, 7777777777, 8888888888, 9999999999]

# global configuration
Host *
`
		output := strings.Join(strings.Split(buffer.String(), "\n")[3:], "\n")
		So(output, ShouldEqual, expected)
	})
}

func TestConfig_ValidateSummary(t *testing.T) {
	Convey("Testing Config.ValidateSummary", t, FailureContinues, func() {
		// no error
		config := New()
		config.Hosts["toto"] = &Host{name: "toto"}
		err := config.ValidateSummary()
		So(err, ShouldBeNil)

		// one error
		config.Hosts["toto"].ControlMaster = "invalid data"
		err = config.ValidateSummary()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `"toto": invalid value for 'ControlMaster': "invalid data"`)

		// multiple errors
		config.Hosts["toto"].AddressFamily = "invalid data"
		config.Hosts["tata"] = &Host{name: "tata"}
		config.Hosts["tata"].AddressFamily = "invalid data"
		errs := config.Validate()
		So(len(errs), ShouldEqual, 3)
	})
}

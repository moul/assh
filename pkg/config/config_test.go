package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	yamlConfig = `
hosts:

  aaa:
    HostName: 1.2.3.4

  bbb:
    Port: ${ENV_VAR_PORT}
    HostName: $ENV_VAR_HOSTNAME
    User: user-$ENV_VAR_USER-user
    LocalCommand: ${ENV_VAR_LOCALCOMMAND:-hello}
    IdentityFile: ${NON_EXISTING_ENV_VAR}

  ccc:
    HostName: 5.6.7.8
    Port: 24
    User: toor
  "*.ddd":
    HostName: 1.3.5.7

  eee:
    Inherits:
    - aaa
    - bbb
    - aaa
    NoControlMasterMkdir: true

  fff:
    Inherits:
    - bbb
    - eee
    - "*.ddd"

  ggg:
    Gateways:
    - direct
    - fff

  hhh:
    Gateways:
    - ggg
    - direct

  iii:
    Gateways:
    - test.ddd

  jjj:
    HostName: "%h.jjjjj"

  "*.kkk":
    HostName: "%h.kkkkk"

  "lll-*":
    HostName: "%h.lll"

  nnn:
    Inherits:
    - mmm
    User: nnnn

templates:

  kkk:
    Port: 25
    User: kkkk

  lll:
    HostName: 5.5.5.5

  mmm:
    Inherits:
    - iii

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
	config.Hosts["toto"] = Host{
		HostName: "1.2.3.4",
	}
	config.Hosts["titi"] = Host{
		HostName:             "tata",
		Port:                 "23",
		User:                 "moul",
		ProxyCommand:         "nc -v 4242",
		NoControlMasterMkdir: "true",
	}
	config.Hosts["tonton"] = Host{
		ResolveNameservers: []string{"a.com", "1.2.3.4"},
	}
	config.Hosts["toutou"] = Host{
		ResolveCommand: "dig -t %h",
	}
	config.Hosts["tutu"] = Host{
		Gateways: []string{"titi", "direct", "1.2.3.4"},
		Inherits: []string{"toto", "tutu", "*.ddd"},
	}
	config.Hosts["empty"] = Host{}
	config.Hosts["tata"] = Host{
		Inherits: []string{"tutu", "titi", "toto", "tutu"},
	}
	config.Hosts["*.ddd"] = Host{
		HostName:               "1.3.5.7",
		PasswordAuthentication: "yes",
	}
	config.Defaults = Host{
		Port: "22",
		User: "root",
	}
	config.Templates["mmm"] = Host{
		Port:     "25",
		User:     "mmmm",
		HostName: "5.5.5.5",
		Inherits: []string{"tata"},
	}
	config.Hosts["nnn"] = Host{
		Port:     "26",
		Inherits: []string{"mmm"},
	}
	config.Hosts["zzz"] = Host{
		// ssh-config fields
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
		Ciphers:                         "aes128-ctr,aes192-ctr,aes256-ctr",
		ClearAllForwardings:             "yes",
		Compression:                     "yes",
		CompressionLevel:                6,
		ConnectionAttempts:              "1",
		ConnectTimeout:                  10,
		ControlMaster:                   "yes",
		ControlPath:                     "/tmp/%L-%l-%n-%p-%u-%r-%C-%h",
		ControlPersist:                  "yes",
		DynamicForward:                  "0.0.0.0:4242",
		EnableSSHKeysign:                "yes",
		EscapeChar:                      "~",
		ExitOnForwardFailure:            "yes",
		FingerprintHash:                 "sha256",
		ForwardAgent:                    "yes",
		ForwardX11:                      "yes",
		ForwardX11Timeout:               42,
		ForwardX11Trusted:               "yes",
		GatewayPorts:                    "yes",
		GlobalKnownHostsFile:            "/etc/ssh/ssh_known_hosts",
		GSSAPIAuthentication:            "no",
		GSSAPIKeyExchange:               "no",
		GSSAPIClientIdentity:            "moul",
		GSSAPIServerIdentity:            "gssapi.example.com",
		GSSAPIDelegateCredentials:       "no",
		GSSAPIRenewalForcesRekey:        "no",
		GSSAPITrustDns:                  "no",
		HashKnownHosts:                  "no",
		HostbasedAuthentication:         "no",
		HostbasedKeyTypes:               "*",
		HostKeyAlgorithms:               "ecdsa-sha2-nistp256-cert-v01@openssh.com",
		HostKeyAlias:                    "z",
		IdentitiesOnly:                  "yes",
		IdentityFile:                    "~/.ssh/identity",
		IgnoreUnknown:                   "testtest", // FIXME: looks very interesting to generate .ssh/config without comments !
		IPQoS:                           "lowdelay",
		KbdInteractiveAuthentication: "yes",
		KbdInteractiveDevices:        "bsdauth",
		KeychainIntegration:          "yes",
		KexAlgorithms:                "curve25519-sha256@libssh.org", // for all algorithms/ciphers, we could have an "assh diagnose" that warns about unsafe connections
		LocalCommand:                 "echo %h > /tmp/logs",
		LocalForward:                 "0.0.0.0:1234", // FIXME: may be a list
		LogLevel:                     "DEBUG3",
		MACs:                         "umac-64-etm@openssh.com,umac-128-etm@openssh.com",
		Match:                        "all",
		NoHostAuthenticationForLocalhost: "yes",
		NumberOfPasswordPrompts:          "3",
		PasswordAuthentication:           "yes",
		PermitLocalCommand:               "yes",
		PKCS11Provider:                   "/a/b/c/pkcs11.so",
		Port:                             "22",
		PreferredAuthentications: "gssapi-with-mic,hostbased,publickey",
		Protocol:                 "2",
		ProxyUseFdpass:           "no",
		PubkeyAuthentication:     "yes",
		RekeyLimit:               "default none",
		RemoteForward:            "0.0.0.0:1234",
		RequestTTY:               "yes",
		RevokedHostKeys:          "/a/revoked-keys",
		RhostsRSAAuthentication:  "no",
		RSAAuthentication:        "yes",
		SendEnv:                  "CUSTOM_*,TEST",
		ServerAliveCountMax:      3,
		ServerAliveInterval:      0,
		StreamLocalBindMask:      "0177",
		StreamLocalBindUnlink:    "no",
		StrictHostKeyChecking:    "ask",
		TCPKeepAlive:             "yes",
		Tunnel:                   "yes",
		TunnelDevice:             "any:any",
		UpdateHostKeys:           "ask",
		UsePrivilegedPort:        "no",
		User:                     "moul",
		UserKnownHostsFile:       "~/.ssh/known_hosts ~/.ssh/known_hosts2",
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
		ResolveNameservers: []string{},
		ResolveCommand:     "",
	}
	config.applyMissingNames()
	return config
}

func TestConfig(t *testing.T) {
	Convey("Testing dummyConfig", t, func() {
		config := dummyConfig()

		So(len(config.Hosts), ShouldEqual, 10)

		So(config.Hosts["toto"].HostName, ShouldEqual, "1.2.3.4")
		So(config.Hosts["toto"].Port, ShouldEqual, "")
		So(config.Hosts["toto"].name, ShouldEqual, "toto")
		So(config.Hosts["toto"].isDefault, ShouldEqual, false)

		So(config.Hosts["titi"].HostName, ShouldEqual, "tata")
		So(config.Hosts["titi"].User, ShouldEqual, "moul")
		So(config.Hosts["titi"].ProxyCommand, ShouldEqual, "nc -v 4242")
		So(BoolVal(config.Hosts["titi"].NoControlMasterMkdir), ShouldBeTrue)
		So(config.Hosts["titi"].Port, ShouldEqual, "23")
		So(config.Hosts["titi"].isDefault, ShouldEqual, false)

		So(config.Hosts["tonton"].isDefault, ShouldEqual, false)
		So(config.Hosts["tonton"].Port, ShouldEqual, "")
		So(config.Hosts["tonton"].ResolveNameservers, ShouldResemble, []string{"a.com", "1.2.3.4"})

		So(config.Hosts["toutou"].isDefault, ShouldEqual, false)
		So(config.Hosts["toutou"].Port, ShouldEqual, "")
		So(config.Hosts["toutou"].ResolveCommand, ShouldEqual, "dig -t %h")

		So(config.Hosts["tutu"].isDefault, ShouldEqual, false)
		So(config.Hosts["tutu"].Port, ShouldEqual, "")
		So(config.Hosts["tutu"].Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})

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
			So(len(config.Hosts), ShouldEqual, 13)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, "")
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(config.Hosts["bbb"].Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(config.Hosts["bbb"].User, ShouldEqual, "user-$ENV_VAR_USER-user")
			So(config.Hosts["bbb"].IdentityFile, ShouldEqual, "${NON_EXISTING_ENV_VAR}")
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

func TestConfig_JsonString(t *testing.T) {
	Convey("Testing Config.JsonString", t, func() {
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
      "NoControlMasterMkdir": "true"
    },
    "tonton": {
      "ResolveNameservers": [
        "a.com",
        "1.2.3.4"
      ]
    },
    "toto": {
      "HostName": "1.2.3.4"
    },
    "toutou": {
      "ResolveCommand": "dig -t %h"
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
      "Ciphers": "aes128-ctr,aes192-ctr,aes256-ctr",
      "ClearAllForwardings": "yes",
      "Compression": "yes",
      "CompressionLevel": 6,
      "ConnectionAttempts": "1",
      "ConnectTimeout": 10,
      "ControlMaster": "yes",
      "ControlPath": "/tmp/%L-%l-%n-%p-%u-%r-%C-%h",
      "ControlPersist": "yes",
      "DynamicForward": "0.0.0.0:4242",
      "EnableSSHKeysign": "yes",
      "EscapeChar": "~",
      "ExitOnForwardFailure": "yes",
      "FingerprintHash": "sha256",
      "ForwardAgent": "yes",
      "ForwardX11": "yes",
      "ForwardX11Timeout": 42,
      "ForwardX11Trusted": "yes",
      "GatewayPorts": "yes",
      "GlobalKnownHostsFile": "/etc/ssh/ssh_known_hosts",
      "GSSAPIAuthentication": "no",
      "GSSAPIClientIdentity": "moul",
      "GSSAPIDelegateCredentials": "no",
      "GSSAPIKeyExchange": "no",
      "GSSAPIRenewalForcesRekey": "no",
      "GSSAPIServerIdentity": "gssapi.example.com",
      "GSSAPITrustDns": "no",
      "HashKnownHosts": "no",
      "HostbasedAuthentication": "no",
      "HostbasedKeyTypes": "*",
      "HostKeyAlgorithms": "ecdsa-sha2-nistp256-cert-v01@openssh.com",
      "HostKeyAlias": "z",
      "IdentitiesOnly": "yes",
      "IdentityFile": "~/.ssh/identity",
      "IgnoreUnknown": "testtest",
      "IPQoS": "lowdelay",
      "KbdInteractiveAuthentication": "yes",
      "KbdInteractiveDevices": "bsdauth",
      "KexAlgorithms": "curve25519-sha256@libssh.org",
      "KeychainIntegration": "yes",
      "LocalCommand": "echo %h \u003e /tmp/logs",
      "LocalForward": "0.0.0.0:1234",
      "LogLevel": "DEBUG3",
      "MACs": "umac-64-etm@openssh.com,umac-128-etm@openssh.com",
      "Match": "all",
      "NoHostAuthenticationForLocalhost": "yes",
      "NumberOfPasswordPrompts": "3",
      "PasswordAuthentication": "yes",
      "PermitLocalCommand": "yes",
      "PKCS11Provider": "/a/b/c/pkcs11.so",
      "Port": "22",
      "PreferredAuthentications": "gssapi-with-mic,hostbased,publickey",
      "Protocol": "2",
      "ProxyUseFdpass": "no",
      "PubkeyAuthentication": "yes",
      "RekeyLimit": "default none",
      "RemoteForward": "0.0.0.0:1234",
      "RequestTTY": "yes",
      "RevokedHostKeys": "/a/revoked-keys",
      "RhostsRSAAuthentication": "no",
      "RSAAuthentication": "yes",
      "SendEnv": "CUSTOM_*,TEST",
      "ServerAliveCountMax": 3,
      "StreamLocalBindMask": "0177",
      "StreamLocalBindUnlink": "no",
      "StrictHostKeyChecking": "ask",
      "TCPKeepAlive": "yes",
      "Tunnel": "yes",
      "TunnelDevice": "any:any",
      "UpdateHostKeys": "ask",
      "UsePrivilegedPort": "no",
      "User": "moul",
      "UserKnownHostsFile": "~/.ssh/known_hosts ~/.ssh/known_hosts2",
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
    "User": "root"
  }
}`
			json, err := config.JsonString()
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
      "IdentityFile": "${NON_EXISTING_ENV_VAR}",
      "LocalCommand": "${ENV_VAR_LOCALCOMMAND:-hello}",
      "Port": "${ENV_VAR_PORT}",
      "User": "user-$ENV_VAR_USER-user",
      "HostName": "$ENV_VAR_HOSTNAME"
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
      "NoControlMasterMkdir": "true"
    },
    "fff": {
      "Inherits": [
        "bbb",
        "eee",
        "*.ddd"
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
    "User": "root"
  },
  "includes": [
    "/path/to/dir/*.yml",
    "/path/to/file.yml"
  ]
}`
			json, err := config.JsonString()
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
			computed, err := computeHost(&host, config, "jjj", false)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "%h.jjjjj")

			computed, err = computeHost(&host, config, "jjj", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "jjj.jjjjj")

			host = config.Hosts["*.kkk"]
			computed, err = computeHost(&host, config, "test.kkk", false)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "%h.kkkkk")

			computed, err = computeHost(&host, config, "test.kkk", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "test.kkk.kkkkk")
		})
		Convey("Do not expand variables twice", func() {
			host := config.Hosts["lll-*"]
			computed, err := computeHost(&host, config, "lll-42", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "lll-42.lll")

			computed, err = computeHost(&host, config, "lll-43.lll", true)
			So(err, ShouldBeNil)
			So(computed.HostName, ShouldEqual, "lll-43.lll")
		})
		Convey("Expand variables using environment", func() {
			host := config.Hosts["bbb"]
			So(host.HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(host.Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(host.IdentityFile, ShouldEqual, "${NON_EXISTING_ENV_VAR}")
			So(host.LocalCommand, ShouldEqual, "${ENV_VAR_LOCALCOMMAND:-hello}")
			So(host.User, ShouldEqual, "user-$ENV_VAR_USER-user")

			os.Setenv("ENV_VAR_HOSTNAME", "aaa")
			os.Setenv("ENV_VAR_PORT", "42")
			os.Unsetenv("NON_EXISTING_ENV_VAR")
			//os.Setenv("ENV_VAR_LOCALCOMMAND", "bbb")
			os.Setenv("ENV_VAR_USER", "ccc")

			computed, err := computeHost(&host, config, "bbb", true)
			So(err, ShouldBeNil)

			So(computed.HostName, ShouldEqual, "aaa")
			So(computed.Port, ShouldEqual, "42")
			So(computed.IdentityFile, ShouldEqual, "")
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

func TestConfig_LoadFiles(t *testing.T) {
	Convey("Testing Config.LoadFiles", t, func() {
		config := New()
		file, err := ioutil.TempFile(os.TempDir(), "assh-tests")
		So(err, ShouldBeNil)
		defer os.Remove(file.Name())
		file.Write([]byte(yamlConfig))

		Convey("Loading a simple file", func() {
			err = config.LoadFiles(file.Name())

			So(err, ShouldBeNil)
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 1)
			So(len(config.Hosts), ShouldEqual, 13)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, "")
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(config.Hosts["bbb"].Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(config.Hosts["bbb"].User, ShouldEqual, "user-$ENV_VAR_USER-user")
			So(config.Hosts["bbb"].IdentityFile, ShouldEqual, "${NON_EXISTING_ENV_VAR}")
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
			config.LoadFiles(file.Name())
			err = config.LoadFiles(file.Name())

			So(err, ShouldBeNil)
			So(config.includedFiles[file.Name()], ShouldEqual, true)
			So(len(config.includedFiles), ShouldEqual, 1)
			So(len(config.Hosts), ShouldEqual, 13)
			So(config.Hosts["aaa"].HostName, ShouldEqual, "1.2.3.4")
			So(config.Hosts["aaa"].Port, ShouldEqual, "")
			So(config.Hosts["aaa"].User, ShouldEqual, "")
			So(config.Hosts["bbb"].HostName, ShouldEqual, "$ENV_VAR_HOSTNAME")
			So(config.Hosts["bbb"].Port, ShouldEqual, "${ENV_VAR_PORT}")
			So(config.Hosts["bbb"].User, ShouldEqual, "user-$ENV_VAR_USER-user")
			So(config.Hosts["bbb"].IdentityFile, ShouldEqual, "${NON_EXISTING_ENV_VAR}")
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
			defer os.Remove(file.Name())
			file.Write([]byte(`
includes:
- $DUMMY_ENV_VAR/assh-tests*
`))
			tempDir, err := ioutil.TempDir(os.TempDir(), "assh-tests")
			So(err, ShouldBeNil)
			defer os.RemoveAll(tempDir)

			file2, err := ioutil.TempFile(tempDir, "assh-tests")
			So(err, ShouldBeNil)
			defer os.Remove(file2.Name())
			os.Setenv("DUMMY_ENV_VAR", tempDir)

			config.LoadFiles(file.Name())

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
			So(host.Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})
			So(host.PasswordAuthentication, ShouldEqual, "yes")

			host, err = config.GetHost("tutu")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"tutu":  true,
				"toto":  true,
				"*.ddd": true,
			})
			So(host.User, ShouldEqual, "root")
			So(host.Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})
			So(host.PasswordAuthentication, ShouldEqual, "yes")

			host, err = config.GetHost("nnn")
			So(err, ShouldBeNil)
			So(host.inherited, ShouldResemble, map[string]bool{
				"nnn": true,
				"mmm": true,
			})
			So(host.User, ShouldEqual, "mmmm")
			So(host.Port, ShouldEqual, "26")
			So(host.Gateways, ShouldResemble, []string{"titi", "direct", "1.2.3.4"})
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

func TestConfig_WriteSshConfig(t *testing.T) {
	Convey("Testing Config.WriteSshConfig", t, func() {
		config := dummyConfig()

		var buffer bytes.Buffer

		err := config.WriteSshConfigTo(&buffer)
		So(err, ShouldBeNil)

		expected := `# ssh config generated by advanced-ssh-config

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
  # NoControlMasterMkdir: true
  # Inherits: [mmm]
  # Gateways: [titi, direct, 1.2.3.4]

Host tata
  PasswordAuthentication yes
  Port 22
  User moul
  # ProxyCommand nc -v 4242
  # HostName: 1.2.3.4
  # NoControlMasterMkdir: true
  # Inherits: [tutu, titi, toto, tutu]
  # Gateways: [titi, direct, 1.2.3.4]

Host titi
  Port 23
  User moul
  # ProxyCommand nc -v 4242
  # HostName: tata
  # NoControlMasterMkdir: true

Host tonton
  # ResolveNameservers: [a.com, 1.2.3.4]

Host toto
  # HostName: 1.2.3.4

Host toutou
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
  EnableSSHKeysign yes
  EscapeChar ~
  ExitOnForwardFailure yes
  FingerprintHash sha256
  ForwardAgent yes
  ForwardX11 yes
  ForwardX11Timeout 42
  ForwardX11Trusted yes
  GatewayPorts yes
  GlobalKnownHostsFile /etc/ssh/ssh_known_hosts
  GSSAPIAuthentication no
  GSSAPIClientIdentity moul
  GSSAPIDelegateCredentials no
  GSSAPIKeyExchange no
  GSSAPIRenewalForcesRekey no
  GSSAPIServerIdentity gssapi.example.com
  GSSAPITrustDns no
  HashKnownHosts no
  HostbasedAuthentication no
  HostbasedKeyTypes *
  HostKeyAlgorithms ecdsa-sha2-nistp256-cert-v01@openssh.com
  HostKeyAlias z
  IdentitiesOnly yes
  IdentityFile ~/.ssh/identity
  IgnoreUnknown testtest
  IPQoS lowdelay
  KbdInteractiveAuthentication yes
  KbdInteractiveDevices bsdauth
  KexAlgorithms curve25519-sha256@libssh.org
  KeychainIntegration yes
  LocalCommand echo %h > /tmp/logs
  LocalForward 0.0.0.0:1234
  LogLevel DEBUG3
  MACs umac-64-etm@openssh.com,umac-128-etm@openssh.com
  Match all
  NoHostAuthenticationForLocalhost yes
  NumberOfPasswordPrompts 3
  PasswordAuthentication yes
  PermitLocalCommand yes
  PKCS11Provider /a/b/c/pkcs11.so
  Port 22
  PreferredAuthentications gssapi-with-mic,hostbased,publickey
  Protocol 2
  ProxyUseFdpass no
  PubkeyAuthentication yes
  RekeyLimit default none
  RemoteForward 0.0.0.0:1234
  RequestTTY yes
  RevokedHostKeys /a/revoked-keys
  RhostsRSAAuthentication no
  RSAAuthentication yes
  SendEnv CUSTOM_*,TEST
  ServerAliveCountMax 3
  StreamLocalBindMask 0177
  StreamLocalBindUnlink no
  StrictHostKeyChecking ask
  TCPKeepAlive yes
  Tunnel yes
  TunnelDevice any:any
  UpdateHostKeys ask
  UsePrivilegedPort no
  User moul
  UserKnownHostsFile ~/.ssh/known_hosts ~/.ssh/known_hosts2
  VerifyHostKeyDNS no
  VisualHostKey yes
  XAuthLocation xauth
  # ProxyCommand nc %h %p
  # HostName: zzz.com

# global configuration
Host *
  Port 22
  User root
  ProxyCommand assh proxy --port=%p %h
`
		So(buffer.String(), ShouldEqual, expected)
	})
}

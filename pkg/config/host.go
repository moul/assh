package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os/user"
	"strings"

	"github.com/moul/advanced-ssh-config/pkg/utils"
)

// Host defines the configuration flags of a host
type Host struct {
	// ssh-config fields
	AddressFamily                    string `yaml:"addressfamily,omitempty,flow" json:"AddressFamily,omitempty"`
	AskPassGUI                       string `yaml:"askpassgui,omitempty,flow" json:"AskPassGUI,omitempty"`
	BatchMode                        string `yaml:"batchmode,omitempty,flow" json:"BatchMode,omitempty"`
	BindAddress                      string `yaml:"bindaddress,omitempty,flow" json:"BindAddress,omitempty"`
	CanonicalDomains                 string `yaml:"canonicaldomains,omitempty,flow" json:"CanonicalDomains,omitempty"`
	CanonicalizeFallbackLocal        string `yaml:"canonicalizefallbacklocal,omitempty,flow" json:"CanonicalizeFallbackLocal,omitempty"`
	CanonicalizeHostname             string `yaml:"canonicalizehostname,omitempty,flow" json:"CanonicalizeHostname,omitempty"`
	CanonicalizeMaxDots              string `yaml:"canonicalizemaxDots,omitempty,flow" json:"CanonicalizeMaxDots,omitempty"`
	CanonicalizePermittedCNAMEs      string `yaml:"canonicalizepermittedcnames,omitempty,flow" json:"CanonicalizePermittedCNAMEs,omitempty"`
	ChallengeResponseAuthentication  string `yaml:"challengeresponseauthentication,omitempty,flow" json:"ChallengeResponseAuthentication,omitempty"`
	CheckHostIP                      string `yaml:"checkhostip,omitempty,flow" json:"CheckHostIP,omitempty"`
	Cipher                           string `yaml:"cipher,omitempty,flow" json:"Cipher,omitempty"`
	Ciphers                          string `yaml:"ciphers,omitempty,flow" json:"Ciphers,omitempty"`
	ClearAllForwardings              string `yaml:"clearallforwardings,omitempty,flow" json:"ClearAllForwardings,omitempty"`
	Compression                      string `yaml:"compression,omitempty,flow" json:"Compression,omitempty"`
	CompressionLevel                 int    `yaml:"compressionlevel,omitempty,flow" json:"CompressionLevel,omitempty"`
	ConnectionAttempts               string `yaml:"connectionattempts,omitempty,flow" json:"ConnectionAttempts,omitempty"`
	ConnectTimeout                   int    `yaml:"connecttimeout,omitempty,flow" json:"ConnectTimeout,omitempty"`
	ControlMaster                    string `yaml:"controlmaster,omitempty,flow" json:"ControlMaster,omitempty"`
	ControlPath                      string `yaml:"controlpath,omitempty,flow" json:"ControlPath,omitempty"`
	ControlPersist                   string `yaml:"controlpersist,omitempty,flow" json:"ControlPersist,omitempty"`
	DynamicForward                   string `yaml:"dynamicforward,omitempty,flow" json:"DynamicForward,omitempty"`
	EnableSSHKeysign                 string `yaml:"enablesshkeysign,omitempty,flow" json:"EnableSSHKeysign,omitempty"`
	EscapeChar                       string `yaml:"escapechar,omitempty,flow" json:"EscapeChar,omitempty"`
	ExitOnForwardFailure             string `yaml:"exitonforwardfailure,omitempty,flow" json:"ExitOnForwardFailure,omitempty"`
	FingerprintHash                  string `yaml:"fingerprinthash,omitempty,flow" json:"FingerprintHash,omitempty"`
	ForwardAgent                     string `yaml:"forwardagent,omitempty,flow" json:"ForwardAgent,omitempty"`
	ForwardX11                       string `yaml:"forwardx11,omitempty,flow" json:"ForwardX11,omitempty"`
	ForwardX11Timeout                int    `yaml:"forwardx11timeout,omitempty,flow" json:"ForwardX11Timeout,omitempty"`
	ForwardX11Trusted                string `yaml:"forwardx11trusted,omitempty,flow" json:"ForwardX11Trusted,omitempty"`
	GatewayPorts                     string `yaml:"gatewayports,omitempty,flow" json:"GatewayPorts,omitempty"`
	GlobalKnownHostsFile             string `yaml:"globalknownhostsfile,omitempty,flow" json:"GlobalKnownHostsFile,omitempty"`
	GSSAPIAuthentication             string `yaml:"gssapiauthentication,omitempty,flow" json:"GSSAPIAuthentication,omitempty"`
	GSSAPIClientIdentity             string `yaml:"gssapiclientidentity,omitempty,flow" json:"GSSAPIClientIdentity,omitempty"`
	GSSAPIDelegateCredentials        string `yaml:"gssapidelegatecredentials,omitempty,flow" json:"GSSAPIDelegateCredentials,omitempty"`
	GSSAPIKeyExchange                string `yaml:"gssapikeyexchange,omitempty,flow" json:"GSSAPIKeyExchange,omitempty"`
	GSSAPIRenewalForcesRekey         string `yaml:"gssapirenewalforcesrekey,omitempty,flow" json:"GSSAPIRenewalForcesRekey,omitempty"`
	GSSAPIServerIdentity             string `yaml:"gssapiserveridentity,omitempty,flow" json:"GSSAPIServerIdentity,omitempty"`
	GSSAPITrustDns                   string `yaml:"gssapitrustdns,omitempty,flow" json:"GSSAPITrustDns,omitempty"`
	HashKnownHosts                   string `yaml:"hashknownhosts,omitempty,flow" json:"HashKnownHosts,omitempty"`
	HostbasedAuthentication          string `yaml:"hostbasedauthentication,omitempty,flow" json:"HostbasedAuthentication,omitempty"`
	HostbasedKeyTypes                string `yaml:"hostbasedkeytypes,omitempty,flow" json:"HostbasedKeyTypes,omitempty"`
	HostKeyAlgorithms                string `yaml:"hostkeyalgorithms,omitempty,flow" json:"HostKeyAlgorithms,omitempty"`
	HostKeyAlias                     string `yaml:"hostkeyalias,omitempty,flow" json:"HostKeyAlias,omitempty"`
	IdentitiesOnly                   string `yaml:"identitiesonly,omitempty,flow" json:"IdentitiesOnly,omitempty"`
	IdentityFile                     string `yaml:"identityfile,omitempty,flow" json:"IdentityFile,omitempty"`
	IgnoreUnknown                    string `yaml:"ignoreunknown,omitempty,flow" json:"IgnoreUnknown,omitempty"`
	IPQoS                            string `yaml:"ipqos,omitempty,flow" json:"IPQoS,omitempty"`
	KbdInteractiveAuthentication     string `yaml:"kbdinteractiveauthentication,omitempty,flow" json:"KbdInteractiveAuthentication,omitempty"`
	KbdInteractiveDevices            string `yaml:"kbdinteractivedevices,omitempty,flow" json:"KbdInteractiveDevices,omitempty"`
	KexAlgorithms                    string `yaml:"kexalgorithms,omitempty,flow" json:"KexAlgorithms,omitempty"`
	KeychainIntegration              string `yaml:"keychainintegration,omitempty,flow" json:"KeychainIntegration,omitempty"`
	LocalCommand                     string `yaml:"localcommand,omitempty,flow" json:"LocalCommand,omitempty"`
	LocalForward                     string `yaml:"localforward,omitempty,flow" json:"LocalForward,omitempty"`
	LogLevel                         string `yaml:"loglevel,omitempty,flow" json:"LogLevel,omitempty"`
	MACs                             string `yaml:"macs,omitempty,flow" json:"MACs,omitempty"`
	Match                            string `yaml:"match,omitempty,flow" json:"Match,omitempty"`
	NoHostAuthenticationForLocalhost string `yaml:"nohostauthenticationforlocalhost,omitempty,flow" json:"NoHostAuthenticationForLocalhost,omitempty"`
	NumberOfPasswordPrompts          string `yaml:"numberofpasswordprompts,omitempty,flow" json:"NumberOfPasswordPrompts,omitempty"`
	PasswordAuthentication           string `yaml:"passwordauthentication,omitempty,flow" json:"PasswordAuthentication,omitempty"`
	PermitLocalCommand               string `yaml:"permitlocalcommand,omitempty,flow" json:"PermitLocalCommand,omitempty"`
	PKCS11Provider                   string `yaml:"pkcs11provider,omitempty,flow" json:"PKCS11Provider,omitempty"`
	Port                             string `yaml:"port,omitempty,flow" json:"Port,omitempty"`
	PreferredAuthentications         string `yaml:"preferredauthentications,omitempty,flow" json:"PreferredAuthentications,omitempty"`
	Protocol                         string `yaml:"protocol,omitempty,flow" json:"Protocol,omitempty"`
	ProxyUseFdpass                   string `yaml:"proxyusefdpass,omitempty,flow" json:"ProxyUseFdpass,omitempty"`
	PubkeyAuthentication             string `yaml:"pubkeyauthentication,omitempty,flow" json:"PubkeyAuthentication,omitempty"`
	RekeyLimit                       string `yaml:"rekeylimit,omitempty,flow" json:"RekeyLimit,omitempty"`
	RemoteForward                    string `yaml:"remoteforward,omitempty,flow" json:"RemoteForward,omitempty"`
	RequestTTY                       string `yaml:"requesttty,omitempty,flow" json:"RequestTTY,omitempty"`
	RevokedHostKeys                  string `yaml:"revokedhostkeys,omitempty,flow" json:"RevokedHostKeys,omitempty"`
	RhostsRSAAuthentication          string `yaml:"rhostsrsaauthentication,omitempty,flow" json:"RhostsRSAAuthentication,omitempty"`
	RSAAuthentication                string `yaml:"rsaauthentication,omitempty,flow" json:"RSAAuthentication,omitempty"`
	SendEnv                          string `yaml:"sendenv,omitempty,flow" json:"SendEnv,omitempty"`
	ServerAliveCountMax              int    `yaml:"serveralivecountmax,omitempty,flow" json:"ServerAliveCountMax,omitempty"`
	ServerAliveInterval              int    `yaml:"serveraliveinterval,omitempty,flow" json:"ServerAliveInterval,omitempty"`
	StreamLocalBindMask              string `yaml:"streamlocalbindmask,omitempty,flow" json:"StreamLocalBindMask,omitempty"`
	StreamLocalBindUnlink            string `yaml:"streamlocalbindunlink,omitempty,flow" json:"StreamLocalBindUnlink,omitempty"`
	StrictHostKeyChecking            string `yaml:"stricthostkeychecking,omitempty,flow" json:"StrictHostKeyChecking,omitempty"`
	TCPKeepAlive                     string `yaml:"tcpkeepalive,omitempty,flow" json:"TCPKeepAlive,omitempty"`
	Tunnel                           string `yaml:"tunnel,omitempty,flow" json:"Tunnel,omitempty"`
	TunnelDevice                     string `yaml:"tunneldevice,omitempty,flow" json:"TunnelDevice,omitempty"`
	UpdateHostKeys                   string `yaml:"updatehostkeys,omitempty,flow" json:"UpdateHostKeys,omitempty"`
	UsePrivilegedPort                string `yaml:"useprivilegedport,omitempty,flow" json:"UsePrivilegedPort,omitempty"`
	User                             string `yaml:"user,omitempty,flow" json:"User,omitempty"`
	UserKnownHostsFile               string `yaml:"userknownhostsfile,omitempty,flow" json:"UserKnownHostsFile,omitempty"`
	VerifyHostKeyDNS                 string `yaml:"verifyhostkeydns,omitempty,flow" json:"VerifyHostKeyDNS,omitempty"`
	VisualHostKey                    string `yaml:"visualhostkey,omitempty,flow" json:"VisualHostKey,omitempty"`
	XAuthLocation                    string `yaml:"xauthlocation,omitempty,flow" json:"XAuthLocation,omitempty"`

	// ssh-config fields with a different behavior
	HostName     string `yaml:"hostname,omitempty,flow" json:"HostName,omitempty"`
	ProxyCommand string `yaml:"proxycommand,omitempty,flow" json:"ProxyCommand,omitempty"`

	// exposed assh fields
	Inherits             []string   `yaml:"inherits,omitempty,flow" json:"Inherits,omitempty"`
	Gateways             []string   `yaml:"gateways,omitempty,flow" json:"Gateways,omitempty"`
	ResolveNameservers   []string   `yaml:"resolvenameservers,omitempty,flow" json:"ResolveNameservers,omitempty"`
	ResolveCommand       string     `yaml:"resolvecommand,omitempty,flow" json:"ResolveCommand,omitempty"`
	NoControlMasterMkdir string     `yaml:"nocontrolmastermkdir,omitempty,flow" json:"NoControlMasterMkdir,omitempty"`
	Aliases              []string   `yaml:"aliases,omitempty,flow" json:"Aliases,omitempty"`
	Hooks                *HostHooks `yaml:"hooks,omitempty,flow" json:"Hooks,omitempty"`

	// private assh fields
	knownHosts []string
	pattern    string
	name       string
	inputName  string
	isDefault  bool
	isTemplate bool
	inherited  map[string]bool
}

// NewHost returns a host with name
func NewHost(name string) *Host {
	return &Host{
		name: name,
	}
}

// String returns the JSON output
func (h *Host) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

// Prototype returns a prototype representation of the host, used in listings
func (h *Host) Prototype() string {
	username := h.User
	if username == "" {
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		username = currentUser.Username
	}

	hostname := h.HostName
	if hostname == "" {
		hostname = "[hostname_not_specified]"
	}

	port := h.Port
	if port == "" {
		port = "22"
	}
	return fmt.Sprintf("%s@%s:%s", username, hostname, port)
}

// Name returns the name of a host
func (h *Host) Name() string {
	return h.name
}

// Clone returns a copy of an existing Host
func (h *Host) Clone() *Host {
	newHost := *h
	return &newHost
}

// Matches returns true if the host matches a given string
func (h *Host) Matches(needle string) bool {
	if matches := strings.Contains(h.Name(), needle); matches {
		return true
	}

	for _, opt := range h.Options() {
		if matches := strings.Contains(opt.Value, needle); matches {
			return true
		}
	}
	return false
}

// Options returns a map of set options
func (h *Host) Options() OptionsList {
	options := make(OptionsList, 0)

	// ssh-config fields
	if h.AddressFamily != "" {
		options = append(options, Option{Name: "AddressFamily", Value: h.AddressFamily})
	}
	if h.AskPassGUI != "" {
		options = append(options, Option{Name: "AskPassGUI", Value: h.AskPassGUI})
	}
	if h.BatchMode != "" {
		options = append(options, Option{Name: "BatchMode", Value: h.BatchMode})
	}
	if h.BindAddress != "" {
		options = append(options, Option{Name: "BindAddress", Value: h.BindAddress})
	}
	if h.CanonicalDomains != "" {
		options = append(options, Option{Name: "CanonicalDomains", Value: h.CanonicalDomains})
	}
	if h.CanonicalizeFallbackLocal != "" {
		options = append(options, Option{Name: "CanonicalizeFallbackLocal", Value: h.CanonicalizeFallbackLocal})
	}
	if h.CanonicalizeHostname != "" {
		options = append(options, Option{Name: "CanonicalizeHostname", Value: h.CanonicalizeHostname})
	}
	if h.CanonicalizeMaxDots != "" {
		options = append(options, Option{Name: "CanonicalizeMaxDots", Value: h.CanonicalizeMaxDots})
	}
	if h.CanonicalizePermittedCNAMEs != "" {
		options = append(options, Option{Name: "CanonicalizePermittedCNAMEs", Value: h.CanonicalizePermittedCNAMEs})
	}
	if h.ChallengeResponseAuthentication != "" {
		options = append(options, Option{Name: "ChallengeResponseAuthentication", Value: h.ChallengeResponseAuthentication})
	}
	if h.CheckHostIP != "" {
		options = append(options, Option{Name: "CheckHostIP", Value: h.CheckHostIP})
	}
	if h.Cipher != "" {
		options = append(options, Option{Name: "Cipher", Value: h.Cipher})
	}
	if h.Ciphers != "" {
		options = append(options, Option{Name: "Ciphers", Value: h.Ciphers})
	}
	if h.ClearAllForwardings != "" {
		options = append(options, Option{Name: "ClearAllForwardings", Value: h.ClearAllForwardings})
	}
	if h.Compression != "" {
		options = append(options, Option{Name: "Compression", Value: h.Compression})
	}
	if h.CompressionLevel != 0 {
		options = append(options, Option{Name: "CompressionLevel", Value: string(h.CompressionLevel)})
	}
	if h.ConnectionAttempts != "" {
		options = append(options, Option{Name: "ConnectionAttempts", Value: h.ConnectionAttempts})
	}
	if h.ConnectTimeout != 0 {
		options = append(options, Option{Name: "ConnectTimeout", Value: string(h.ConnectTimeout)})
	}
	if h.ControlMaster != "" {
		options = append(options, Option{Name: "ControlMaster", Value: h.ControlMaster})
	}
	if h.ControlPath != "" {
		options = append(options, Option{Name: "ControlPath", Value: h.ControlPath})
	}
	if h.ControlPersist != "" {
		options = append(options, Option{Name: "ControlPersist", Value: h.ControlPersist})
	}
	if h.DynamicForward != "" {
		options = append(options, Option{Name: "DynamicForward", Value: h.DynamicForward})
	}
	if h.EnableSSHKeysign != "" {
		options = append(options, Option{Name: "EnableSSHKeysign", Value: h.EnableSSHKeysign})
	}
	if h.EscapeChar != "" {
		options = append(options, Option{Name: "EscapeChar", Value: h.EscapeChar})
	}
	if h.ExitOnForwardFailure != "" {
		options = append(options, Option{Name: "ExitOnForwardFailure", Value: h.ExitOnForwardFailure})
	}
	if h.FingerprintHash != "" {
		options = append(options, Option{Name: "FingerprintHash", Value: h.FingerprintHash})
	}
	if h.ForwardAgent != "" {
		options = append(options, Option{Name: "ForwardAgent", Value: h.ForwardAgent})
	}
	if h.ForwardX11 != "" {
		options = append(options, Option{Name: "ForwardX11", Value: h.ForwardX11})
	}
	if h.ForwardX11Timeout != 0 {
		options = append(options, Option{Name: "ForwardX11Timeout", Value: string(h.ForwardX11Timeout)})
	}
	if h.ForwardX11Trusted != "" {
		options = append(options, Option{Name: "ForwardX11Trusted", Value: h.ForwardX11Trusted})
	}
	if h.GatewayPorts != "" {
		options = append(options, Option{Name: "GatewayPorts", Value: h.GatewayPorts})
	}
	if h.GlobalKnownHostsFile != "" {
		options = append(options, Option{Name: "GlobalKnownHostsFile", Value: h.GlobalKnownHostsFile})
	}
	if h.GSSAPIAuthentication != "" {
		options = append(options, Option{Name: "GSSAPIAuthentication", Value: h.GSSAPIAuthentication})
	}
	if h.GSSAPIClientIdentity != "" {
		options = append(options, Option{Name: "GSSAPIClientIdentity", Value: h.GSSAPIClientIdentity})
	}
	if h.GSSAPIDelegateCredentials != "" {
		options = append(options, Option{Name: "GSSAPIDelegateCredentials", Value: h.GSSAPIDelegateCredentials})
	}
	if h.GSSAPIKeyExchange != "" {
		options = append(options, Option{Name: "GSSAPIKeyExchange", Value: h.GSSAPIKeyExchange})
	}
	if h.GSSAPIRenewalForcesRekey != "" {
		options = append(options, Option{Name: "GSSAPIRenewalForcesRekey", Value: h.GSSAPIRenewalForcesRekey})
	}
	if h.GSSAPIServerIdentity != "" {
		options = append(options, Option{Name: "GSSAPIServerIdentity", Value: h.GSSAPIServerIdentity})
	}
	if h.GSSAPITrustDns != "" {
		options = append(options, Option{Name: "GSSAPITrustDns", Value: h.GSSAPITrustDns})
	}
	if h.HashKnownHosts != "" {
		options = append(options, Option{Name: "HashKnownHosts", Value: h.HashKnownHosts})
	}
	if h.HostbasedAuthentication != "" {
		options = append(options, Option{Name: "HostbasedAuthentication", Value: h.HostbasedAuthentication})
	}
	if h.HostbasedKeyTypes != "" {
		options = append(options, Option{Name: "HostbasedKeyTypes", Value: h.HostbasedKeyTypes})
	}
	if h.HostKeyAlgorithms != "" {
		options = append(options, Option{Name: "HostKeyAlgorithms", Value: h.HostKeyAlgorithms})
	}
	if h.HostKeyAlias != "" {
		options = append(options, Option{Name: "HostKeyAlias", Value: h.HostKeyAlias})
	}
	if h.IdentitiesOnly != "" {
		options = append(options, Option{Name: "IdentitiesOnly", Value: h.IdentitiesOnly})
	}
	if h.IdentityFile != "" {
		options = append(options, Option{Name: "IdentityFile", Value: h.IdentityFile})
	}
	if h.IgnoreUnknown != "" {
		options = append(options, Option{Name: "IgnoreUnknown", Value: h.IgnoreUnknown})
	}
	if h.IPQoS != "" {
		options = append(options, Option{Name: "IPQoS", Value: h.IPQoS})
	}
	if h.KbdInteractiveAuthentication != "" {
		options = append(options, Option{Name: "KbdInteractiveAuthentication", Value: h.KbdInteractiveAuthentication})
	}
	if h.KbdInteractiveDevices != "" {
		options = append(options, Option{Name: "KbdInteractiveDevices", Value: h.KbdInteractiveDevices})
	}
	if h.KexAlgorithms != "" {
		options = append(options, Option{Name: "KexAlgorithms", Value: h.KexAlgorithms})
	}
	if h.KeychainIntegration != "" {
		options = append(options, Option{Name: "KeychainIntegration", Value: h.KeychainIntegration})
	}
	if h.LocalCommand != "" {
		options = append(options, Option{Name: "LocalCommand", Value: h.LocalCommand})
	}
	if h.LocalForward != "" {
		options = append(options, Option{Name: "LocalForward", Value: h.LocalForward})
	}
	if h.LogLevel != "" {
		options = append(options, Option{Name: "LogLevel", Value: h.LogLevel})
	}
	if h.MACs != "" {
		options = append(options, Option{Name: "MACs", Value: h.MACs})
	}
	if h.Match != "" {
		options = append(options, Option{Name: "Match", Value: h.Match})
	}
	if h.NoHostAuthenticationForLocalhost != "" {
		options = append(options, Option{Name: "NoHostAuthenticationForLocalhost", Value: h.NoHostAuthenticationForLocalhost})
	}
	if h.NumberOfPasswordPrompts != "" {
		options = append(options, Option{Name: "NumberOfPasswordPrompts", Value: h.NumberOfPasswordPrompts})
	}
	if h.PasswordAuthentication != "" {
		options = append(options, Option{Name: "PasswordAuthentication", Value: h.PasswordAuthentication})
	}
	if h.PermitLocalCommand != "" {
		options = append(options, Option{Name: "PermitLocalCommand", Value: h.PermitLocalCommand})
	}
	if h.PKCS11Provider != "" {
		options = append(options, Option{Name: "PKCS11Provider", Value: h.PKCS11Provider})
	}
	if h.Port != "" {
		options = append(options, Option{Name: "Port", Value: h.Port})
	}
	if h.PreferredAuthentications != "" {
		options = append(options, Option{Name: "PreferredAuthentications", Value: h.PreferredAuthentications})
	}
	if h.Protocol != "" {
		options = append(options, Option{Name: "Protocol", Value: h.Protocol})
	}
	if h.ProxyUseFdpass != "" {
		options = append(options, Option{Name: "ProxyUseFdpass", Value: h.ProxyUseFdpass})
	}
	if h.PubkeyAuthentication != "" {
		options = append(options, Option{Name: "PubkeyAuthentication", Value: h.PubkeyAuthentication})
	}
	if h.RekeyLimit != "" {
		options = append(options, Option{Name: "RekeyLimit", Value: h.RekeyLimit})
	}
	if h.RemoteForward != "" {
		options = append(options, Option{Name: "RemoteForward", Value: h.RemoteForward})
	}
	if h.RequestTTY != "" {
		options = append(options, Option{Name: "RequestTTY", Value: h.RequestTTY})
	}
	if h.RevokedHostKeys != "" {
		options = append(options, Option{Name: "RevokedHostKeys", Value: h.RevokedHostKeys})
	}
	if h.RhostsRSAAuthentication != "" {
		options = append(options, Option{Name: "RhostsRSAAuthentication", Value: h.RhostsRSAAuthentication})
	}
	if h.RSAAuthentication != "" {
		options = append(options, Option{Name: "RSAAuthentication", Value: h.RSAAuthentication})
	}
	if h.SendEnv != "" {
		options = append(options, Option{Name: "SendEnv", Value: h.SendEnv})
	}
	if h.ServerAliveCountMax != 0 {
		options = append(options, Option{Name: "ServerAliveCountMax", Value: string(h.ServerAliveCountMax)})
	}
	if h.ServerAliveInterval != 0 {
		options = append(options, Option{Name: "ServerAliveInterval", Value: string(h.ServerAliveInterval)})
	}
	if h.StreamLocalBindMask != "" {
		options = append(options, Option{Name: "StreamLocalBindMask", Value: h.StreamLocalBindMask})
	}
	if h.StreamLocalBindUnlink != "" {
		options = append(options, Option{Name: "StreamLocalBindUnlink", Value: h.StreamLocalBindUnlink})
	}
	if h.StrictHostKeyChecking != "" {
		options = append(options, Option{Name: "StrictHostKeyChecking", Value: h.StrictHostKeyChecking})
	}
	if h.TCPKeepAlive != "" {
		options = append(options, Option{Name: "TCPKeepAlive", Value: h.TCPKeepAlive})
	}
	if h.Tunnel != "" {
		options = append(options, Option{Name: "Tunnel", Value: h.Tunnel})
	}
	if h.TunnelDevice != "" {
		options = append(options, Option{Name: "TunnelDevice", Value: h.TunnelDevice})
	}
	if h.UpdateHostKeys != "" {
		options = append(options, Option{Name: "UpdateHostKeys", Value: h.UpdateHostKeys})
	}
	if h.UsePrivilegedPort != "" {
		options = append(options, Option{Name: "UsePrivilegedPort", Value: h.UsePrivilegedPort})
	}
	if h.User != "" {
		options = append(options, Option{Name: "User", Value: h.User})
	}
	if h.UserKnownHostsFile != "" {
		options = append(options, Option{Name: "UserKnownHostsFile", Value: h.UserKnownHostsFile})
	}
	if h.VerifyHostKeyDNS != "" {
		options = append(options, Option{Name: "VerifyHostKeyDNS", Value: h.VerifyHostKeyDNS})
	}
	if h.VisualHostKey != "" {
		options = append(options, Option{Name: "VisualHostKey", Value: h.VisualHostKey})
	}
	if h.XAuthLocation != "" {
		options = append(options, Option{Name: "XAuthLocation", Value: h.XAuthLocation})
	}

	// ssh-config fields with a different behavior
	//HostName
	//ProxyCommand

	// exposed assh fields
	//Inherits
	//Gateways
	//ResolveNameservers
	//ResolveCommand
	//NoControlMasterMkdir
	//Aliases
	//Hooks

	// private assh fields
	//knownHosts
	//pattern
	//name
	//inputName
	//isDefault
	//isTemplate
	//inherited

	return options
}

// ApplyDefaults ensures a Host is valid by filling the missing fields with defaults
func (h *Host) ApplyDefaults(defaults *Host) {
	// ssh-config fields
	if h.AddressFamily == "" {
		h.AddressFamily = defaults.AddressFamily
	}
	h.AddressFamily = utils.ExpandField(h.AddressFamily)

	if h.AskPassGUI == "" {
		h.AskPassGUI = defaults.AskPassGUI
	}
	h.AskPassGUI = utils.ExpandField(h.AskPassGUI)

	if h.BatchMode == "" {
		h.BatchMode = defaults.BatchMode
	}
	h.BatchMode = utils.ExpandField(h.BatchMode)

	if h.BindAddress == "" {
		h.BindAddress = defaults.BindAddress
	}
	h.BindAddress = utils.ExpandField(h.BindAddress)

	if h.CanonicalDomains == "" {
		h.CanonicalDomains = defaults.CanonicalDomains
	}
	h.CanonicalDomains = utils.ExpandField(h.CanonicalDomains)

	if h.CanonicalizeFallbackLocal == "" {
		h.CanonicalizeFallbackLocal = defaults.CanonicalizeFallbackLocal
	}
	h.CanonicalizeFallbackLocal = utils.ExpandField(h.CanonicalizeFallbackLocal)

	if h.CanonicalizeHostname == "" {
		h.CanonicalizeHostname = defaults.CanonicalizeHostname
	}
	h.CanonicalizeHostname = utils.ExpandField(h.CanonicalizeHostname)

	if h.CanonicalizeMaxDots == "" {
		h.CanonicalizeMaxDots = defaults.CanonicalizeMaxDots
	}
	h.CanonicalizeMaxDots = utils.ExpandField(h.CanonicalizeMaxDots)

	if h.CanonicalizePermittedCNAMEs == "" {
		h.CanonicalizePermittedCNAMEs = defaults.CanonicalizePermittedCNAMEs
	}
	h.CanonicalizePermittedCNAMEs = utils.ExpandField(h.CanonicalizePermittedCNAMEs)

	if h.ChallengeResponseAuthentication == "" {
		h.ChallengeResponseAuthentication = defaults.ChallengeResponseAuthentication
	}
	h.ChallengeResponseAuthentication = utils.ExpandField(h.ChallengeResponseAuthentication)

	if h.CheckHostIP == "" {
		h.CheckHostIP = defaults.CheckHostIP
	}
	h.CheckHostIP = utils.ExpandField(h.CheckHostIP)

	if h.Cipher == "" {
		h.Cipher = defaults.Cipher
	}
	h.Cipher = utils.ExpandField(h.Cipher)

	if h.Ciphers == "" {
		h.Ciphers = defaults.Ciphers
	}
	h.Ciphers = utils.ExpandField(h.Ciphers)

	if h.ClearAllForwardings == "" {
		h.ClearAllForwardings = defaults.ClearAllForwardings
	}
	h.ClearAllForwardings = utils.ExpandField(h.ClearAllForwardings)

	if h.Compression == "" {
		h.Compression = defaults.Compression
	}
	h.Compression = utils.ExpandField(h.Compression)

	if h.CompressionLevel == 0 {
		h.CompressionLevel = defaults.CompressionLevel
	}
	// h.CompressionLevel = utils.ExpandField(h.CompressionLevel)

	if h.ConnectionAttempts == "" {
		h.ConnectionAttempts = defaults.ConnectionAttempts
	}
	h.ConnectionAttempts = utils.ExpandField(h.ConnectionAttempts)

	if h.ConnectTimeout == 0 {
		h.ConnectTimeout = defaults.ConnectTimeout
	}
	// h.ConnectTimeout = utils.ExpandField(h.ConnectTimeout)

	if h.ControlMaster == "" {
		h.ControlMaster = defaults.ControlMaster
	}
	h.ControlMaster = utils.ExpandField(h.ControlMaster)

	if h.ControlPath == "" {
		h.ControlPath = defaults.ControlPath
	}
	h.ControlPath = utils.ExpandField(h.ControlPath)

	if h.ControlPersist == "" {
		h.ControlPersist = defaults.ControlPersist
	}
	h.ControlPersist = utils.ExpandField(h.ControlPersist)

	if h.DynamicForward == "" {
		h.DynamicForward = defaults.DynamicForward
	}
	h.DynamicForward = utils.ExpandField(h.DynamicForward)

	if h.EnableSSHKeysign == "" {
		h.EnableSSHKeysign = defaults.EnableSSHKeysign
	}
	h.EnableSSHKeysign = utils.ExpandField(h.EnableSSHKeysign)

	if h.EscapeChar == "" {
		h.EscapeChar = defaults.EscapeChar
	}
	h.EscapeChar = utils.ExpandField(h.EscapeChar)

	if h.ExitOnForwardFailure == "" {
		h.ExitOnForwardFailure = defaults.ExitOnForwardFailure
	}
	h.ExitOnForwardFailure = utils.ExpandField(h.ExitOnForwardFailure)

	if h.FingerprintHash == "" {
		h.FingerprintHash = defaults.FingerprintHash
	}
	h.FingerprintHash = utils.ExpandField(h.FingerprintHash)

	if h.ForwardAgent == "" {
		h.ForwardAgent = defaults.ForwardAgent
	}
	h.ForwardAgent = utils.ExpandField(h.ForwardAgent)

	if h.ForwardX11 == "" {
		h.ForwardX11 = defaults.ForwardX11
	}
	h.ForwardX11 = utils.ExpandField(h.ForwardX11)

	if h.ForwardX11Timeout == 0 {
		h.ForwardX11Timeout = defaults.ForwardX11Timeout
	}
	// h.ForwardX11Timeout = utils.ExpandField(h.ForwardX11Timeout)

	if h.ForwardX11Trusted == "" {
		h.ForwardX11Trusted = defaults.ForwardX11Trusted
	}
	h.ForwardX11Trusted = utils.ExpandField(h.ForwardX11Trusted)

	if h.GatewayPorts == "" {
		h.GatewayPorts = defaults.GatewayPorts
	}
	h.GatewayPorts = utils.ExpandField(h.GatewayPorts)

	if h.GlobalKnownHostsFile == "" {
		h.GlobalKnownHostsFile = defaults.GlobalKnownHostsFile
	}
	h.GlobalKnownHostsFile = utils.ExpandField(h.GlobalKnownHostsFile)

	if h.GSSAPIAuthentication == "" {
		h.GSSAPIAuthentication = defaults.GSSAPIAuthentication
	}
	h.GSSAPIAuthentication = utils.ExpandField(h.GSSAPIAuthentication)

	if h.GSSAPIClientIdentity == "" {
		h.GSSAPIClientIdentity = defaults.GSSAPIClientIdentity
	}
	h.GSSAPIClientIdentity = utils.ExpandField(h.GSSAPIClientIdentity)

	if h.GSSAPIDelegateCredentials == "" {
		h.GSSAPIDelegateCredentials = defaults.GSSAPIDelegateCredentials
	}
	h.GSSAPIDelegateCredentials = utils.ExpandField(h.GSSAPIDelegateCredentials)

	if h.GSSAPIKeyExchange == "" {
		h.GSSAPIKeyExchange = defaults.GSSAPIKeyExchange
	}
	h.GSSAPIKeyExchange = utils.ExpandField(h.GSSAPIKeyExchange)

	if h.GSSAPIRenewalForcesRekey == "" {
		h.GSSAPIRenewalForcesRekey = defaults.GSSAPIRenewalForcesRekey
	}
	h.GSSAPIRenewalForcesRekey = utils.ExpandField(h.GSSAPIRenewalForcesRekey)

	if h.GSSAPIServerIdentity == "" {
		h.GSSAPIServerIdentity = defaults.GSSAPIServerIdentity
	}
	h.GSSAPIServerIdentity = utils.ExpandField(h.GSSAPIServerIdentity)

	if h.GSSAPITrustDns == "" {
		h.GSSAPITrustDns = defaults.GSSAPITrustDns
	}
	h.GSSAPITrustDns = utils.ExpandField(h.GSSAPITrustDns)

	if h.HashKnownHosts == "" {
		h.HashKnownHosts = defaults.HashKnownHosts
	}
	h.HashKnownHosts = utils.ExpandField(h.HashKnownHosts)

	if h.HostbasedAuthentication == "" {
		h.HostbasedAuthentication = defaults.HostbasedAuthentication
	}
	h.HostbasedAuthentication = utils.ExpandField(h.HostbasedAuthentication)

	if h.HostbasedKeyTypes == "" {
		h.HostbasedKeyTypes = defaults.HostbasedKeyTypes
	}
	h.HostbasedKeyTypes = utils.ExpandField(h.HostbasedKeyTypes)

	if h.HostKeyAlgorithms == "" {
		h.HostKeyAlgorithms = defaults.HostKeyAlgorithms
	}
	h.HostKeyAlgorithms = utils.ExpandField(h.HostKeyAlgorithms)

	if h.HostKeyAlias == "" {
		h.HostKeyAlias = defaults.HostKeyAlias
	}
	h.HostKeyAlias = utils.ExpandField(h.HostKeyAlias)

	if h.HostName == "" {
		h.HostName = defaults.HostName
	}
	h.HostName = utils.ExpandField(h.HostName)

	if h.IdentitiesOnly == "" {
		h.IdentitiesOnly = defaults.IdentitiesOnly
	}
	h.IdentitiesOnly = utils.ExpandField(h.IdentitiesOnly)

	if h.IdentityFile == "" {
		h.IdentityFile = defaults.IdentityFile
	}
	h.IdentityFile = utils.ExpandField(h.IdentityFile)

	if h.IgnoreUnknown == "" {
		h.IgnoreUnknown = defaults.IgnoreUnknown
	}
	h.IgnoreUnknown = utils.ExpandField(h.IgnoreUnknown)

	if h.IPQoS == "" {
		h.IPQoS = defaults.IPQoS
	}
	h.IPQoS = utils.ExpandField(h.IPQoS)

	if h.KbdInteractiveAuthentication == "" {
		h.KbdInteractiveAuthentication = defaults.KbdInteractiveAuthentication
	}
	h.KbdInteractiveAuthentication = utils.ExpandField(h.KbdInteractiveAuthentication)

	if h.KbdInteractiveDevices == "" {
		h.KbdInteractiveDevices = defaults.KbdInteractiveDevices
	}
	h.KbdInteractiveDevices = utils.ExpandField(h.KbdInteractiveDevices)

	if h.KexAlgorithms == "" {
		h.KexAlgorithms = defaults.KexAlgorithms
	}
	h.KexAlgorithms = utils.ExpandField(h.KexAlgorithms)

	if h.KeychainIntegration == "" {
		h.KeychainIntegration = defaults.KeychainIntegration
	}
	h.KeychainIntegration = utils.ExpandField(h.KeychainIntegration)

	if h.LocalCommand == "" {
		h.LocalCommand = defaults.LocalCommand
	}
	h.LocalCommand = utils.ExpandField(h.LocalCommand)

	if h.LocalForward == "" {
		h.LocalForward = defaults.LocalForward
	}
	h.LocalForward = utils.ExpandField(h.LocalForward)

	if h.LogLevel == "" {
		h.LogLevel = defaults.LogLevel
	}
	h.LogLevel = utils.ExpandField(h.LogLevel)

	if h.MACs == "" {
		h.MACs = defaults.MACs
	}
	h.MACs = utils.ExpandField(h.MACs)

	if h.Match == "" {
		h.Match = defaults.Match
	}
	h.Match = utils.ExpandField(h.Match)

	if h.NoHostAuthenticationForLocalhost == "" {
		h.NoHostAuthenticationForLocalhost = defaults.NoHostAuthenticationForLocalhost
	}
	h.NoHostAuthenticationForLocalhost = utils.ExpandField(h.NoHostAuthenticationForLocalhost)

	if h.NumberOfPasswordPrompts == "" {
		h.NumberOfPasswordPrompts = defaults.NumberOfPasswordPrompts
	}
	h.NumberOfPasswordPrompts = utils.ExpandField(h.NumberOfPasswordPrompts)

	if h.PasswordAuthentication == "" {
		h.PasswordAuthentication = defaults.PasswordAuthentication
	}
	h.PasswordAuthentication = utils.ExpandField(h.PasswordAuthentication)

	if h.PermitLocalCommand == "" {
		h.PermitLocalCommand = defaults.PermitLocalCommand
	}
	h.PermitLocalCommand = utils.ExpandField(h.PermitLocalCommand)

	if h.PKCS11Provider == "" {
		h.PKCS11Provider = defaults.PKCS11Provider
	}
	h.PKCS11Provider = utils.ExpandField(h.PKCS11Provider)

	if h.Port == "" {
		h.Port = defaults.Port
	}
	h.Port = utils.ExpandField(h.Port)

	if h.PreferredAuthentications == "" {
		h.PreferredAuthentications = defaults.PreferredAuthentications
	}
	h.PreferredAuthentications = utils.ExpandField(h.PreferredAuthentications)

	if h.Protocol == "" {
		h.Protocol = defaults.Protocol
	}
	h.Protocol = utils.ExpandField(h.Protocol)

	if h.ProxyCommand == "" {
		h.ProxyCommand = defaults.ProxyCommand
	}
	h.ProxyCommand = utils.ExpandField(h.ProxyCommand)

	if h.ProxyUseFdpass == "" {
		h.ProxyUseFdpass = defaults.ProxyUseFdpass
	}
	h.ProxyUseFdpass = utils.ExpandField(h.ProxyUseFdpass)

	if h.PubkeyAuthentication == "" {
		h.PubkeyAuthentication = defaults.PubkeyAuthentication
	}
	h.PubkeyAuthentication = utils.ExpandField(h.PubkeyAuthentication)

	if h.RekeyLimit == "" {
		h.RekeyLimit = defaults.RekeyLimit
	}
	h.RekeyLimit = utils.ExpandField(h.RekeyLimit)

	if h.RemoteForward == "" {
		h.RemoteForward = defaults.RemoteForward
	}
	h.RemoteForward = utils.ExpandField(h.RemoteForward)

	if h.RequestTTY == "" {
		h.RequestTTY = defaults.RequestTTY
	}
	h.RequestTTY = utils.ExpandField(h.RequestTTY)

	if h.RevokedHostKeys == "" {
		h.RevokedHostKeys = defaults.RevokedHostKeys
	}
	h.RevokedHostKeys = utils.ExpandField(h.RevokedHostKeys)

	if h.RhostsRSAAuthentication == "" {
		h.RhostsRSAAuthentication = defaults.RhostsRSAAuthentication
	}
	h.RhostsRSAAuthentication = utils.ExpandField(h.RhostsRSAAuthentication)

	if h.RSAAuthentication == "" {
		h.RSAAuthentication = defaults.RSAAuthentication
	}
	h.RSAAuthentication = utils.ExpandField(h.RSAAuthentication)

	if h.SendEnv == "" {
		h.SendEnv = defaults.SendEnv
	}
	h.SendEnv = utils.ExpandField(h.SendEnv)

	if h.ServerAliveCountMax == 0 {
		h.ServerAliveCountMax = defaults.ServerAliveCountMax
	}
	// h.ServerAliveCountMax = utils.ExpandField(h.ServerAliveCountMax)

	if h.ServerAliveInterval == 0 {
		h.ServerAliveInterval = defaults.ServerAliveInterval
	}
	// h.ServerAliveInterval = utils.ExpandField(h.ServerAliveInterval)

	if h.StreamLocalBindMask == "" {
		h.StreamLocalBindMask = defaults.StreamLocalBindMask
	}
	h.StreamLocalBindMask = utils.ExpandField(h.StreamLocalBindMask)

	if h.StreamLocalBindUnlink == "" {
		h.StreamLocalBindUnlink = defaults.StreamLocalBindUnlink
	}
	h.StreamLocalBindUnlink = utils.ExpandField(h.StreamLocalBindUnlink)

	if h.StrictHostKeyChecking == "" {
		h.StrictHostKeyChecking = defaults.StrictHostKeyChecking
	}
	h.StrictHostKeyChecking = utils.ExpandField(h.StrictHostKeyChecking)

	if h.TCPKeepAlive == "" {
		h.TCPKeepAlive = defaults.TCPKeepAlive
	}
	h.TCPKeepAlive = utils.ExpandField(h.TCPKeepAlive)

	if h.Tunnel == "" {
		h.Tunnel = defaults.Tunnel
	}
	h.Tunnel = utils.ExpandField(h.Tunnel)

	if h.TunnelDevice == "" {
		h.TunnelDevice = defaults.TunnelDevice
	}
	h.TunnelDevice = utils.ExpandField(h.TunnelDevice)

	if h.UpdateHostKeys == "" {
		h.UpdateHostKeys = defaults.UpdateHostKeys
	}
	h.UpdateHostKeys = utils.ExpandField(h.UpdateHostKeys)

	if h.UsePrivilegedPort == "" {
		h.UsePrivilegedPort = defaults.UsePrivilegedPort
	}
	h.UsePrivilegedPort = utils.ExpandField(h.UsePrivilegedPort)

	if h.User == "" {
		h.User = defaults.User
	}
	h.User = utils.ExpandField(h.User)

	if h.UserKnownHostsFile == "" {
		h.UserKnownHostsFile = defaults.UserKnownHostsFile
	}
	h.UserKnownHostsFile = utils.ExpandField(h.UserKnownHostsFile)

	if h.VerifyHostKeyDNS == "" {
		h.VerifyHostKeyDNS = defaults.VerifyHostKeyDNS
	}
	h.VerifyHostKeyDNS = utils.ExpandField(h.VerifyHostKeyDNS)

	if h.VisualHostKey == "" {
		h.VisualHostKey = defaults.VisualHostKey
	}
	h.VisualHostKey = utils.ExpandField(h.VisualHostKey)

	if h.XAuthLocation == "" {
		h.XAuthLocation = defaults.XAuthLocation
	}
	h.XAuthLocation = utils.ExpandField(h.XAuthLocation)

	// ssh-config fields with a different behavior
	if h.ProxyCommand == "" {
		h.ProxyCommand = defaults.ProxyCommand
	}
	h.ProxyCommand = utils.ExpandField(h.ProxyCommand)

	// exposed assh fields
	if len(h.ResolveNameservers) == 0 {
		h.ResolveNameservers = defaults.ResolveNameservers
	}
	// h.ResolveNameservers = utils.ExpandField(h.ResolveNameservers)

	if h.ResolveCommand == "" {
		h.ResolveCommand = defaults.ResolveCommand
	}
	h.ResolveCommand = utils.ExpandField(h.ResolveCommand)

	if h.NoControlMasterMkdir == "" {
		h.NoControlMasterMkdir = defaults.NoControlMasterMkdir
	}
	h.NoControlMasterMkdir = utils.ExpandField(h.NoControlMasterMkdir)

	if len(h.Gateways) == 0 {
		h.Gateways = defaults.Gateways
	}
	// h.Gateways = utils.ExpandField(h.Gateways)

	if len(h.Aliases) == 0 {
		h.Aliases = defaults.Aliases
	}

	if h.Hooks.Length() == 0 {
		h.Hooks = defaults.Hooks
	}

	if len(h.Inherits) == 0 {
		h.Inherits = defaults.Inherits
	}
	// h.Inherits = utils.ExpandField(h.Inherits)

	// private assh fields
	// h.inherited = make(map[string]bool, 0)
	if h.inputName == "" {
		h.inputName = h.name
	}
	h.inputName = utils.ExpandField(h.inputName)

	// Extra defaults
	if h.Port == "" {
		h.Port = "22"
	}
}

// AddKnownHost append target to the host' known hosts list
func (h *Host) AddKnownHost(target string) {
	h.knownHosts = append(h.knownHosts, target)
}

// WriteSSHConfigTo writes an ~/.ssh/config file compatible host definition to a writable stream
func (h *Host) WriteSSHConfigTo(w io.Writer) error {
	aliases := append([]string{h.Name()}, h.Aliases...)
	aliases = append(aliases, h.knownHosts...)
	aliasIdx := 0
	for _, alias := range aliases {
		// FIXME: skip complex patterns

		if aliasIdx > 0 {
			fmt.Fprint(w, "\n")
		}

		fmt.Fprintf(w, "Host %s\n", alias)

		// ssh-config fields
		if h.AddressFamily != "" {
			fmt.Fprintf(w, "  AddressFamily %s\n", h.AddressFamily)
		}
		if h.AskPassGUI != "" {
			fmt.Fprintf(w, "  AskPassGUI %s\n", h.AskPassGUI)
		}
		if h.BatchMode != "" {
			fmt.Fprintf(w, "  BatchMode %s\n", h.BatchMode)
		}
		if h.BindAddress != "" {
			fmt.Fprintf(w, "  BindAddress %s\n", h.BindAddress)
		}
		if h.CanonicalDomains != "" {
			fmt.Fprintf(w, "  CanonicalDomains %s\n", h.CanonicalDomains)
		}
		if h.CanonicalizeFallbackLocal != "" {
			fmt.Fprintf(w, "  CanonicalizeFallbackLocal %s\n", h.CanonicalizeFallbackLocal)
		}
		if h.CanonicalizeHostname != "" {
			fmt.Fprintf(w, "  CanonicalizeHostname %s\n", h.CanonicalizeHostname)
		}
		if h.CanonicalizeMaxDots != "" {
			fmt.Fprintf(w, "  CanonicalizeMaxDots %s\n", h.CanonicalizeMaxDots)
		}
		if h.CanonicalizePermittedCNAMEs != "" {
			fmt.Fprintf(w, "  CanonicalizePermittedCNAMEs %s\n", h.CanonicalizePermittedCNAMEs)
		}
		if h.ChallengeResponseAuthentication != "" {
			fmt.Fprintf(w, "  ChallengeResponseAuthentication %s\n", h.ChallengeResponseAuthentication)
		}
		if h.CheckHostIP != "" {
			fmt.Fprintf(w, "  CheckHostIP %s\n", h.CheckHostIP)
		}
		if h.Cipher != "" {
			fmt.Fprintf(w, "  Cipher %s\n", h.Cipher)
		}
		if h.Ciphers != "" {
			fmt.Fprintf(w, "  Ciphers %s\n", h.Ciphers)
		}
		if h.ClearAllForwardings != "" {
			fmt.Fprintf(w, "  ClearAllForwardings %s\n", h.ClearAllForwardings)
		}
		if h.Compression != "" {
			fmt.Fprintf(w, "  Compression %s\n", h.Compression)
		}
		if h.CompressionLevel != 0 {
			fmt.Fprintf(w, "  CompressionLevel %d\n", h.CompressionLevel)
		}
		if h.ConnectionAttempts != "" {
			fmt.Fprintf(w, "  ConnectionAttempts %s\n", h.ConnectionAttempts)
		}
		if h.ConnectTimeout != 0 {
			fmt.Fprintf(w, "  ConnectTimeout %d\n", h.ConnectTimeout)
		}
		if h.ControlMaster != "" {
			fmt.Fprintf(w, "  ControlMaster %s\n", h.ControlMaster)
		}
		if h.ControlPath != "" {
			fmt.Fprintf(w, "  ControlPath %s\n", h.ControlPath)
		}
		if h.ControlPersist != "" {
			fmt.Fprintf(w, "  ControlPersist %s\n", h.ControlPersist)
		}
		if h.DynamicForward != "" {
			fmt.Fprintf(w, "  DynamicForward %s\n", h.DynamicForward)
		}
		if h.EnableSSHKeysign != "" {
			fmt.Fprintf(w, "  EnableSSHKeysign %s\n", h.EnableSSHKeysign)
		}
		if h.EscapeChar != "" {
			fmt.Fprintf(w, "  EscapeChar %s\n", h.EscapeChar)
		}
		if h.ExitOnForwardFailure != "" {
			fmt.Fprintf(w, "  ExitOnForwardFailure %s\n", h.ExitOnForwardFailure)
		}
		if h.FingerprintHash != "" {
			fmt.Fprintf(w, "  FingerprintHash %s\n", h.FingerprintHash)
		}
		if h.ForwardAgent != "" {
			fmt.Fprintf(w, "  ForwardAgent %s\n", h.ForwardAgent)
		}
		if h.ForwardX11 != "" {
			fmt.Fprintf(w, "  ForwardX11 %s\n", h.ForwardX11)
		}
		if h.ForwardX11Timeout != 0 {
			fmt.Fprintf(w, "  ForwardX11Timeout %d\n", h.ForwardX11Timeout)
		}
		if h.ForwardX11Trusted != "" {
			fmt.Fprintf(w, "  ForwardX11Trusted %s\n", h.ForwardX11Trusted)
		}
		if h.GatewayPorts != "" {
			fmt.Fprintf(w, "  GatewayPorts %s\n", h.GatewayPorts)
		}
		if h.GlobalKnownHostsFile != "" {
			fmt.Fprintf(w, "  GlobalKnownHostsFile %s\n", h.GlobalKnownHostsFile)
		}
		if h.GSSAPIAuthentication != "" {
			fmt.Fprintf(w, "  GSSAPIAuthentication %s\n", h.GSSAPIAuthentication)
		}
		if h.GSSAPIClientIdentity != "" {
			fmt.Fprintf(w, "  GSSAPIClientIdentity %s\n", h.GSSAPIClientIdentity)
		}
		if h.GSSAPIDelegateCredentials != "" {
			fmt.Fprintf(w, "  GSSAPIDelegateCredentials %s\n", h.GSSAPIDelegateCredentials)
		}
		if h.GSSAPIKeyExchange != "" {
			fmt.Fprintf(w, "  GSSAPIKeyExchange %s\n", h.GSSAPIKeyExchange)
		}
		if h.GSSAPIRenewalForcesRekey != "" {
			fmt.Fprintf(w, "  GSSAPIRenewalForcesRekey %s\n", h.GSSAPIRenewalForcesRekey)
		}
		if h.GSSAPIServerIdentity != "" {
			fmt.Fprintf(w, "  GSSAPIServerIdentity %s\n", h.GSSAPIServerIdentity)
		}
		if h.GSSAPITrustDns != "" {
			fmt.Fprintf(w, "  GSSAPITrustDns %s\n", h.GSSAPITrustDns)
		}
		if h.HashKnownHosts != "" {
			fmt.Fprintf(w, "  HashKnownHosts %s\n", h.HashKnownHosts)
		}
		if h.HostbasedAuthentication != "" {
			fmt.Fprintf(w, "  HostbasedAuthentication %s\n", h.HostbasedAuthentication)
		}
		if h.HostbasedKeyTypes != "" {
			fmt.Fprintf(w, "  HostbasedKeyTypes %s\n", h.HostbasedKeyTypes)
		}
		if h.HostKeyAlgorithms != "" {
			fmt.Fprintf(w, "  HostKeyAlgorithms %s\n", h.HostKeyAlgorithms)
		}
		if h.HostKeyAlias != "" {
			fmt.Fprintf(w, "  HostKeyAlias %s\n", h.HostKeyAlias)
		}
		if h.IdentitiesOnly != "" {
			fmt.Fprintf(w, "  IdentitiesOnly %s\n", h.IdentitiesOnly)
		}
		if h.IdentityFile != "" {
			fmt.Fprintf(w, "  IdentityFile %s\n", h.IdentityFile)
		}
		if h.IgnoreUnknown != "" {
			fmt.Fprintf(w, "  IgnoreUnknown %s\n", h.IgnoreUnknown)
		}
		if h.IPQoS != "" {
			fmt.Fprintf(w, "  IPQoS %s\n", h.IPQoS)
		}
		if h.KbdInteractiveAuthentication != "" {
			fmt.Fprintf(w, "  KbdInteractiveAuthentication %s\n", h.KbdInteractiveAuthentication)
		}
		if h.KbdInteractiveDevices != "" {
			fmt.Fprintf(w, "  KbdInteractiveDevices %s\n", h.KbdInteractiveDevices)
		}
		if h.KexAlgorithms != "" {
			fmt.Fprintf(w, "  KexAlgorithms %s\n", h.KexAlgorithms)
		}
		if h.KeychainIntegration != "" {
			fmt.Fprintf(w, "  KeychainIntegration %s\n", h.KeychainIntegration)
		}
		if h.LocalCommand != "" {
			fmt.Fprintf(w, "  LocalCommand %s\n", h.LocalCommand)
		}
		if h.LocalForward != "" {
			fmt.Fprintf(w, "  LocalForward %s\n", h.LocalForward)
		}
		if h.LogLevel != "" {
			fmt.Fprintf(w, "  LogLevel %s\n", h.LogLevel)
		}
		if h.MACs != "" {
			fmt.Fprintf(w, "  MACs %s\n", h.MACs)
		}
		if h.Match != "" {
			fmt.Fprintf(w, "  Match %s\n", h.Match)
		}
		if h.NoHostAuthenticationForLocalhost != "" {
			fmt.Fprintf(w, "  NoHostAuthenticationForLocalhost %s\n", h.NoHostAuthenticationForLocalhost)
		}
		if h.NumberOfPasswordPrompts != "" {
			fmt.Fprintf(w, "  NumberOfPasswordPrompts %s\n", h.NumberOfPasswordPrompts)
		}
		if h.PasswordAuthentication != "" {
			fmt.Fprintf(w, "  PasswordAuthentication %s\n", h.PasswordAuthentication)
		}
		if h.PermitLocalCommand != "" {
			fmt.Fprintf(w, "  PermitLocalCommand %s\n", h.PermitLocalCommand)
		}
		if h.PKCS11Provider != "" {
			fmt.Fprintf(w, "  PKCS11Provider %s\n", h.PKCS11Provider)
		}
		if h.Port != "" {
			fmt.Fprintf(w, "  Port %s\n", h.Port)
		}
		if h.PreferredAuthentications != "" {
			fmt.Fprintf(w, "  PreferredAuthentications %s\n", h.PreferredAuthentications)
		}
		if h.Protocol != "" {
			fmt.Fprintf(w, "  Protocol %s\n", h.Protocol)
		}
		if h.ProxyUseFdpass != "" {
			fmt.Fprintf(w, "  ProxyUseFdpass %s\n", h.ProxyUseFdpass)
		}
		if h.PubkeyAuthentication != "" {
			fmt.Fprintf(w, "  PubkeyAuthentication %s\n", h.PubkeyAuthentication)
		}
		if h.RekeyLimit != "" {
			fmt.Fprintf(w, "  RekeyLimit %s\n", h.RekeyLimit)
		}
		if h.RemoteForward != "" {
			fmt.Fprintf(w, "  RemoteForward %s\n", h.RemoteForward)
		}
		if h.RequestTTY != "" {
			fmt.Fprintf(w, "  RequestTTY %s\n", h.RequestTTY)
		}
		if h.RevokedHostKeys != "" {
			fmt.Fprintf(w, "  RevokedHostKeys %s\n", h.RevokedHostKeys)
		}
		if h.RhostsRSAAuthentication != "" {
			fmt.Fprintf(w, "  RhostsRSAAuthentication %s\n", h.RhostsRSAAuthentication)
		}
		if h.RSAAuthentication != "" {
			fmt.Fprintf(w, "  RSAAuthentication %s\n", h.RSAAuthentication)
		}
		if h.SendEnv != "" {
			fmt.Fprintf(w, "  SendEnv %s\n", h.SendEnv)
		}
		if h.ServerAliveCountMax != 0 {
			fmt.Fprintf(w, "  ServerAliveCountMax %d\n", h.ServerAliveCountMax)
		}
		if h.ServerAliveInterval != 0 {
			fmt.Fprintf(w, "  ServerAliveInterval %d\n", h.ServerAliveInterval)
		}
		if h.StreamLocalBindMask != "" {
			fmt.Fprintf(w, "  StreamLocalBindMask %s\n", h.StreamLocalBindMask)
		}
		if h.StreamLocalBindUnlink != "" {
			fmt.Fprintf(w, "  StreamLocalBindUnlink %s\n", h.StreamLocalBindUnlink)
		}
		if h.StrictHostKeyChecking != "" {
			fmt.Fprintf(w, "  StrictHostKeyChecking %s\n", h.StrictHostKeyChecking)
		}
		if h.TCPKeepAlive != "" {
			fmt.Fprintf(w, "  TCPKeepAlive %s\n", h.TCPKeepAlive)
		}
		if h.Tunnel != "" {
			fmt.Fprintf(w, "  Tunnel %s\n", h.Tunnel)
		}
		if h.TunnelDevice != "" {
			fmt.Fprintf(w, "  TunnelDevice %s\n", h.TunnelDevice)
		}
		if h.UpdateHostKeys != "" {
			fmt.Fprintf(w, "  UpdateHostKeys %s\n", h.UpdateHostKeys)
		}
		if h.UsePrivilegedPort != "" {
			fmt.Fprintf(w, "  UsePrivilegedPort %s\n", h.UsePrivilegedPort)
		}
		if h.User != "" {
			fmt.Fprintf(w, "  User %s\n", h.User)
		}
		if h.UserKnownHostsFile != "" {
			fmt.Fprintf(w, "  UserKnownHostsFile %s\n", h.UserKnownHostsFile)
		}
		if h.VerifyHostKeyDNS != "" {
			fmt.Fprintf(w, "  VerifyHostKeyDNS %s\n", h.VerifyHostKeyDNS)
		}
		if h.VisualHostKey != "" {
			fmt.Fprintf(w, "  VisualHostKey %s\n", h.VisualHostKey)
		}
		if h.XAuthLocation != "" {
			fmt.Fprintf(w, "  XAuthLocation %s\n", h.XAuthLocation)
		}

		// ssh-config fields with a different behavior
		if h.isDefault {
			fmt.Fprintf(w, "  ProxyCommand %s connect --port=%%p %%h\n", asshBinaryPath)
		} else {
			if h.ProxyCommand != "" {
				fmt.Fprintf(w, "  # ProxyCommand %s\n", h.ProxyCommand)
			}
		}

		// assh fields
		if h.HostName != "" {
			fmt.Fprintf(w, "  # HostName: %s\n", h.HostName)
		}
		if BoolVal(h.NoControlMasterMkdir) {
			fmt.Fprintf(w, "  # NoControlMasterMkdir: true\n")
		}
		if len(h.Inherits) > 0 {
			fmt.Fprintf(w, "  # Inherits: [%s]\n", strings.Join(h.Inherits, ", "))
		}
		if len(h.Gateways) > 0 {
			fmt.Fprintf(w, "  # Gateways: [%s]\n", strings.Join(h.Gateways, ", "))
		}
		if len(h.Aliases) > 0 {
			if aliasIdx == 0 {
				fmt.Fprintf(w, "  # Aliases: [%s]\n", strings.Join(h.Aliases, ", "))
			} else {
				fmt.Fprintf(w, "  # AliasOf: %s\n", h.Name())
			}
		}
		if h.Hooks.Length() > 0 {
			fmt.Fprintf(w, "  # Hooks: [%s]\n", h.Hooks.String())
		}
		if len(h.knownHosts) > 0 {
			if aliasIdx == 0 {
				fmt.Fprintf(w, "  # KnownHosts: [%s]\n", strings.Join(h.knownHosts, ", "))
			} else {
				fmt.Fprintf(w, "  # KnownHostOf: %s\n", h.Name())
			}
		}

		if len(h.ResolveNameservers) > 0 {
			fmt.Fprintf(w, "  # ResolveNameservers: [%s]\n", strings.Join(h.ResolveNameservers, ", "))
		}
		if h.ResolveCommand != "" {
			fmt.Fprintf(w, "  # ResolveCommand: %s\n", h.ResolveCommand)
		}

		aliasIdx++
	}
	return nil
}

func (h *Host) ExpandString(input string) string {
	output := input

	// name of the host in config
	output = strings.Replace(output, "%name", h.Name(), -1)

	// original target host name specified on the command line
	output = strings.Replace(output, "%n", h.inputName, -1)

	// target host name
	output = strings.Replace(output, "%h", h.HostName, -1)

	// port
	output = strings.Replace(output, "%p", fmt.Sprintf("%s", h.Port), -1)

	// FIXME: add
	//   %L -> first component of the local host name
	//   %l -> local host name
	//   %r -> remote login username
	//   %u -> username of the user running assh
	//   %r -> remote login username

	return output
}

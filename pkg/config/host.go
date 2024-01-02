package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os/user"
	"reflect"
	"strings"

	composeyaml "github.com/docker/libcompose/yaml"
	"moul.io/assh/v2/pkg/utils"
)

// EnableEscapeCommandline          string                    `yaml:"enableescapecommandline,omitempty,flow" json:"EnableSSHKeysign,omitempty"`
// Host defines the configuration flags of a host
type Host struct {
	// ssh-config fields
	AddKeysToAgent                   string                    `yaml:"addkeystoagent,omitempty,flow" json:"AddKeysToAgent,omitempty" passthrough:""`
	AddressFamily                    string                    `yaml:"addressfamily,omitempty,flow" json:"AddressFamily,omitempty" passthrough:""`
	AskPassGUI                       string                    `yaml:"askpassgui,omitempty,flow" json:"AskPassGUI,omitempty" passthrough:""`
	BatchMode                        string                    `yaml:"batchmode,omitempty,flow" json:"BatchMode,omitempty" passthrough:""`
	BindAddress                      string                    `yaml:"bindaddress,omitempty,flow" json:"BindAddress,omitempty" passthrough:""`
	CanonicalDomains                 string                    `yaml:"canonicaldomains,omitempty,flow" json:"CanonicalDomains,omitempty" passthrough:""`
	CanonicalizeFallbackLocal        string                    `yaml:"canonicalizefallbacklocal,omitempty,flow" json:"CanonicalizeFallbackLocal,omitempty" passthrough:""`
	CanonicalizeHostname             string                    `yaml:"canonicalizehostname,omitempty,flow" json:"CanonicalizeHostname,omitempty" passthrough:""`
	CanonicalizeMaxDots              string                    `yaml:"canonicalizemaxDots,omitempty,flow" json:"CanonicalizeMaxDots,omitempty" passthrough:""`
	CanonicalizePermittedCNAMEs      string                    `yaml:"canonicalizepermittedcnames,omitempty,flow" json:"CanonicalizePermittedCNAMEs,omitempty" passthrough:""`
	CASignatureAlgorithms            composeyaml.Stringorslice `yaml:"casignaturealgorithms,omitempty,flow" json:"CASignatureAlgorithms,omitempty" passthrough:"joinByComma"`
	CertificateFile                  composeyaml.Stringorslice `yaml:"certificatefile,omitempty,flow" json:"CertificateFile,omitempty" passthrough:"perLine"`
	ChallengeResponseAuthentication  string                    `yaml:"challengeresponseauthentication,omitempty,flow" json:"ChallengeResponseAuthentication,omitempty" passthrough:""`
	CheckHostIP                      string                    `yaml:"checkhostip,omitempty,flow" json:"CheckHostIP,omitempty" passthrough:""`
	Cipher                           string                    `yaml:"cipher,omitempty,flow" json:"Cipher,omitempty" passthrough:""`
	Ciphers                          composeyaml.Stringorslice `yaml:"ciphers,omitempty,flow" json:"Ciphers,omitempty" passthrough:"joinByComma"`
	ClearAllForwardings              string                    `yaml:"clearallforwardings,omitempty,flow" json:"ClearAllForwardings,omitempty" passthrough:""`
	Compression                      string                    `yaml:"compression,omitempty,flow" json:"Compression,omitempty" passthrough:""`
	CompressionLevel                 int                       `yaml:"compressionlevel,omitempty,flow" json:"CompressionLevel,omitempty" passthrough:""`
	ConnectionAttempts               string                    `yaml:"connectionattempts,omitempty,flow" json:"ConnectionAttempts,omitempty" passthrough:""`
	ConnectTimeout                   int                       `yaml:"connecttimeout,omitempty,flow" json:"ConnectTimeout,omitempty" passthrough:""`
	ControlMaster                    string                    `yaml:"controlmaster,omitempty,flow" json:"ControlMaster,omitempty" passthrough:""`
	ControlPath                      string                    `yaml:"controlpath,omitempty,flow" json:"ControlPath,omitempty" passthrough:""`
	ControlPersist                   string                    `yaml:"controlpersist,omitempty,flow" json:"ControlPersist,omitempty" passthrough:""`
	DynamicForward                   composeyaml.Stringorslice `yaml:"dynamicforward,omitempty,flow" json:"DynamicForward,omitempty" passthrough:"perLine"`
	EnableEscapeCommandline          string                    `yaml:"enableescapecommandline,omitempty,flow" json:"EnableEscapeCommandline,omitempty" passthrough:""`
	EnableSSHKeysign                 string                    `yaml:"enablesshkeysign,omitempty,flow" json:"EnableSSHKeysign,omitempty" passthrough:""`
	EscapeChar                       string                    `yaml:"escapechar,omitempty,flow" json:"EscapeChar,omitempty" passthrough:""`
	ExitOnForwardFailure             string                    `yaml:"exitonforwardfailure,omitempty,flow" json:"ExitOnForwardFailure,omitempty" passthrough:""`
	FingerprintHash                  string                    `yaml:"fingerprinthash,omitempty,flow" json:"FingerprintHash,omitempty" passthrough:""`
	ForwardAgent                     string                    `yaml:"forwardagent,omitempty,flow" json:"ForwardAgent,omitempty" passthrough:""`
	ForwardX11                       string                    `yaml:"forwardx11,omitempty,flow" json:"ForwardX11,omitempty" passthrough:""`
	ForwardX11Timeout                int                       `yaml:"forwardx11timeout,omitempty,flow" json:"ForwardX11Timeout,omitempty" passthrough:""`
	ForwardX11Trusted                string                    `yaml:"forwardx11trusted,omitempty,flow" json:"ForwardX11Trusted,omitempty" passthrough:""`
	GatewayPorts                     string                    `yaml:"gatewayports,omitempty,flow" json:"GatewayPorts,omitempty" passthrough:""`
	GlobalKnownHostsFile             composeyaml.Stringorslice `yaml:"globalknownhostsfile,omitempty,flow" json:"GlobalKnownHostsFile,omitempty" passthrough:"joinByComma"`
	GSSAPIAuthentication             string                    `yaml:"gssapiauthentication,omitempty,flow" json:"GSSAPIAuthentication,omitempty" passthrough:""`
	GSSAPIClientIdentity             string                    `yaml:"gssapiclientidentity,omitempty,flow" json:"GSSAPIClientIdentity,omitempty" passthrough:""`
	GSSAPIDelegateCredentials        string                    `yaml:"gssapidelegatecredentials,omitempty,flow" json:"GSSAPIDelegateCredentials,omitempty" passthrough:""`
	GSSAPIKeyExchange                string                    `yaml:"gssapikeyexchange,omitempty,flow" json:"GSSAPIKeyExchange,omitempty" passthrough:""`
	GSSAPIRenewalForcesRekey         string                    `yaml:"gssapirenewalforcesrekey,omitempty,flow" json:"GSSAPIRenewalForcesRekey,omitempty" passthrough:""`
	GSSAPIServerIdentity             string                    `yaml:"gssapiserveridentity,omitempty,flow" json:"GSSAPIServerIdentity,omitempty" passthrough:""`
	GSSAPITrustDNS                   string                    `yaml:"gssapitrustdns,omitempty,flow" json:"GSSAPITrustDNS,omitempty" passthrough:""`
	HashKnownHosts                   string                    `yaml:"hashknownhosts,omitempty,flow" json:"HashKnownHosts,omitempty" passthrough:""`
	HostbasedAuthentication          string                    `yaml:"hostbasedauthentication,omitempty,flow" json:"HostbasedAuthentication,omitempty" passthrough:""`
	HostbasedKeyTypes                string                    `yaml:"hostbasedkeytypes,omitempty,flow" json:"HostbasedKeyTypes,omitempty" passthrough:""`
	HostKeyAlgorithms                composeyaml.Stringorslice `yaml:"hostkeyalgorithms,omitempty,flow" json:"HostKeyAlgorithms,omitempty" passthrough:"joinByComma"`
	HostKeyAlias                     string                    `yaml:"hostkeyalias,omitempty,flow" json:"HostKeyAlias,omitempty" passthrough:""`
	IdentitiesOnly                   string                    `yaml:"identitiesonly,omitempty,flow" json:"IdentitiesOnly,omitempty" passthrough:""`
	IdentityAgent                    string                    `yaml:"identityagent,omitempty,flow" json:"IdentityAgent,omitempty" passthrough:""`
	IdentityFile                     composeyaml.Stringorslice `yaml:"identityfile,omitempty,flow" json:"IdentityFile,omitempty" passthrough:"perLine"`
	IgnoreUnknown                    string                    `yaml:"ignoreunknown,omitempty,flow" json:"IgnoreUnknown,omitempty" passthrough:""`
	IPQoS                            composeyaml.Stringorslice `yaml:"ipqos,omitempty,flow" json:"IPQoS,omitempty" passthrough:"joinByComma"`
	KbdInteractiveAuthentication     string                    `yaml:"kbdinteractiveauthentication,omitempty,flow" json:"KbdInteractiveAuthentication,omitempty" passthrough:""`
	KbdInteractiveDevices            composeyaml.Stringorslice `yaml:"kbdinteractivedevices,omitempty,flow" json:"KbdInteractiveDevices,omitempty" passthrough:"joinByComma"`
	KexAlgorithms                    composeyaml.Stringorslice `yaml:"kexalgorithms,omitempty,flow" json:"KexAlgorithms,omitempty" passthrough:"joinByComma"`
	KeychainIntegration              string                    `yaml:"keychainintegration,omitempty,flow" json:"KeychainIntegration,omitempty" passthrough:""`
	LocalCommand                     string                    `yaml:"localcommand,omitempty,flow" json:"LocalCommand,omitempty" passthrough:""`
	RemoteCommand                    string                    `yaml:"remotecommand,omitempty,flow" json:"RemoteCommand,omitempty" passthrough:""`
	LocalForward                     composeyaml.Stringorslice `yaml:"localforward,omitempty,flow" json:"LocalForward,omitempty" passthrough:"perLine"`
	LogLevel                         string                    `yaml:"loglevel,omitempty,flow" json:"LogLevel,omitempty" passthrough:""`
	MACs                             composeyaml.Stringorslice `yaml:"macs,omitempty,flow" json:"MACs,omitempty" passthrough:"joinByComma"`
	Match                            string                    `yaml:"match,omitempty,flow" json:"Match,omitempty" passthrough:""`
	NoHostAuthenticationForLocalhost string                    `yaml:"nohostauthenticationforlocalhost,omitempty,flow" json:"NoHostAuthenticationForLocalhost,omitempty" passthrough:""`
	NumberOfPasswordPrompts          string                    `yaml:"numberofpasswordprompts,omitempty,flow" json:"NumberOfPasswordPrompts,omitempty" passthrough:""`
	PasswordAuthentication           string                    `yaml:"passwordauthentication,omitempty,flow" json:"PasswordAuthentication,omitempty" passthrough:""`
	PermitLocalCommand               string                    `yaml:"permitlocalcommand,omitempty,flow" json:"PermitLocalCommand,omitempty" passthrough:""`
	PKCS11Provider                   string                    `yaml:"pkcs11provider,omitempty,flow" json:"PKCS11Provider,omitempty" passthrough:""`
	Port                             string                    `yaml:"port,omitempty,flow" json:"Port,omitempty" passthrough:""`
	PreferredAuthentications         string                    `yaml:"preferredauthentications,omitempty,flow" json:"PreferredAuthentications,omitempty" passthrough:""`
	Protocol                         composeyaml.Stringorslice `yaml:"protocol,omitempty,flow" json:"Protocol,omitempty" passthrough:"joinByComma"`
	ProxyJump                        string                    `yaml:"proxyjump,omitempty,flow" json:"ProxyJump,omitempty" passthrough:""`
	ProxyUseFdpass                   string                    `yaml:"proxyusefdpass,omitempty,flow" json:"ProxyUseFdpass,omitempty" passthrough:""`
	PubkeyAcceptedAlgorithms         string                    `yaml:"pubkeyacceptedalgorithms,omitempty,flow" json:"PubkeyAcceptedAlgorithms,omitempty" passthrough:""`
	PubkeyAcceptedKeyTypes           string                    `yaml:"pubkeyacceptedkeytypes,omitempty,flow" json:"PubkeyAcceptedKeyTypes,omitempty" passthrough:""`
	PubkeyAuthentication             string                    `yaml:"pubkeyauthentication,omitempty,flow" json:"PubkeyAuthentication,omitempty" passthrough:""`
	RekeyLimit                       string                    `yaml:"rekeylimit,omitempty,flow" json:"RekeyLimit,omitempty" passthrough:""`
	RemoteForward                    composeyaml.Stringorslice `yaml:"remoteforward,omitempty,flow" json:"RemoteForward,omitempty" passthrough:"perLine"`
	RequestTTY                       string                    `yaml:"requesttty,omitempty,flow" json:"RequestTTY,omitempty" passthrough:""`
	RevokedHostKeys                  string                    `yaml:"revokedhostkeys,omitempty,flow" json:"RevokedHostKeys,omitempty" passthrough:""`
	RhostsRSAAuthentication          string                    `yaml:"rhostsrsaauthentication,omitempty,flow" json:"RhostsRSAAuthentication,omitempty" passthrough:""`
	RSAAuthentication                string                    `yaml:"rsaauthentication,omitempty,flow" json:"RSAAuthentication,omitempty" passthrough:""`
	SendEnv                          composeyaml.Stringorslice `yaml:"sendenv,omitempty,flow" json:"SendEnv,omitempty" passthrough:"perLine"`
	ServerAliveCountMax              int                       `yaml:"serveralivecountmax,omitempty,flow" json:"ServerAliveCountMax,omitempty" passthrough:""`
	ServerAliveInterval              int                       `yaml:"serveraliveinterval,omitempty,flow" json:"ServerAliveInterval,omitempty" passthrough:""`
	StreamLocalBindMask              string                    `yaml:"streamlocalbindmask,omitempty,flow" json:"StreamLocalBindMask,omitempty" passthrough:""`
	StreamLocalBindUnlink            string                    `yaml:"streamlocalbindunlink,omitempty,flow" json:"StreamLocalBindUnlink,omitempty" passthrough:""`
	StrictHostKeyChecking            string                    `yaml:"stricthostkeychecking,omitempty,flow" json:"StrictHostKeyChecking,omitempty" passthrough:""`
	TCPKeepAlive                     string                    `yaml:"tcpkeepalive,omitempty,flow" json:"TCPKeepAlive,omitempty" passthrough:""`
	Tunnel                           string                    `yaml:"tunnel,omitempty,flow" json:"Tunnel,omitempty" passthrough:""`
	TunnelDevice                     string                    `yaml:"tunneldevice,omitempty,flow" json:"TunnelDevice,omitempty" passthrough:""`
	UpdateHostKeys                   string                    `yaml:"updatehostkeys,omitempty,flow" json:"UpdateHostKeys,omitempty" passthrough:""`
	UseKeychain                      string                    `yaml:"usekeychain,omitempty,flow" json:"UseKeychain,omitempty" passthrough:""`
	UsePrivilegedPort                string                    `yaml:"useprivilegedport,omitempty,flow" json:"UsePrivilegedPort,omitempty" passthrough:""`
	User                             string                    `yaml:"user,omitempty,flow" json:"User,omitempty" passthrough:""`
	UserKnownHostsFile               composeyaml.Stringorslice `yaml:"userknownhostsfile,omitempty,flow" json:"UserKnownHostsFile,omitempty" passthrough:"joinByComma"`
	VerifyHostKeyDNS                 string                    `yaml:"verifyhostkeydns,omitempty,flow" json:"VerifyHostKeyDNS,omitempty" passthrough:""`
	VisualHostKey                    string                    `yaml:"visualhostkey,omitempty,flow" json:"VisualHostKey,omitempty" passthrough:""`
	XAuthLocation                    string                    `yaml:"xauthlocation,omitempty,flow" json:"XAuthLocation,omitempty" passthrough:""`

	// ssh-config fields with a different behavior
	HostName     string `yaml:"hostname,omitempty,flow" json:"HostName,omitempty"`
	ProxyCommand string `yaml:"proxycommand,omitempty,flow" json:"ProxyCommand,omitempty"`

	// exposed assh fields
	Inherits              composeyaml.Stringorslice `yaml:"inherits,omitempty,flow" json:"Inherits,omitempty"`
	Gateways              composeyaml.Stringorslice `yaml:"gateways,omitempty,flow" json:"Gateways,omitempty"`
	ResolveNameservers    composeyaml.Stringorslice `yaml:"resolvenameservers,omitempty,flow" json:"ResolveNameservers,omitempty"`
	ResolveCommand        string                    `yaml:"resolvecommand,omitempty,flow" json:"ResolveCommand,omitempty"`
	ControlMasterMkdir    string                    `yaml:"controlmastermkdir,omitempty,flow" json:"ControlMasterMkdir,omitempty"`
	Aliases               composeyaml.Stringorslice `yaml:"aliases,omitempty,flow" json:"Aliases,omitempty"`
	Hooks                 *HostHooks                `yaml:"hooks,omitempty,flow" json:"Hooks,omitempty"`
	Comment               composeyaml.Stringorslice `yaml:"comment,omitempty,flow" json:"Comment,omitempty"`
	RateLimit             string                    `yaml:"ratelimit,omitempty,flow" json:"RateLimit,omitempty"`
	GatewayConnectTimeout int                       `yaml:"gatewayconnecttimeout,omitempty,flow" json:"GatewayConnectTimeout,omitempty"`

	// private assh fields
	noAutomaticRewrite bool
	knownHosts         []string
	pattern            string
	name               string
	inputName          string
	isDefault          bool
	isTemplate         bool
	inherited          map[string]bool
}

// NewHost returns a host with name
func NewHost(name string) *Host {
	return &Host{
		name: name,
	}
}

// Validate checks for values errors
func (h *Host) Validate() []error {
	errs := []error{}

	switch cleanupValue(h.AddressFamily) {
	case "", "any", "inet", "inet4", "inet6":
		break
	default:
		errs = append(errs, fmt.Errorf("%q: invalid value for 'AddressFamily': %q", h.name, h.ControlMaster))
	}

	switch cleanupValue(h.ControlMaster) {
	case "", "yes", "no", "ask", "auto", "autoask":
		break
	default:
		errs = append(errs, fmt.Errorf("%q: invalid value for 'ControlMaster': %q", h.name, h.ControlMaster))
	}

	return errs
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
		if isDynamicHostname(h.name) {
			hostname = "[dynamic]"
		} else {
			hostname = h.name
		}
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

// RawName returns the raw name of a host without pattern computing
func (h *Host) RawName() string {
	return h.pattern
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
// nolint:gocyclo
func (h *Host) Options() OptionsList {
	options := make(OptionsList, 0)

	// ssh-config fields

	hostReflect := reflect.ValueOf(h).Elem()
	hostTypeReflect := reflect.TypeOf(h).Elem()

	// go through each field in Host. If 'passthrough' tag exists, process it.
	for i := 0; i < hostReflect.NumField(); i++ {
		field := hostReflect.Field(i)
		fieldType := hostTypeReflect.Field(i)

		fieldName := fieldType.Name
		passthroughTag, found := fieldType.Tag.Lookup("passthrough")
		if !found {
			continue
		}

		switch field.Type().Kind() {
		case reflect.String:
			if fieldVal := field.String(); fieldVal == "" {
				options = append(options, Option{Name: fieldName, Value: fieldVal})
			}
		case reflect.Int:
			if fieldVal := field.Int(); fieldVal != 0 {
				options = append(options, Option{Name: fieldName, Value: fmt.Sprintf("%d", fieldVal)})
			}
		case reflect.TypeOf(composeyaml.Stringorslice(nil)).Kind():
			fieldVal := field.Interface().(composeyaml.Stringorslice)

			if passthroughTag == "joinByComma" {
				if len(fieldVal) > 0 {
					options = append(options, Option{Name: fieldName, Value: strings.Join(fieldVal, " ")})
				}
			} else if passthroughTag == "perLine" {
				for _, entry := range fieldVal {
					options = append(options, Option{Name: fieldName, Value: entry})
				}
			} else {
				panic(fmt.Sprintf("Undefined Range type for field '%s'.\nIt must contains either of the following value: [ joinByComma, perLine ]", fieldName))
			}
		default:
			panic(fmt.Sprintf("Unimplemented datatype '%s' for field '%s'", field.Type(), fieldName))
		}
	}

	// ssh-config fields with a different behavior
	// HostName
	// ProxyCommand

	// exposed assh fields
	// Inherits
	// Gateways
	// ResolveNameservers
	// ResolveCommand
	// ControlMasterMkdir
	// Aliases
	// Comment
	// Hooks

	// private assh fields
	// knownHosts
	// pattern
	// name
	// inputName
	// isDefault
	// isTemplate
	// inherited

	return options
}

// minimal host preparation after loading
func (h *Host) prepare() {
	for key, name := range h.Inherits {
		h.Inherits[key] = strings.ToLower(name)
	}
	for key, name := range h.Gateways {
		h.Gateways[key] = strings.ToLower(name)
	}
}

// ApplyDefaults ensures a Host is valid by filling the missing fields with defaults
// nolint:gocyclo
func (h *Host) ApplyDefaults(defaults *Host) {

	// ssh-config fields
	hostReflect := reflect.ValueOf(h).Elem()
	hostTypeReflect := reflect.TypeOf(h).Elem()

	defaultsReflecet := reflect.ValueOf(defaults).Elem()

	// go through each field in Host. If 'passthrough' tag exists, process it.
	for i := 0; i < hostReflect.NumField(); i++ {
		field := hostReflect.Field(i)
		fieldType := hostTypeReflect.Field(i)

		fieldName := fieldType.Name
		_, found := fieldType.Tag.Lookup("passthrough")
		if !found {
			continue
		}

		switch field.Type().Kind() {
		case reflect.String:
			fieldVal := field.String()

			if fieldVal == "" {
				fieldVal = defaultsReflecet.FieldByName(fieldName).String()
			}
			field.SetString(utils.ExpandField(fieldVal))
		case reflect.Int:
			fieldVal := field.Int()

			if fieldVal == 0 {
				field.SetInt(defaultsReflecet.FieldByName(fieldName).Int())
			}
		case reflect.TypeOf(composeyaml.Stringorslice(nil)).Kind():
			fieldVal := field.Interface().(composeyaml.Stringorslice)

			if len(fieldVal) == 0 {
				fieldVal = defaultsReflecet.FieldByName(fieldName).Interface().(composeyaml.Stringorslice)
			}
			field.Set(reflect.ValueOf(utils.ExpandSliceField(fieldVal)).Convert(field.Type()))
		default:
			panic(fmt.Sprintf("Unimplemented datatype '%s' for field '%s'", field.Type(), fieldName))
		}
	}

	// ssh-config fields with a different behavior
	if h.HostName == "" {
		h.HostName = defaults.HostName
	}
	h.HostName = utils.ExpandField(h.HostName)

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

	if h.ControlMasterMkdir == "" {
		h.ControlMasterMkdir = defaults.ControlMasterMkdir
	}
	h.ControlMasterMkdir = utils.ExpandField(h.ControlMasterMkdir)

	if len(h.Gateways) == 0 {
		h.Gateways = defaults.Gateways
	}
	// h.Gateways = utils.ExpandField(h.Gateways)

	if len(h.Aliases) == 0 {
		h.Aliases = defaults.Aliases
	}

	if len(h.Comment) == 0 {
		h.Comment = defaults.Comment
	}

	if len(h.RateLimit) == 0 {
		h.RateLimit = defaults.RateLimit
	}

	if h.GatewayConnectTimeout == 0 {
		h.GatewayConnectTimeout = defaults.GatewayConnectTimeout
	}

	if h.Hooks == nil {
		h.Hooks = defaults.Hooks
		if h.Hooks == nil {
			h.Hooks = &HostHooks{}
		}
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
// nolint:gocyclo
func (h *Host) WriteSSHConfigTo(w io.Writer) error {
	aliases := append([]string{h.Name()}, h.Aliases...)
	aliases = append(aliases, h.knownHosts...)
	aliasIdx := 0
	for _, alias := range aliases {
		// FIXME: skip complex patterns

		if aliasIdx > 0 {
			_, _ = fmt.Fprint(w, "\n")
		}

		_, _ = fmt.Fprintf(w, "Host %s\n", alias)

		hostReflect := reflect.ValueOf(h).Elem()
		hostTypeReflect := reflect.TypeOf(h).Elem()

		// go through each field in Host. If 'passthrough' tag exists, process it.
		for i := 0; i < hostReflect.NumField(); i++ {
			field := hostReflect.Field(i)
			fieldType := hostTypeReflect.Field(i)

			fieldName := fieldType.Name
			passthroughTag, found := fieldType.Tag.Lookup("passthrough")
			if !found {
				continue
			}

			switch field.Type().Kind() {
			case reflect.String:
				fieldVal := field.String()
				if fieldVal != "" {
					_, _ = fmt.Fprintf(w, "  %s %s\n", fieldName, fieldVal)
				}
			case reflect.Int:
				fieldVal := field.Int()
				if fieldVal != 0 {
					_, _ = fmt.Fprintf(w, "  %s %d\n", fieldName, fieldVal)
				}
			case reflect.TypeOf(composeyaml.Stringorslice(nil)).Kind():
				fieldVal := field.Interface().(composeyaml.Stringorslice)

				if passthroughTag == "joinByComma" {
					if len(fieldVal) > 0 {
						_, _ = fmt.Fprintf(w, "  %s %s\n", fieldName, strings.Join(fieldVal, ","))
					}
				} else if passthroughTag == "perLine" {
					for _, entry := range fieldVal {
						_, _ = fmt.Fprintf(w, "   %s %s\n", fieldName, entry)
					}
				} else {
					panic(fmt.Sprintf("Undefined Range type for field '%s'.\nIt must contains either of the following value: [ joinByComma, perLine ]", fieldName))
				}
			default:
				panic(fmt.Sprintf("Unimplemented datatype '%s' for field '%s'", field.Type(), fieldName))
			}
		}

		// ssh-config fields with a different behavior
		if h.isDefault {
			if h.noAutomaticRewrite {
				_, _ = fmt.Fprintf(w, "  ProxyCommand %s connect --no-rewrite --port=%%p %%h\n", asshBinaryPath)
			} else {
				_, _ = fmt.Fprintf(w, "  ProxyCommand %s connect --port=%%p %%h\n", asshBinaryPath)
			}
		} else {
			if h.ProxyCommand != "" {
				_, _ = fmt.Fprintf(w, "  # ProxyCommand %s\n", h.ProxyCommand)
			}
		}

		// assh fields
		if h.HostName != "" {
			_, _ = fmt.Fprint(w, stringComment("HostName", h.HostName))
		}
		if BoolVal(h.ControlMasterMkdir) {
			_, _ = fmt.Fprint(w, "  # ControlMasterMkdir: true\n")
		}
		if len(h.Inherits) > 0 {
			_, _ = fmt.Fprint(w, sliceComment("Inherits", h.Inherits))
		}
		if len(h.Gateways) > 0 {
			_, _ = fmt.Fprint(w, sliceComment("Gateways", h.Gateways))
		}
		if len(h.Comment) > 0 {
			_, _ = fmt.Fprint(w, sliceComment("Comment", h.Comment))
		}
		if h.GatewayConnectTimeout > 0 {
			_, _ = fmt.Fprint(w, stringComment("GatewayConnectTimeout", fmt.Sprintf("%d", h.GatewayConnectTimeout)))
		}
		if len(h.Aliases) > 0 {
			if aliasIdx == 0 {
				_, _ = fmt.Fprint(w, sliceComment("Aliases", h.Aliases))
			} else {
				_, _ = fmt.Fprint(w, stringComment("AliasOf", h.Name()))
			}
		}
		if h.Hooks.Length() > 0 {
			_, _ = fmt.Fprint(w, stringComment("Hooks", h.Hooks.String()))
		}
		if len(h.knownHosts) > 0 {
			if aliasIdx == 0 {
				_, _ = fmt.Fprint(w, sliceComment("KnownHosts", h.knownHosts))
			} else {
				_, _ = fmt.Fprint(w, stringComment("KnownHostOf", h.Name()))
			}
		}

		if len(h.ResolveNameservers) > 0 {
			_, _ = fmt.Fprint(w, sliceComment("ResolveNameservers", h.ResolveNameservers))
		}
		if h.ResolveCommand != "" {
			_, _ = fmt.Fprint(w, stringComment("ResolveCommand", h.ResolveCommand))
		}
		if h.RateLimit != "" {
			_, _ = fmt.Fprint(w, stringComment("RateLimit", h.RateLimit))
		}

		aliasIdx++
	}
	return nil
}

// ExpandString replaces elements in a format string with host variables.
func (h *Host) ExpandString(input string, gateway string) string {
	output := input

	// name of the host in config
	output = strings.ReplaceAll(output, "%name", h.Name())

	// original target host name specified on the command line
	output = strings.ReplaceAll(output, "%n", h.inputName)

	// target host name
	output = strings.ReplaceAll(output, "%h", h.HostName)

	// port
	output = strings.ReplaceAll(output, "%p", h.Port)

	// gateway
	output = strings.ReplaceAll(output, "%g", gateway)

	// FIXME: add
	//   %L -> first component of the local host name
	//   %l -> local host name
	//   %r -> remote login username
	//   %u -> username of the user running assh
	//   %r -> remote login username

	return output
}

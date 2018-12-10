# Changelog

## master (unreleased)

  * Remove invalid warn ([#283](https://github.com/moul/assh/issues/283))
  * Avoid double-connection when having chained gateways ([#285](https://github.com/moul/assh/pull/285)) by [@4wrxb](https://github.com/4wrxb)
  * Add `GatewayConnectTimeout` option ([#289](https://github.com/moul/assh/issues/289))
  * Fix hook parsing to support `{{.}}` ([#297](https://github.com/moul/assh/issues/297)) by [@ahhx](https://github.com/ahhx)
  * Support inherits on regex hosts ([#298](https://github.com/moul/assh/issues/298)) by [@alenn-m](https://github.com/alenn-m)
  * Print version in `assh info`
  * Switch to go modules
  * Switch go moul.io/assh canonical domain
  * Move changelog to its own file
  * Switch to golangci-lint for linting
  * Move old webapp to contrib
  * Cleanup old files
  * Bump deps
  * Switch to go.uber.org/zap for logging

[Full commits list](https://github.com/moul/assh/compare/v2.8.0...master)

## v2.8.0 (2018-07-03)

  * Add more shell propositions for the 'exec' hook ([#254](https://github.com/moul/assh/issues/254))
  * Support SSH tokens and ~ expansion in ControlPaths ([#276](https://github.com/moul/assh/pull/276)) by [@stk0vrfl0w](https://github.com/stk0vrfl0w)
  * Ensure ControlPath directories are properly created when using syntax such as "ssh host1/host2" ([#276](https://github.com/moul/assh/pull/276)) by [@stk0vrfl0w](https://github.com/stk0vrfl0w)
  * Change panic() to a warning statement when removing the temporary file. Since delete is deferred, the config file should have already been renamed and would no longer exist ([#276](https://github.com/moul/assh/pull/276)) by [@stk0vrfl0w](https://github.com/stk0vrfl0w)
  * Fix check when ControlPath is empty, to avoid creating control socket directory in some cases ([#281](https://github.com/moul/assh/pull/281)) by [@frezbo](https://github.com/frezbo)

[Full commits list](https://github.com/moul/assh/compare/v2.7.0...v2.8.0)

## v2.7.0 (2017-10-24)

  * Automatically detect available shell when using the 'exec' hook ([#254](https://github.com/moul/assh/issues/254))
  * Automatically detect if `-q` is passed to the parent ssh process to disable logging ([#253](https://github.com/moul/assh/pull/253)) by [@cao](https://github.com/cao)
  * Add a new `%g` (gateway) parameter to `ResolveCommand` and `ProxyCommand` ([#247](https://github.com/moul/assh/pull/247)
  * Fix panic on particular `assh.yml` files
  * Fix build issue on OSX ([#238](https://github.com/moul/assh/pull/238) by [@jcftang](https://github.com/jcftang))
  * Rewrite .ssh/config file atomically ([#215](https://github.com/moul/assh/issues/215))
  * Support inheritance links in Graphviz ([#235](https://github.com/moul/assh/issues/235))
  * Support wildcards in Graphviz config export ([#228](https://github.com/moul/assh/issues/228))
  * Fix error message on first config build ([#230](https://github.com/moul/assh/issues/230))
  * Fix NetBSD, FreeBSD and Windows builds ([#182](https://github.com/moul/assh/issues/182))
  * Add RateLimit support ([#65](https://github.com/moul/assh/issues/65))
  * Add possibility to disable automatic configuration rewrite ([#239](https://github.com/moul/assh/issues/239))
  * Add `BeforeConfigWrite` and `AfterConfigWrite` new hooks ([#239](https://github.com/moul/assh/issues/239))
  * Generate full assh binary path in ~/.ssh/config ([#148](https://github.com/moul/assh/issues/148))
  * Initial version of `assh ping` command

[Full commits list](https://github.com/moul/assh/compare/v2.6.0...v2.7.0)

## v2.6.0 (2017-02-03)

  * Support `UseKeychain` option ([#189](https://github.com/moul/assh/pull/189)) ([@ocean90](https://github.com/ocean90))
  * Support `ConnectTimeout` option ([#132](https://github.com/moul/assh/issues/132))
  * `.ssh/config`: Wrap long comments to avoid syntax errors ([#191](https://github.com/moul/assh/issues/191))
  * Fix integers output in `assh config list` ([#181](https://github.com/moul/assh/issues/181))
  * Initial graphviz support ([#32](https://github.com/moul/assh/issues/32))
  * Remove case-sensitivity for `Inherits` and `Gateways` ([#178](https://github.com/moul/assh/issues/178))
  * Loads hosts from `~/.ssh/assh_known_hosts` file when calling `assh config build`, can be ignored using `--ignore-known-hosts` ([#178](https://github.com/moul/assh/issues/178))
  * Add `assh config graphviz --show-isolated-hosts` flag
  * Fix nil dereference when calling `assh config serach` without providing needle
  * Add [sprig](https://github.com/Masterminds/sprig) helpers to the template engine ([#206](https://github.com/moul/assh/issues/206))
  * Improve readability of `assh config list` ([#203](https://github.com/moul/assh/issues/203))
  * Add support for the `AddKeysToAgent` key ([#210](https://github.com/moul/assh/pull/210)) ([@bachya](https://github.com/bachya))
  * OpenBSD support ([#182](https://github.com/moul/assh/issues/182))
  * Improve hostname output in `assh config list` ([#204](https://github.com/moul/assh/issues/204))
  * Support for inline comments ([#34](https://github.com/moul/assh/issues/34))
  * Initial support of values validation to avoid writing invalid .ssh/config file ([#92](https://github.com/moul/assh/issues/92))
  * Alpha version of the webapp ([#69](https://github.com/moul/assh/issues/69))

[Full commits list](https://github.com/moul/assh/compare/v2.5.0...v2.6.0)

## v2.5.0 (2017-01-04)

  * Support multiple string arguments of the same type on `assh wrapper ssh` ([#185](https://github.com/moul/assh/issues/185))
  * Remove the `NoControlMasterMkdir` option, and add the `ControlMasterMkdir` option instead ([#173](https://github.com/moul/assh/issues/173))
  * Accepting string or slices for list options ([#119](https://github.com/moul/assh/issues/119))
  * Add new `PubkeyAcceptedKeyTypes` OpenSSH 7+ field ([#175](https://github.com/moul/assh/issues/175))
  * Gracefully report an error when calling assh without configuration file ([#171](https://github.com/moul/assh/issues/171))
  * Fix `written bytes` calculation ([@quentinperez](https://github.com/quentinperez))
  * Add template functions: `json`, `prettyjson`, `split`, `join`, `title`, `lower`, `upper`
  * Support of `BeforeConnect`, `OnConnect`, `OnConnectError` and `OnDisconnect` hooks
  * Support of `write`, `notify` and `exec` hook drivers
  * Add `assh config json` command
  * Add `assh config {build,json} --expand` option
  * Round the hook's `ConnectionDuration` variable value

[Full commits list](https://github.com/moul/assh/compare/v2.4.1...v2.5.0)

## v2.4.1 (2016-07-19)

  * Fix panic in `assh wrapper` ([#157](https://github.com/moul/assh/issues/157))

[Full commits list](https://github.com/moul/assh/compare/v2.4.0...v2.4.1)

## v2.4.0 (2016-07-14)

  * Add a control socket manager `assh sockets {list,flush,master}` ([#152](https://github.com/moul/assh/pull/152))
  * Add a `assh --config=/path/to/assh.yml` option
  * Add storm-like `assh config list` and `assh config search {keyword}` commands ([#151](https://github.com/moul/assh/pull/151))
  * Add an optional `ASSHBinaryPath` variable in the `assh.yml` file ([#148](https://github.com/moul/assh/issues/148))
  * Rename `assh proxy -> assh connect`
  * Hide `assh connect` and `assh wrapper` from the help
  * Support built-in ssh netcat mode, may fail with older SSH clients ([#146](https://github.com/moul/assh/issues/146))

[Full commits list](https://github.com/moul/assh/compare/v2.3.0...v2.4.0)

## v2.3.0 (2016-04-27)

  * Add wrapper and `known_hosts` support to handle *advanced patterns* ([#122](https://github.com/moul/assh/issues/122))
  * Add build information in .ssh/config header ([#49](https://github.com/moul/assh/issues/49))
  * Add Autocomplete support ([#48](https://github.com/moul/assh/issues/48))
  * Initial `Aliases` support ([#133](https://github.com/moul/assh/issues/133))
  * Use args[0] as ProxyCommand ([#134](https://github.com/moul/assh/issues/134))
  * Add `NoControlMasterMkdir` option to disable automatic creation of directories for gateways ([#124](https://github.com/moul/assh/issues/124))
  * Fix: Allow `$(...)` syntax in the `ResolveCommand` function ([#117](https://github.com/moul/assh/issues/117))
  * Printing the error of a failing `ResolveCommand` ([#117](https://github.com/moul/assh/issues/117))
  * Fix: `Gateways` field is no longer ignored when the `HostName` field is present ([#102](https://github.com/moul/assh/issues/102))
  * Ignore SIGHUP, close goroutines and export written bytes ([#112](https://github.com/moul/assh/pull/112)) ([@QuentinPerez](https://github.com/QuentinPerez))
  * Various documentation improvements ([@ashmatadeen](https://github.com/ashmatadeen), [@loliee](https://github.com/loliee), [@cerisier](https://github.com/cerisier))
  * Support of new SSH configuration fields (`AskPassGUI`, `GSSAPIClientIdentity`, `GSSAPIKeyExchange`, `GSSAPIRenewalForcesRekey`, `GSSAPIServerIdentity`, `GSSAPITrustDns`, `KeychainIntegration`)

[Full commits list](https://github.com/moul/assh/compare/v2.2.0...v2.3.0)

## v2.2.0 (2016-02-03)

  * Avoid exiting when an included file contains errors ([#95](https://github.com/moul/assh/issues/95))
  * Anonymize paths in `assh info`
  * Support of `assh proxy --dry-run` option
  * Fix: do not resolve variables in hostnames twice ([#103](https://github.com/moul/assh/issues/103))

[Full commits list](https://github.com/moul/assh/compare/v2.1.0...v2.2.0)

## v2.1.0 (2015-10-05)

  * Expand environment variables ([#86](https://github.com/moul/assh/issues/86))
  * Add homebrew support ([#73](https://github.com/moul/assh/issues/73))
  * Add a 'ssh info' command ([#71](https://github.com/moul/assh/issues/71))
  * Templates support ([#52](https://github.com/moul/assh/issues/52))
  * Configuration is now case insensitive ([#51](https://github.com/moul/assh/issues/51))
  * Fix: resolving host fields for gateways ([#79](https://github.com/moul/assh/issues/79))
  * Fix: inheritance was not working for non assh-related fields ([#54](https://github.com/moul/assh/issues/54))
  * Fix: expanding variables in HostName ([#56](https://github.com/moul/assh/issues/56))

[Full commits list](https://github.com/moul/assh/compare/v2.0.0...v2.1.0)

## v2.0.0 (2015-09-07)

  * First Golang version
  * Compatibility issue: complete switch from `.ini` file format to `.yml`, the `~/.ssh/assh.yml` file needs to be manually crafted
  * Features
    * Parses `~/.ssh/assh.yml` and generates `~/.ssh/config` dynamically
    * CLI: Use gateways from CLI without any configuration needed
    * Config: Declares gateways in coniguration
    * Config: Host inheritance
    * Config: Support of `includes`
    * Config: Support of Regex
    * Config: Handling all sshconfig fields
    * Config: Support of host `ProxyCommand` (inception)
    * Under the hood: Inspecting parent process **verbose**/**debug** mode
    * Under the hook: dynamic proxy using **raw TCP**, **netcat**

[Full commits list](https://github.com/moul/assh/compare/be4fea1632b1e9f8aa60585187338777baaf1210...v2.0.0)

## [v1](https://github.com/moul/assh/tree/v1.1.0) (2015-07-22)

  * Last Python version

## [POC](https://github.com/moul/assh/commit/550f86c225d30292728ad24bc883b6d3a3e3f1b1) (2010-08-26)

  * First Python version (POC)

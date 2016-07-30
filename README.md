# assh - *Advanced SSH config*
[![Travis](https://img.shields.io/travis/moul/advanced-ssh-config.svg)](https://travis-ci.org/moul/advanced-ssh-config)
[![GoDoc](https://godoc.org/github.com/moul/advanced-ssh-config?status.svg)](https://godoc.org/github.com/moul/advanced-ssh-config)
![License](https://img.shields.io/github/license/moul/advanced-ssh-config.svg)
[![GitHub release](https://img.shields.io/github/release/moul/advanced-ssh-config.svg)](https://github.com/moul/advanced-ssh-config/releases)

<img src="https://raw.githubusercontent.com/moul/advanced-ssh-config/master/resources/assh.png" width="400" />

#### Table of Contents

1. [Overview](#overview)
2. [Features](#features)
  * [Configuration features](#configuration-features)
  * [Using Gateway from command line](#using-gateway-from-command-line)
  * [Using Gateway from configuration file](#using-gateways-from-configuration-file)
  * [Under the hood features](#under-the-hood-features)
  * [Hooks](#hooks)
3. [Configuration](#configuration)
4. [Usage](#usage)
  * [Usage Examples](#usage-examples)
5. [Install](#install)
  * [Register the wrapper (optional)](#register-the-wrapper-optional)
6. [Getting started](#getting-started)
7. [Changelog](#changelog)
8. [Alternative version](#alternative-version)
9. [Troubleshooting](#troubleshooting)
10. [License](#license)

## Overview

A *transparent wrapper* that adds **regex**, **aliases**, **gateways**, **includes**, **dynamic hostnames** to **SSH**.

[lib-ssh](https://www.libssh.org) wraps `assh` as a [ProxyCommand](https://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#ProxyCommand_with_Netcat); it means that it works seamlessly with:

* [ssh](http://www.openbsd.org/cgi-bin/man.cgi?query=ssh&sektion=1)
* [scp](http://www.openbsd.org/cgi-bin/man.cgi?query=scp&sektion=1)
* [rsync](http://linuxcommand.org/man_pages/rsync1.html)
* [git](https://www.kernel.org/pub/software/scm/git/docs/)
* Desktop applications depending on `lib-ssh` or `ssh` (i.e., [Tower](http://www.git-tower.com), [Atom.io](https://atom.io), [SSH Tunnel Manager](http://projects.tynsoe.org/fr/stm/))

## Features

### Configuration features

* **regex** support
* **aliases** `gate` -> `gate.domain.tld`
* **gateways** -> transparent ssh connection chaining
* **includes**: split configuration in multiple files
* **local command execution**: finally the reverse of **RemoteCommand**
* **templates**: equivalent to host but you can't connect directly to a template, perfect for inheritance
* **inheritance**: make hosts inherits from host hosts or templates
* **variable expansion**: resolve variables from the environment
* **smart proxycommand**: RAW tcp connection when possible with `netcat` and `socat` as default fallbacks

### Using Gateway from command line

*assh* can use the [ProxyCommand with netcat](https://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#ProxyCommand_with_Netcat) feature of OpenSSH **transparently** and without the pain of using extended configuration.

Connect to `hosta` using `hostb` as a gateway.

```
  ┌─────┐
  │ YOU │─ ─ ─ ─ ─
  └─────┘         │
     ┃            ▽
     ┃         ┌─────┐
 firewall      │hostb│
     ┃         └─────┘
     ▼            │
  ┌─────┐
  │hosta│◁─ ─ ─ ─ ┘
  └─────┘
```

```console
$ ssh hosta/hostb
user@hosta $
```

Equivalent to `ssh -o ProxyCommand="ssh hostb nc %h %p" hosta`

---

Connect to `hosta` using `hostb` as a gateway using `hostc` as a gateway.

```
  ┌─────┐              ┌─────┐
  │ YOU │─ ─ ─ ─ ─ ─ ─▷│hostc│
  └─────┘              └─────┘
     ┃                    │
     ┃
 firewall                 │
     ┃
     ┃                    │
     ▼                    ▽
  ┌─────┐              ┌─────┐
  │hosta│◁─ ─ ─ ─ ─ ─ ─│hostb│
  └─────┘              └─────┘
```

```console
$ ssh hosta/hostb/hostc
user@hosta $
```

Equivalent to `ssh -o ProxyCommand="ssh -o ProxyCommand='ssh hostc nc %h %p' hostb nc %h %p" hosta`

### Using Gateways from configuration file

You can define an equivalent of the ["ProxyCommand with netcat"](https://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#ProxyCommand_with_Netcat) feature of OpenSSH, with a simpler syntax, more advanced workflows, and a unique fallback feature.

Let's consider the following `assh.yml` file
```yaml
hosts:
  hosta:
    Hostname: 1.2.3.4

  hostb:
    Hostname: 5.6.7.8
    Gateways:
    - hosta

  hostc:
    Hostname: 9.10.11.12
    Gateways:
    - hostb

  hostd:
    Hostname: 13.14.15.16
    Gateways:
    - direct
    - hosta
```

* `ssh hosta` -> `ssh 1.2.3.4`
* `ssh hostb` -> `ssh -o ProxyCommand="ssh hostb nc %h %p" hosta`
* `ssh hostc` -> `ssh -o ProxyCommand="ssh -o ProxyCommand='ssh hostc nc %h %p' hostb nc %h %p" hosta`
* `ssh hostd` ->
  * assh will try to `ssh 13.14.15.16`
  * then, fallback on `ssh -o ProxyCommand="ssh hostd nc %h %p" hosta`
  * this method allows you to have the best performances when it is possible, but ensure your commands will work if you are outside of your company for instance

### Under the hood features

* Automatically regenerates `~/.ssh/config` file when needed
* Inspect parent process to determine log level (if you use `ssh -vv`, **assh** will automatically run in debug mode)
* Automatically creates `ControlPath` directories so you can use *slashes* in your `ControlPath` option, can be disabled with the `NoControlMasterMkdir: true` configuration in host or globally.

### Hooks

#### Events

##### OnConnect

`OnConnect` is called as soon as assh is connected to the remote SSH port.

Note: `OnConnect` is not aware of the authentication process and will always be raised.

---

Example of Golang template variables:

```golang
// Host: http://godoc.org/github.com/moul/advanced-ssh-config/pkg/config/#Host
{{.Host.Name}}                                  //  localhost
{{.Host.HostName}}                              //  127.0.0.1
{{.Host.Port}}                                  //  22
{{.Host.User}}                                  //  moul
{{.Host.Prototype}}                             //  moul@127.0.0.1:22
{{.Host}}                                       //  {"HostName":"localhost","Port":22","User":"moul","ControlPersist":"yes",...}
{{printf "%s:%s" .Host.HostName .Host.Port}}    //  localhost:22

// Stats: http://godoc.org/github.com/moul/advanced-ssh-config/pkg/commands/#ConnectionStats
{{.Stats.ConnectedAt}}                           //  2016-07-20 11:19:23.467900594 +0200 CEST
{{.Stats.WrittenBytes}}                          //  3613
{{.Stats.WrittenBytesHuman}}                     //  3.6kb
{{.Stats.DisconnectAt}}                          //  2016-07-20 11:19:29,520515792 +0200 CEST
{{.Stats.ConnectionDuration}}                    //  6.052615198s
{{.Stats.ConnectionDurationHuman}}               //  6s
{{.Stats.AverageSpeed}}                          //  596.933bps
{{.Stats.AverageSpeedHuman}}                     //  3.4kb/s
```

##### OnDisconnect

`OnDisconnect` is called as the assh socket is closed.

---

Example of Golang template variables:

```golang
// Host: http://godoc.org/github.com/moul/advanced-ssh-config/pkg/config/#Host
{{.Host.Name}}                                  //  localhost
{{.Host.HostName}}                              //  127.0.0.1
{{.Host.Port}}                                  //  22
{{.Host.User}}                                  //  moul
{{.Host.Prototype}}                             //  moul@127.0.0.1:22
{{.Host}}                                       //  {"HostName":"localhost","Port":22","User":"moul","ControlPersist":"yes",...}
{{printf "%s:%s" .Host.HostName .Host.Port}}    //  localhost:22

// Stats: http://godoc.org/github.com/moul/advanced-ssh-config/pkg/commands/#ConnectionStats
{{.Stats.ConnectedAt}}                           //  2016-07-20 11:19:23.467900594 +0200 CEST
{{.Stats.WrittenBytes}}                          //  3613
{{.Stats.WrittenBytesHuman}}                     //  3.6kb
{{.Stats.DisconnectAt}}                          //  2016-07-20 11:19:29,520515792 +0200 CEST
{{.Stats.ConnectionDuration}}                    //  6.052615198s
{{.Stats.ConnectionDurationHuman}}               //  6s
{{.Stats.AverageSpeed}}                          //  596.933bps
{{.Stats.AverageSpeedHuman}}                     //  3.4kb/s
```

#### Hooks drivers

##### Exec driver

Exec driver uses [Golang's template system](https://golang.org/pkg/text/template/) to execute a shell command

Usage: `exec <binary> [args...]`

```yaml
defaults:
  Hooks:
    OnConnect:
    - exec echo '{{.Host}}' | jq .
# executes: `echo '{"HostName":"localhost","Port":"22","User":"moul","ControlPersist":"yes",...}' | jq .
# which results in printing a pretty JSON of the host
# {
#   "HostName": "localhost",
#   "Port": "22",
#   "User": "moul",
#   "ControlPersist": "yes",
#   ...
# }
```

```yaml
defaults:
  Hooks:
    OnConnect:
    - exec echo 'New SSH connection to {{.Host.Prototype}}.' | mail -s "SSH connection journal" m+assh@42.am
# send an email with the connection prototype
```

---

The `exec` commands are blocking, a new driver for background tasks is planned. For now, you can run a job in background like this:

```yaml
defaults:
  Hooks:
    OnConnect:
    - exec sleep 60 &
# execute the `sleep 60` command in background (non-blocking)
# if you quit your ssh connection, the process will continue in background.
```

##### Write driver

Write driver uses [Golang's template system](https://golang.org/pkg/text/template/) to write out data to stdout

Usage: `write <line:string...>`

```yaml
defaults:
  Hooks:
    OnConnect:
    - write New SSH connection to {{.Host.Prototype}}.
# writes: "New SSH connection to moul@127.0.0.1:22." on the terminal on connection
```

```yaml
defaults:
  Hooks:
    OnDisconnect:
    - "write SSH connection to {{.Host.Name}} closed, {{ .Stats.WrittenBytes }} bytes written in {{ .Stats.ConnectionDuration }} ({{ .Stats.AverageSpeed }})"
# writes: SSH connection to localhost closed, 40 bytes written.
```

##### Notify driver

Notify driver uses [Golang's template system](https://golang.org/pkg/text/template/) to open Desktop notifications.

Note: OS X only for now.

Usage: `notify <line:string...>`

```yaml
defaults:
  Hooks:
    OnConnect:
    - notify New SSH connection to {{.Host.Prototype}}.
```

![](https://github.com/moul/advanced-ssh-config/raw/master/resources/new_connection_notification.png)

```yaml
defaults:
  Hooks:
    OnDisconnect:
    - "notify SSH connection to {{.Host.Name}} closed, {{ .Stats.WrittenBytes }} bytes written in {{ .Stats.ConnectionDuration }} ({{ .Stats.AverageSpeed }})"
```

![](https://github.com/moul/advanced-ssh-config/raw/master/resources/closed_connection_notification.png)

## Configuration

`assh` now manages the `~/.ssh/config` file, take care to keep a backup your `~/.ssh/config` file.

`~/.ssh/assh.yml` is a [YAML](http://www.yaml.org/spec/1.2/spec.html) file containing:

* a `hosts` dictionary containing multiple *HOST* definitions
* a `defaults` section containing global flags
* and an `includes` section containing path to other configuration files

```yaml
hosts:

  homer:
    # ssh homer ->  ssh 1.2.3.4 -p 2222 -u robert
    Hostname: 1.2.3.4
    User: robert
    Port: 2222

  bart:
    # ssh bart ->   ssh 5.6.7.8 -u bart           <- direct access
    #            or ssh 5.6.7.8/homer -u bart     <- using homer as a gateway
    Hostname: 5.6.7.8
    User: bart
    Gateways:
    - direct                   # tries a direct access first
    - homer                    # fallback on homer gateway

  maggie:
    # ssh maggie ->   ssh 5.6.7.8 -u maggie       <- direct access
    #              or ssh 5.6.7.8/homer -u maggie   <- using homer as a gateway
    User: maggie
    Inherits:
    - bart                     # inherits rules from "bart"

  bart-access:
    # ssh bart-access ->  ssh home.simpson.springfield.us -u bart
    Inherits:
    - bart-template
    - simpson-template

  lisa-access:
    # ssh lisa-access ->  ssh home.simpson.springfield.us -u lisa
    Inherits:
    - lisa-template
    - simpson-template

  marvin:
    # ssh marvin    -> ssh marvin    -p 23
    # ssh sad-robot -> ssh sad-robot -p 23
    # ssh bighead   -> ssh bighead   -p 23
    # aliases inherit everything from marvin, except hostname
    Port: 23
    Aliases:
    - sad-robot
    - bighead

  dolphin:
    # ssh dolphin   -> ssh dolphin -p 24
    # ssh ecco      -> ssh dolphin -p 24
    # same as above, but with fixed hostname
    Port: 24
    Hostname: dolphin
    Aliases:
    - ecco

  schooltemplate:
    User: student
    IdentityFile: ~/.ssh/school-rsa
    ForwardX11: yes

  schoolgw:
    # ssh school ->   ssh gw.school.com -l student -o ForwardX11=no -i ~/.ssh/school-rsa
    Hostname: gw.school.com
    ForwardX11: no
    Inherits:
    - schooltemplate

  "expanded-host[0-7]*":
    # ssh somehost2042 ->       ssh somehost2042.some.zone
    Hostname: "%h.some.zone"

  vm-*.school.com:
    # ssh vm-42.school.com ->   ssh vm-42.school.com/gw.school.com -l student -o ForwardX11=yes -i ~/.ssh/school-rsa
    Gateways:
    - schoolgw
    Inherits:
    - schooltemplate
    # do not automatically create `ControlPath` -> may result in error
    NoControlMasterMkdir: true

  "*.shortcut1":
    ResolveCommand: /bin/sh -c "echo %h | sed s/.shortcut1/.my-long-domain-name.com/"

  "*.shortcut2":
    ResolveCommand: /bin/sh -c "echo $(echo %h | sed s/.shortcut2//).my-other-long-domain-name.com"

  "*.scw":
    # ssh toto.scw -> 1. dynamically resolves the IP address
    #                 2. ssh {resolved ip address} -u root -p 22 -o UserKnownHostsFile=null -o StrictHostKeyChecking=no
    # requires github.com/scaleway/scaleway-cli
    ResolveCommand: /bin/sh -c "scw inspect -f {{.PublicAddress.IP}} server:$(echo %h | sed s/.scw//)"
    User: root
    Port: 22
    UserKnownHostsFile: /dev/null
    StrictHostKeyChecking: no

  my-env-host:
    User: user-$USER
    Hostname: ${HOSTNAME}${HOSTNAME_SUFFIX}

templates:
  # Templates are similar to Hosts; you can inherit from them
  # but you cannot ssh to a template
  bart-template:
    User: bart
  lisa-template:
    User: lisa
  simpson-template:
    Host: home.simpson.springfield.us

defaults:
  # Defaults are applied to each hosts
  ControlMaster: auto
  ControlPath: ~/tmp/.ssh/cm/%h-%p-%r.sock
  ControlPersist: yes
  Port: 22
  User: bob

includes:
- ~/.ssh/assh.d/*.yml
- /etc/assh.yml
- $ENV_VAR/blah-blah-*/*.yml

ASSHBinaryPath: ~/bin/assh  # optionally set the path of assh
```

---

A *HOST* and the `defaults` section may



## Usage

`assh` usage

```
NAME:
   assh - advanced ssh config

USAGE:
   assh [global options] command [command options] [arguments...]

VERSION:
   2.4.1 (HEAD)

AUTHOR(S):
   Manfred Touron <https://github.com/moul/advanced-ssh-config>

COMMANDS:
   info          Display system-wide information
   config        Manage ssh and assh configuration
   sockets       Manage control sockets
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
  --config value, -c value       Location of config file (default: "~/.ssh/assh.yml") [$ASSH_CONFIG]
  --debug, -D                    Enable debug mode [$ASSH_DEBUG]
  --verbose, -V                  Enable verbose mode
  --help, -h                     show help
  --version, -v                  print the version
```

### Usage examples

##### `assh config build`

Rewrites and replaces the existing ~/.ssh/config file.

This action is automatically done by assh when detecting configuration changes.
Running this command is useful to set up assh or repair the configuration file.

```console
$ assh config build > ~/.ssh/config
```

##### `assh config list`

List hosts and options.

```console
$ assh config list
Listing entries

    *.scw -> root@[hostname_not_specified]:22
        [custom options] StrictHostKeyChecking=no UserKnownHostsFile=/dev/null

    *.shortcut1 -> bob@[hostname_not_specified]:22

    *.shortcut2 -> bob@[hostname_not_specified]:22

    bart -> bart@5.6.7.8:22

    bart-access -> bob@[hostname_not_specified]:22

    dolphin -> bob@dolphin:24

    expanded-host[0-7]* -> bob@%h.some.zone:22

    homer -> robert@1.2.3.4:2222

    lisa-access -> bob@[hostname_not_specified]:22

    maggie -> maggie@[hostname_not_specified]:22

    marvin -> bob@[hostname_not_specified]:23

    my-env-host -> user-moul@[hostname_not_specified]:22

    schoolgw -> bob@gw.school.com:22
        [custom options] ForwardX11=no

    schooltemplate -> student@[hostname_not_specified]:22
        [custom options] ForwardX11=yes IdentityFile=~/.ssh/school-rsa

    vm-*.school.com -> bob@[hostname_not_specified]:22

    (*) General options:
        ControlMaster: auto
        ControlPath: ~/tmp/.ssh/cm/%h-%p-%r.sock
        ControlPersist: yes
        Port: 22
        User: bob
```

##### `assh config search <keyword>`

Search for `<keyword>` in hosts and host options.

```console
$ assh config search bart
Listing results for bart:
    bart -> bart@5.6.7.8:22
    bart-access -> moul@[hostname_not_specified]:22
```

##### `assh info`

Display system-wide information.

```console
$ assh info
Debug mode (client): false
CLI Path: /path/to/assh
Go version: go1.6.2
OS/Arch: darwin/amd64

RC files:
- ~/.ssh/assh.yml
- ~/.ssh/assh.d/hosts.yml
- ~/.ssh/assh.d/moul.yml
- ~/.ssh/assh.d/test.yml

Statistics:
- 299 hosts
- 2 templates
- 4 included files
```

##### `assh sockets list`

List active control sockets.

```console
$ assh sockets list
4 active control sockets in "~/.ssh/cm/":

- bart/homer/lisa-22-root.sock (14 minutes)
- bart/homer-22-root.sock (14 minutes)
- bart-22-root.sock (14 minutes)
- marge-22-bart.sock (1 hour)
```

##### `assh sockets flush`

Close active control sockets.

```console
$ assh sockets flush
Closed 4 control sockets.
```

##### `assh sockets master`

Create a master control sockets.

```console
$ assh sockets master
```

## Install

Get the latest version using GO (recommended way):

```bash
go get -u github.com/moul/advanced-ssh-config/cmd/assh
```

---

Get the latest released version using homebrew (Mac OS X):

```bash
brew install assh
```

Build the latest version

```bash
brew install assh --HEAD
```

---

Get a released version on: https://github.com/moul/advanced-ssh-config/releases

---

#### Register the wrapper (optional)

To improve experience when using advanced pattern matching, add the following at the end of your .bashrc or .zshrc:

```bash
alias ssh="assh wrapper ssh"
```

This step is not *mandatory* but highly *recommended*.

---

**Note**: `ssh` does not understand advanced patterns;
To bypass this limitation, `assh` maintains a list of *known hosts* and regenerate the `~/.ssh/config` with all those expanded *known hosts*.

Without the wrapper, the `~/.ssh/config` risks to be outdated when connecting to a new host for the first time and you will need to launch the command again.

With the wrapper, `ssh` will *always* be called with an updated `~/.ssh/config` file.

## Getting started

1. Backup your old `~/.ssh/config`: `cp ~/.ssh/config ~/.ssh/config.backup`
2. Create a new `~/.ssh/assh.yml` file
3. Run `assh config build > ~/.ssh/config` to validate the syntax of your `~/.ssh/assh.yml` file and automatically build your `~/.ssh/config` file
4. You are ready!

## Changelog

### master (unreleased)

* Fix `written bytes` calculation ([@quentinperez](https://github.com/quentinperez))
* Add template functions: `json`, `prettyjson`, `split`, `join`, `title`, `lower`, `upper`
* Support of `OnConnect` and `OnDisconnect` hooks
* Support of `write`, `notify` and `exec` hook drivers
* Add `assh config json` command
* Add `assh config {build,json} --expand` option

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.4.1...master)

### v2.4.1 (2016-07-19)

* Fix panic in `assh wrapper` ([#157](https://github.com/moul/advanced-ssh-config/issues/157))

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.4.0...v2.4.1)

### v2.4.0 (2016-07-14)

* Add a control socket manager `assh sockets {list,flush,master}` ([#152](https://github.com/moul/advanced-ssh-config/pull/152))
* Add a `assh --config=/path/to/assh.yml` option
* Add storm-like `assh config list` and `assh config search {keyword}` commands ([#151](https://github.com/moul/advanced-ssh-config/pull/151))
* Add an optional `ASSHBinaryPath` variable in the `assh.yml` file ([#148](https://github.com/moul/advanced-ssh-config/issues/148))
* Rename `assh proxy -> assh connect`
* Hide `assh connect` and `assh wrapper` from the help
* Support built-in ssh netcat mode, may fail with older SSH clients ([#146](https://github.com/moul/advanced-ssh-config/issues/146))

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.3.0...v2.4.0)

### v2.3.0 (2016-04-27)

* Add wrapper and `known_hosts` support to handle *advanced patterns* ([#122](https://github.com/moul/advanced-ssh-config/issues/122))
* Add build information in .ssh/config header ([#49](https://github.com/moul/advanced-ssh-config/issues/49))
* Add Autocomplete support ([#48](https://github.com/moul/advanced-ssh-config/issues/48))
* Initial `Aliases` support ([#133](https://github.com/moul/advanced-ssh-config/issues/133))
* Use args[0] as ProxyCommand ([#134](https://github.com/moul/advanced-ssh-config/issues/134))
* Add `NoControlMasterMkdir` option to disable automatic creation of directories for gateways ([#124](https://github.com/moul/advanced-ssh-config/issues/124))
* Fix: Allow `$(...)` syntax in the `ResolveCommand` function ([#117](https://github.com/moul/advanced-ssh-config/issues/117))
* Printing the error of a failing `ResolveCommand` ([#117](https://github.com/moul/advanced-ssh-config/issues/117))
* Fix: `Gateways` field is no longer ignored when the `HostName` field is present ([#102](https://github.com/moul/advanced-ssh-config/issues/102))
* Ignore SIGHUP, close goroutines and export written bytes ([#112](https://github.com/moul/advanced-ssh-config/pull/112)) ([@QuentinPerez](https://github.com/QuentinPerez))
* Various documentation improvements ([@ashmatadeen](https://github.com/ashmatadeen), [@loliee](https://github.com/loliee), [@cerisier](https://github.com/cerisier))
* Support of new SSH configuration fields (`AskPassGUI`, `GSSAPIClientIdentity`, `GSSAPIKeyExchange`, `GSSAPIRenewalForcesRekey`, `GSSAPIServerIdentity`, `GSSAPITrustDns`, `KeychainIntegration`)

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.2.0...v2.3.0)

### v2.2.0 (2016-02-03)

* Avoid exiting when an included file contains errors ([#95](https://github.com/moul/advanced-ssh-config/issues/95))
* Anonymize paths in `assh info`
* Support of `assh proxy --dry-run` option
* Fix: do not resolve variables in hostnames twice ([#103](https://github.com/moul/advanced-ssh-config/issues/103))

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.1.0...v2.2.0)

### v2.1.0 (2015-10-05)

* Expand environment variables ([#86](https://github.com/moul/advanced-ssh-config/issues/86))
* Add homebrew support ([#73](https://github.com/moul/advanced-ssh-config/issues/73))
* Add a 'ssh info' command ([#71](https://github.com/moul/advanced-ssh-config/issues/71))
* Templates support ([#52](https://github.com/moul/advanced-ssh-config/issues/52))
* Configuration is now case insensitive ([#51](https://github.com/moul/advanced-ssh-config/issues/51))
* Fix: resolving host fields for gateways ([#79](https://github.com/moul/advanced-ssh-config/issues/79))
* Fix: inheritance was not working for non assh-related fields ([#54](https://github.com/moul/advanced-ssh-config/issues/54))
* Fix: expanding variables in HostName ([#56](https://github.com/moul/advanced-ssh-config/issues/56))

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.0.0...v2.1.0)

### v2.0.0 (2015-09-07)

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

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/be4fea1632b1e9f8aa60585187338777baaf1210...v2.0.0)

### [v1](https://github.com/moul/advanced-ssh-config/tree/v1.1.0) (2015-07-22)

* Last Python version

### [POC](https://github.com/moul/advanced-ssh-config/commit/550f86c225d30292728ad24bc883b6d3a3e3f1b1) (2010-08-26)

* First Python version (POC)

## Docker

Experimental: `assh` may run in Docker, however you will have limitations:

* The `assh` containers does not have any binaries except `assh`, you can't use `ProxyCommand`, `ResolveCommand`...
* Docker may run on another host, `ssh localhost` will ssh to Docker host

```console
docker run -it --rm -v ~/.ssh:/.ssh moul/assh --help
```

`assh` in Docker is slower and has more limitations, but it may be useful for testing or if you plan to use a Docker host as a remote Gateway

## Alternative version

* [v1](https://github.com/moul/advanced-ssh-config/tree/v1) (2009-2015) - The original implementation. It worked quite well, but was a lot slower, less portable, harder to install for the user and harder to work on to develop new features and fix bugs

## Troubleshooting

#### I can't use gateways

`assh` uses the [built-in netcat mode of OpenSSH (shipped with OpenSSH 5.4)](https://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#Passing_through_a_gateway_using_netcat_mode) by default.
If your ssh client doesn't support this feature, you can configure a custom `ProxyCommand` configuration, i.e.,

```yaml
hosts:
  myserver:
    host: 1.2.3.4
    gateways:
    - mygateway
    # configure a custom proxycommand
    proxycommand: /bin/nc %h %p

  mygateway:
    host: 5.6.7.8
```

---

You can configure this rule globally:

```yaml
defaults:
  proxycommand: nc %h %p
```

---

Also, be sure to have netcat installed on your system, or use an alternative proxy binary, i.e., `socat`.

#### How to Configure resolver to parse `/etc/hosts` and/or handle **mDNS** requests?

**assh** resolves hostnames using the system built-in resolver, depending on the OS, you can enable new features and/or change modules order.

* [Linux - nsswitch documentation](http://man7.org/linux/man-pages/man5/nsswitch.conf.5.html)
* [Linux - mDNS support (nss-mdns)](http://0pointer.de/lennart/projects/nss-mdns/)
* [Mac OS X - `/etc/resolv.conf` documentation](https://developer.apple.com/library/mac/documentation/Darwin/Reference/ManPages/man5/resolver.5.html)

#### `unix_listener: "/Users/.../.ssh/cm/..." too long for Unix domain socket`

Starting with `OpenSSH v6.7` the socket name can be shortened by configuring `%C` for the name expansion.

```yaml
defaults:
  ControlPath: ~/tmp/.ssh/cm/%C.sock
```

`%C` is a unique identifier based on a hash of the tuple of (local host, remote user, hostname, port).

## License

© 2009-2016 Manfred Touron - MIT License

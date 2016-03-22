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
  * [Under the hood features](#under-the-hood-features)
3. [Configuration](#configuration)
4. [Usage](#usage)
5. [Install](#install)
6. [Getting started](#getting-started)
7. [Changelog](#changelog)
8. [Alternative version](#alternative-version)
9. [License](#license)

## Overview

A *transparent wrapper* that adds **regex**, **aliases**, **gateways**, **includes**, **dynamic hostnames** to **SSH**.

**Advanced SSH config** is wrapped in [lib-ssh](https://www.libssh.org) as a [ProxyCommand](https://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#ProxyCommand_with_Netcat), it means that it works seamlessly with:

* [ssh](http://www.openbsd.org/cgi-bin/man.cgi?query=ssh&sektion=1)
* [scp](http://www.openbsd.org/cgi-bin/man.cgi?query=scp&sektion=1)
* [rsync](http://linuxcommand.org/man_pages/rsync1.html)
* [git](https://www.kernel.org/pub/software/scm/git/docs/)
* Desktop applications depending on `lib-ssh` or `ssh` (i.e: [Tower](http://www.git-tower.com), [Atom.io](https://atom.io), [SSH Tunnel Manager](http://projects.tynsoe.org/fr/stm/))

## Features

### Configuration features

* **regex** support
* **aliases** `gate` -> `gate.domain.tld`
* **gateways** -> transparent ssh connection chaining
* **includes**: split configuration in multiple files
* **local command execution**: finally the reverse of **RemoteCommand**
* **templates**: equivalent to host but you can't connect directly to a template, perfect for inheritance
* **inheritance**: make hosts inherits from host hosts or templates
* **variable expansion**: resolve variables from environment
* **smart proxycommand**: RAW tcp connection when possible with `netcat` and `socat` as default fallbacks

### Using Gateway from command line

Connect to `hosta` using `hostb` as gateway.

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

### Under the hood features

* Automatically regenerates `~/.ssh/config` file when needed
* Inspect parent process to determine log level (if you use `ssh -vv`, **assh** will automatically be ran in debug mode)
* Automatically creates `ControlPath` directories so you can use *slashes* in your `ControlPath` option, can be disabled with the `NoControlMasterMkdir: true` configuration in host or globally.

## Configuration

The `~/.ssh/config` file is now managed by `assh`, take care to keep a backup your `~/.ssh/config` file.

`~/.ssh/assh.yml` is a [YAML](http://www.yaml.org/spec/1.2/spec.html) file containing:

* an `hosts` dictionary containing multiple *HOST* definitions
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
    - sad-robot
    - bighead

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
  # Templates are similar to Hosts, you can inherits from them
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
   2.2.0 (HEAD)

AUTHOR(S):
   Manfred Touron <https://github.com/moul/advanced-ssh-config>

COMMANDS:
   proxy         Connect to host SSH socket, used by ProxyCommand
   build         Build .ssh/config
   info          Display system-wide information
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
  --debug, -D       Enable debug mode [$ASSH_DEBUG]
  --verbose, -V     Enable verbose mode
  --help, -h        show help
  --version, -v     print the version
  ```

## Install

Get the latest version using GO (recommended way):

```bash
go get -u github.com/moul/advanced-ssh-config/cmd/assh
```

---

Get the latest version using homebrew (Mac OS X):

```bash
brew install assh
```

---

Get a released version on: https://github.com/moul/advanced-ssh-config/releases

## Getting started

1. Backup your old `~/.ssh/config`: `cp ~/.ssh/config ~/.ssh/config.backup`
2. Create a new `~/.ssh/assh.yml` file
3. Run `assh build > ~/.ssh/config` to validate the syntax of your `~/.ssh/assh.yml` file and automatically build your `~/.ssh/config` file
4. You are ready!

## Changelog

### master (unreleased)

* Add build information in .ssh/config header ([#49](https://github.com/moul/advanced-ssh-config/issues/49))
* Initial `Aliases` support ([#133](https://github.com/moul/advanced-ssh-config/issues/133))
* Use args[0] as ProxyCommand ([#134](https://github.com/moul/advanced-ssh-config/issues/134))
* Add `NoControlMasterMkdir` option to disable automatic creation of directories for gateways ([#124](https://github.com/moul/advanced-ssh-config/issues/124))
* Fix: Allow `$(...)` syntax in the `ResolveCommand` function ([#117](https://github.com/moul/advanced-ssh-config/issues/117))
* Printing the error of a failing `ResolveCommand` ([#117](https://github.com/moul/advanced-ssh-config/issues/117))
* Fix: `Gateways` field is no longer ignored when the `HostName` field is present ([#102](https://github.com/moul/advanced-ssh-config/issues/102))
* Ignore SIGHUP, close goroutines and export written bytes ([#112](https://github.com/moul/advanced-ssh-config/pull/112)) ([@QuentinPerez](https://github.com/QuentinPerez))
* Various documentation improvements ([@ashmatadeen](https://github.com/ashmatadeen), [@loliee](https://github.com/loliee), [@cerisier](https://github.com/cerisier))
* Support of new SSH configuration fields (`AskPassGUI`, `GSSAPIClientIdentity`, `GSSAPIKeyExchange`, `GSSAPIRenewalForcesRekey`, `GSSAPIServerIdentity`, `GSSAPITrustDns`, `KeychainIntegration`)

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.2.0...master)

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

## License

© 2009-2016 Manfred Touron - MIT License

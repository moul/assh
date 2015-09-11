# Advanced SSH config
[![Travis](https://img.shields.io/travis/moul/advanced-ssh-config.svg)](https://travis-ci.org/moul/advanced-ssh-config)
[![GoDoc](https://godoc.org/github.com/moul/advanced-ssh-config?status.svg)](https://godoc.org/github.com/moul/advanced-ssh-config)
![License](https://img.shields.io/github/license/moul/advanced-ssh-config.svg)


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

**Advanced SSH config** is wrapped in [lib-ssh](https://www.libssh.org) as a [ProxyCommand](http://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#ProxyCommand_with_Netcat), it means that it works seamlessly with:

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
* **inheritance**: make hosts inherits from host templates
* **variable expansion**: resolve variables from environment
* **smart proxycommand**: RAW tcp connection when possible with `netcat` and `socat` as default fallbacks

### Using Gateway from command line

Connect to `hosta` using `hostb` as gateway.

```console
$ ssh hosta/hostb
user@hosta $
```

Equivalent to `ssh -o ProxyCommand="ssh hostb nc %h %p" hosta`

---

Connect to `host` using `hostb` as a gateway using `hostc` as a gateway.

```console
$ ssh hosta/hostb/hostc
user@hosta $
```

Equivalent to `ssh -o ProxyCommand="ssh -o ProxyCommand='ssh hostc nc %h %p' hostb nc %h %p" hosta`

### Under the hood features

* Automatically regenerates `~/.ssh/config` file when needed
* Inspect parent process to determine log level (if you use `ssh -vv`, **assh** will automatically be ran in debug mode)
* Automatically creates `ControlPath` directories so you can use *slashes* in your `ControlPath` option

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
    HostName: 1.2.3.4
    User: robert
    Port: 2222

  bart:
    # ssh bart ->   ssh 5.6.7.8 -u bart           <- direct access
    #            or ssh 5.6.7.8/homer -u bart     <- using homer as a gateway
    HostName: 5.6.7.8
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

  "somehost[0-7]*":
    # ssh somehost2042 ->       ssh somehost2042.some.zone
    HostName: "%h.some.zone"

  vm-*.school.com:
    # ssh vm-42.school.com ->   ssh vm-42.school.com/gw.school.com -l student -o ForwardX11=yes -i ~/.ssh/school-rsa
    Gateways:
    - schoolgw
    Inherits:
    - schooltemplate

  "*.scw":
    # ssh toto.scw -> 1. dynamically resolves the IP address
    #                 2. ssh {resolved ip address} -u root -p 22 -o UserKnownHostsFile=null -o StrictHostKeyChecking=no
    # requires github.com/scaleway/scaleway-cli
    ResolveCommand: /bin/sh -c "scw inspect -f {{.PublicAddress.IP}} server:$(echo %h | sed s/.scw//)"
    User: root
    Port: 22
    UserKnownHostsFile: /dev/null
    StrictHostKeyChecking: no

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
   2.0.0 (HEAD)

AUTHOR(S):
   Manfred Touron <https://github.com/moul/advanced-ssh-config>

COMMANDS:
   proxy         Open an SSH connection to HOST
   stats         Print statistics
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
  --debug, -D       Enable debug mode
  --verbose, -V     Enable verbose mode
  --help, -h        show help
  --version, -v     print the version
  ```

## Install

Get the latest version using:

```bash
go get -u github.com/moul/advanced-ssh-config/...
```

Get a released version on: https://github.com/moul/advanced-ssh-config/releases

## Getting started

1. Backup your old `~/.ssh/config`: `cp ~/.ssh/config ~/.ssh/config.backup`
2. Create a new `~/.ssh/assh.yml` file
3. Run `assh proxy localhost` to validates the syntax of your `~/.ssh/assh.yml` file and automatically build your `~/.ssh/config` file
4. You are ready!

## Changelog

### master (unreleased)

* Fix: inheritance was not working for non assh-related fields ([#54](https://github.com/moul/advanced-ssh-config/issues/54))
* Fix: expanding variables in HostName ([#56](https://github.com/moul/advanced-ssh-config/issues/56))

[Full commits list](https://github.com/moul/advanced-ssh-config/compare/v2.0.0...master)

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

## Alternative version

* [v1](https://github.com/moul/advanced-ssh-config/tree/v1) (2009-2015) - The original implementation. It worked quite well, but was a lot slower, less portable, harder to install for the user and harder to work on to develop new features and fix bugs

## License

© 2009-2015 Manfred Touron - MIT License


[![ASSH logo - Advanced SSH Config logo](https://raw.githubusercontent.com/moul/advanced-ssh-config/develop/assets/assh.jpg)](https://github.com/moul/advanced-ssh-config)

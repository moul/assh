Advanced SSH config
===================

[![PyPI version](https://badge.fury.io/py/advanced-ssh-config.png)](http://badge.fury.io/py/advanced-ssh-config)
[![Build Status](https://travis-ci.org/moul/advanced-ssh-config.png?branch=develop)](https://travis-ci.org/moul/advanced-ssh-config)
[![authors](https://sourcegraph.com/api/repos/github.com/moul/advanced-ssh-config/badges/authors.png)](https://sourcegraph.com/github.com/moul/advanced-ssh-config)
[![library users](https://sourcegraph.com/api/repos/github.com/moul/advanced-ssh-config/badges/library-users.png)](https://sourcegraph.com/github.com/moul/advanced-ssh-config)
[![Total views](https://sourcegraph.com/api/repos/github.com/moul/advanced-ssh-config/counters/views.png)](https://sourcegraph.com/github.com/moul/advanced-ssh-config)
[![Views in the last 24 hours](https://sourcegraph.com/api/repos/github.com/moul/advanced-ssh-config/counters/views-24h.png)](https://sourcegraph.com/github.com/moul/advanced-ssh-config)
[![PyPi downloads](https://pypip.in/d/advanced-ssh-config/badge.png)](https://crate.io/packages/advanced-ssh-config/)
[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/moul/advanced-ssh-config/trend.png)](https://bitdeli.com/free "Bitdeli Badge")
[![Gitter chat](https://badges.gitter.im/moul/advanced-ssh-config.png)](https://gitter.im/moul/advanced-ssh-config)


Enhances `ssh_config` file capabilities

**NOTE**: This program is called by [ProxyCommand](http://en.wikibooks.org/wiki/OpenSSH/Cookbook/Proxies_and_Jump_Hosts#ProxyCommand_with_Netcat) from [lib-ssh](https://www.libssh.org).

---

It works *transparently* with :

- ssh
- scp
- rsync
- git
- and even desktop applications depending on `lib-ssh` (for instance [Tower](http://www.git-tower.com), [Atom.io](https://atom.io), [SSH Tunnel Manager](http://projects.tynsoe.org/fr/stm/))

---

The `.ssh/config` file is automatically generated, you need to update
`.ssh/config.advanced` file instead;
With new features and a better regex engine for the hostnames.

Commmand line features
----------------------

**Gateway chaining**

    ssh foo.com/bar.com

Connect to `bar.com` using ssh and create a proxy on `bar.com` to `foo.com`. Then connect to `foo.com` using the created proxy on `bar.com`.

    ssh foo.com/bar.com/baz.com

Connect to `foo.com` using `bar.com/baz.com` which itself uses `baz.com`.

Configuration features
----------------------

- **regex for hostnames**: `gw.school-.*.domain.net`
- **aliases**: `gate` -> `gate.domain.tld`
- **gateways**: transparent ssh connections chaining
- **includes**: split configuration into multiple files
- **local command execution**: finally a way to execute a command locally on connection
- **inheritance**: `inherits = gate.domain.tld`
- **variable expansion**: `User = $USER` (take $USER from environment)
- **smart proxycommand**: connect using `netcat`, `socat` or custom handler

Config example
--------------

`~/.ssh/config.advanced`

    # Simple example
    [foo.com]
    user = pacman
    port = 2222

    [bar]
    hostname = 1.2.3.4
    gateways = foo.com   # `ssh bar` will use `foo.com` as gateway

    [^vm-[0-9]*\.joe\.com$]
    gateways = bar       # `ssh vm-42.joe.com will use `bar` as gateway which
                         # itself will use `foo.com` as gateway

    [default]
    ProxyCommand = assh --port=%p connect %h

---

    # Complete example
    [foo]
    user = pacman
    port = 2222
    hostname = foo.com

    [bar]
    hostname = 1.2.3.4
    gateways = foo
    # By running `ssh bar`, you will ssh to `bar` through a `ssh foo`

    [^vm-[0-9]*\.joe\.com$]
    IdentityFile = ~/.ssh/root-joe
    gateways = direct joe.com joe.com/bar
    # Will try to ssh without proxy, then fallback to joe.com proxy, then
    # fallback to joe.com through bar
    DynamicForward = 43217
    LocalForward = 1723 localhost:1723
    ForwardX11 = yes

    [default]
    Includes = ~/.ssh/config.advanced2 ~/.ssh/config.advanced3
    Port = 22
    User = root
    IdentityFile = ~/.ssh/id_rsa
    ProxyCommand = assh connect %h --port=%p
    Gateways = direct
    PubkeyAuthentication = yes
    VisualHostKey = yes
    ControlMaster = auto
    ControlPath = ~/.ssh/controlmaster/%h-%p-%r.sock
    EscapeChar = ~

Installation
------------

From Pypi

    # pip install advanced-ssh-config

Download the latest build

    # wget https://github.com/moul/advanced-ssh-config/releases/download/v0.16.0/advanced-ssh-config-Darwin-x86_64 -O /usr/local/bin/advanced-ssh-config && chmod +x /usr/local/bin/advanced-ssh-config
    # wget https://github.com/moul/advanced-ssh-config/releases/download/v0.16.0/advanced-ssh-config-$(uname -s)-$(uname -p) -O /usr/local/bin/advanced-ssh-config && chmod +x /usr/local/bin/advanced-ssh-config

Or by cloning

    # git clone https://github.com/moul/advanced-ssh-config
    # cd advanced-ssh-config
    # make install

First run
---------

Automatically generate a new `.ssh/config.advanced` based on your
current `.ssh/config` file:

    # advanced-ssh-config init

Tests
-----

    # make test

Docker
------

Build

    # docker build -t moul/advanced-ssh-config .

Run

    # docker run -rm -i -t moul/advanced-ssh-config
    or
    # docker run -rm -i -t -v $(pwd)/:/advanced_ssh_config moul/advanced-ssh-config
    or
    # docker run -rm -i -t -v moul/advanced-ssh-config python setup.py test

Contributors
------------

- [Christo DeLange](https://github.com/dldinternet)

--

© 2009-2015 Manfred Touron - [MIT License](https://github.com/moul/advanced-ssh-config/blob/master/License.txt).

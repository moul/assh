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


Enhances `ssh_config` file capabilities

This program is called by `ProxyCommand` from `lib-ssh`
It works with `ssh`, `scp`, `rsync`, `git`, etc

The new `.ssh/config` file become `.ssh/config.advanced` with new features and a better regex engine for the hostnames.
Each time the script is called, it recreate a whole new `.ssh/config`, so be careful, backup your old .ssh/config file !

[![Gitter chat](https://badges.gitter.im/moul/advanced-ssh-config.png)](https://gitter.im/moul/advanced-ssh-config)

Features
--------

- regex for hostnames - gw.school-*.*.domain.net
- aliases
- gateways - transparent ssh connections chaining
- includes files
- real local command - executes a command on the local shell when connecting
- intelligent proxycommand with fallbacks
- inherits configuration
- variable expansion

Config example
--------------

`~/.ssh/config.advanced`

    # Simple example
    [foo.com]
    user = pacman
    port = 2222
    
    [bar]
    hostname = 1.2.3.4
    gateways = foo.com
    
    [default]
    ProxyCommand = /path/to/advanced-ssh-config.py --hostname=%h --port=%p -u

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
    
    [vm-.*\.joe\.com]
    IdentityFile = ~/.ssh/root-joe
    gateways = direct joe.com joe.com/bar
    DynamicForward = 43217
    LocalForward = 1723 localhost:1723
    ForwardX11 = yes

    [default]
    Includes = ~/.ssh/config.advanced2 ~/.ssh/config.advanced3
    Port = 22
    User = root
    IdentityFile = ~/.ssh/id_rsa
    ProxyCommand = ~/.ssh/scripts/dynamic_ssh.py --hostname=%h --port=%p -u
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

Or by cloning

    # git clone https://github.com/moul/advanced-ssh-config
    # cd advanced-ssh-config
    # python setup.py install

Backup your old ~/.ssh/config file

    # cp ~/.ssh/config ~/.ssh/config.backup

Generate the new config file

    # advanced-ssh-config -u

Or add this line manually in your ~/.ssh/config file

    ...
    ProxyCommand = advanced-ssh-config --hostname=%h --port=%p -u
    ...

Tests
-----

Install test dependencies and run tests

    # python setup.py test
    
Pep8

    # pep8 advanced_ssh_config | grep -v /tests/

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

© 2009-2014 Manfred Touron - [MIT License](https://github.com/moul/advanced-ssh-config/blob/master/License.txt).

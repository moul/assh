Advanced SSH config
===================

Enhances `ssh_config` file capabilities

This program is called by `ProxyCommand` from `lib-ssh`
It works with `ssh`, `scp`, `rsync`, `git`, etc

The new `.ssh/config` file become `.ssh/config.advanced` with new features and a better regex engine for the hostnames.
Each time the script is called, it recreate a whole new `.ssh/config`, so be careful, backup your old .ssh/config file !

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

- install test dependencies and run tests

    # python setup.py test

Docker
------

- build

    # docker build -t moul/advanced-ssh-config .

- run

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


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/moul/advanced-ssh-config/trend.png)](https://bitdeli.com/free "Bitdeli Badge")


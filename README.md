Advanced SSH config
===================

Enhance `ssh_config` file capabilities

This program is called by `ProxyCommand` from `lib-ssh`
It works with `ssh`, `scp`, `rsync`, `git`, etc

The new `.ssh/config` file become `.ssh/config.advanced` with new features and a better regex engine for the hostnames.
Each time the script is called, it recreate a whole new `.ssh/config`, so be careful, backup your old .ssh/config file !

Features
========

- regex for hostnames (gw.school-*.*.domain.net)
- aliases
- gateways (chains your ssh connections)
- includes (includes sub files)
- real local command (executes a command on the local shell)

Contributors
============

- [Christo DeLange](https://github.com/dldinternet)

Installation
============

From Pypi

    pip install advanced-ssh-config

Or by cloning

    git clone https://github.com/moul/advanced-ssh-config
    cd advanced-ssh-config
    python setup.py install

Backup your old ~/.ssh/config file

    cp ~/.ssh/config ~/.ssh/config.backup

Add this line in your ~/.ssh/config file

    ProxyCommand = advanced-ssh-config --hostname=%h --port=%p -u

© 2009-2013 Manfred Touron - [MIT License](https://github.com/moul/advanced-ssh-config/blob/master/License.txt).

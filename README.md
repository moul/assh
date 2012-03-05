Advanced SSH config
===================

by Manfred Touron

Using ssh_config ProxyCommand, ssh calls advanced-ssh-config.py.

The new .ssh/config file become .ssh/config.advanced with new features and a better regex engine for the hostnames.
Each time the script is called, it recreate a whole new .ssh/config, so be careful, backup your old .ssh/config file !

Features
========

* better regex engine (gw.school-*.*.domain.net)
* aliases
* gateways (chain your ssh connections)
* includes (include sub files)
* real local command (execute a command on the local shell)

Installation
============

Backup your old ~/.ssh/config file

`cp ~/.ssh/config ~/.ssh/config.backup`

Add this line in your ~/.ssh/config file

`ProxyCommand = /path/to/advanced-ssh-config.py --hostname=%h --port=%p -u`


© 2009-2012 Manfred Touron - [MIT License](https://github.com/moul/advanced-ssh-config/blob/master/License.txt).
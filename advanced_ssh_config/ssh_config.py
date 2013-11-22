# -*- coding: utf-8 -*-

import re

proxy_re = re.compile(r"^(proxycommand)\s*=*\s*(.*)", re.I)


def parse_ssh_config(file_obj):
    """
    Read an OpenSSH config from the given file object.

    Small adaptation of the paramiko.config.SSH_Config.parse method
    https://github.com/paramiko/paramiko/blob/master/paramiko/config.py

    @param file_obj: a file-like object to read the config file from
    @type file_obj: file
    """
    hosts = {}
    host = {"host": ['*'], "config": {}}
    for line in file_obj:
        line = line.rstrip('\n').lstrip()
        if (line == '') or (line[0] == '#'):
            continue
        if '=' in line:
            # Ensure ProxyCommand gets properly split
            if line.lower().strip().startswith('proxycommand'):
                match = proxy_re.match(line)
                key, value = match.group(1).lower(), match.group(2)
            else:
                key, value = line.split('=', 1)
                key = key.strip().lower()
        else:
            # find first whitespace, and split there
            i = 0
            while (i < len(line)) and not line[i].isspace():
                i += 1
            if i == len(line):
                raise Exception('Unparsable line: %r' % line)
            key = line[:i].lower()
            value = line[i:].lstrip()

        if key == 'host':
            hosts[host['host'][0]] = host['config']
            value = value.split()
            host = {key: value, 'config': {}}
        #identityfile, localforward, remoteforward keys are special cases, since they are allowed to be
        # specified multiple times and they should be tried in order
        # of specification.

        elif key in ['identityfile', 'localforward', 'remoteforward']:
            if key in host['config']:
                host['config'][key].append(value)
            else:
                host['config'][key] = [value]
        elif key not in host['config']:
            host['config'].update({key: value})
    hosts[host['host'][0]] = host['config']
    return hosts

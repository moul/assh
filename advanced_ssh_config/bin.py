# -*- coding: utf-8 -*-

import sys
import optparse
import logging
import re
import os

from . import __version__
from .utils import (
    validate_host, validate_port, parent_ssh_process_info, setup_logging
    )
from .config import Config
from .exceptions import ConfigError
from .advanced_ssh_config import AdvancedSshConfig
from .ssh_config import parse_ssh_config


def advanced_ssh_config_parse_options():
    parser = optparse.OptionParser(usage='%prog [-v] [-l 9] -H host [-p 22]',
                                   version='%prog {0}'.format(__version__))

    parser.add_option('-H', '--hostname',
                      dest='hostname',
                      help='Host')

    parser.add_option('-p', '--port',
                      dest='port',
                      default=None)

    parser.add_option('--file',
                      dest='file',
                      default='~/.ssh/config',
                      help='ssh_config file')

    parser.add_option('-f', '--force',
                      dest='force',
                      default=False,
                      action='store_true',
                      help='force update if versions differ')

    parser.add_option('-v', '--verbose',
                      dest='verbose',
                      action='store_true')

    parser.add_option('-l', '--log_level',
                      dest='log_level')

    parser.add_option('-u', '--update-sshconfig',
                      dest='update_sshconfig',
                      action='store_true')

    parser.add_option('--dry-run',
                      action='store_true',
                      dest='dry_run')

    (options, args) = parser.parse_args()

    if len(args):
        raise ValueError('This program only takes options, not args')

    if options.update_sshconfig and not options.hostname:
        return options

    validate_host(options.hostname)
    if options.port is not None:
        validate_port(options.port)
        options.port = int(options.port)

    return options


def keyboard_interrupt(fn):
    """ KeyboardInterrupt interceptor decorator. """
    logger = logging.getLogger('assh.advanced_ssh_config')
    try:
        fn()
    except KeyboardInterrupt:
        logging.fatal('Advanced SSH interrupted, bye.')
        sys.exit(1)


@keyboard_interrupt
def advanced_ssh_config():
    """ advanced-ssh-config entry-point. """
    try:
        options = advanced_ssh_config_parse_options()
    except ValueError as err:
        logging.fatal(err.message)
        sys.exit(1)

    parent = parent_ssh_process_info()
    setup_logging(options, parent)
    logger = logging.getLogger('assh.advanced_ssh_config')

    try:
        ssh = AdvancedSshConfig(hostname=options.hostname,
                                port=options.port,
                                verbose=options.verbose,
                                dry_run=options.dry_run,
                                ssh_config_file=options.file,
                                force=options.force)
        if options.update_sshconfig:
            ssh.write_sshconfig()

        if ssh.hostname:
            routing = ssh.get_routing()
            ssh.connect(routing)
        elif not options.update_sshconfig:
            print 'Must specify a host!\n'

    except ConfigError as err:
        sys.stderr.write(err.message)

    except Exception as err:
        sys.stderr.write(err.__str__())


def ssh_config_to_advanced_ssh_config_parse_options():
    parser = optparse.OptionParser(usage='%prog',
                                   version='%prog {0}'.format(__version__))

    parser.add_option('-f', '--file',
                      dest='file',
                      default='~/.ssh/config',
                      help='ssh_config file to parse')

    parser.add_option('-a', '--all',
                      action='store_true',
                      help='include hosts without configuration')

    parser.add_option('--no-escape',
                      action='store_false',
                      default=True,
                      dest='escape',
                      help='escape host for regexp')

    (options, args) = parser.parse_args()

    options.file = os.path.expanduser(options.file)
    if not os.path.exists(options.file):
        raise ValueError('File not found: {0}'.format(options.file))
    return options


@keyboard_interrupt
def ssh_config_to_advanced_ssh_config():
    """ ssh-config-to-advanced-ssh-config entry-point. """
    try:
        options = ssh_config_to_advanced_ssh_config_parse_options()
    except ValueError as err:
        logger = logging.getLogger('assh')
        logger.fatal(err.message)
        sys.exit(1)

    with open(options.file, 'r') as file_descriptor:
        config = parse_ssh_config(file_descriptor)
        print(options.escape)
        for host, config in config.iteritems():
            if not config and not options.all:
                continue
            if options.escape:
                host = host.replace('.', '\\.')
                host = host.replace('*', '.*')
                host = '^{0}$'.format(host)
            print('[{0}]'.format(host))
            for key, value in config.iteritems():
                print('  {0} = {1}'.format(key, value))
            print('')


@keyboard_interrupt
def assh_to_etchosts():
    """ assh-to-etchosts entry-point. """
    print('')
    print('## Automatically generated with assh-to-etchosts')
    configfiles = [
        '/etc/ssh/config.advanced',
        '~/.ssh/config.advanced',
    ]
    config = Config(configfiles=configfiles)
    hosts = {}
    for sect in config.parser.sections():
        ip = config.get('hostname', sect)
        if ip:
            if '$' in ip:
                continue  # command
            # FIXME: handle IPV6
            if not re.match(r'[0-9\.]+', ip):
                # FIXME: try to resolve non-ip hostnames
                continue
            if ip not in hosts:
                hosts[ip] = []
            if sect not in hosts[ip]:
                hosts[ip].append(sect)
    for ip in sorted(hosts.keys()):
        hostnames = hosts[ip]
        print("{:40} {}".format(ip, ' '.join(hostnames)))

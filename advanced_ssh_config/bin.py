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


AVAILABLE_COMMANDS = (
    'build',
    'connect',
    'generate-etc-hosts',
    'help',
    'info',
    'init',
    'stats',
)


def advanced_ssh_config_parse_options():
    parser = optparse.OptionParser(
        usage="""%prog [OPTIONS] COMMAND [arg...]

Commands:
  build                 Build .ssh/config based on .ssh/config.advanced
  connect <host>        Open a connection to <host>
  info <host>           Print connection informations
  init                  Build a .ssh/config.advanced file based on .ssh/config
  generate-etc-hosts    Print a /etc/hosts file of .ssh/config.advanced
  stats                 Print statistics""",
        version='%prog {0}'.format(__version__)
    )
    parser.add_option('-p', '--port',
                      dest='port',
                      default=None,
                      help='SSH port')

    parser.add_option('-c', '--config',
                      dest='config_file',
                      default='~/.ssh/config',
                      help='ssh_config file')

    parser.add_option('-f', '--force',
                      dest='force',
                      default=False,
                      action='store_true')

    parser.add_option('-v', '--verbose',
                      dest='verbose',
                      action='store_true')

    parser.add_option('-l', '--log_level',
                      dest='log_level')

    parser.add_option('--dry-run',
                      action='store_true',
                      dest='dry_run')

    (options, args) = parser.parse_args()

    # Must specify a COMMAND
    if not len(args):
        parser.print_help()
        exit(-1)

    options.command = args[0]

    # COMMAND must exists
    if options.command not in AVAILABLE_COMMANDS:
        raise ValueError("'{}' is not a valid command."
                         .format(options.command))

    if options.command == 'help':
        parser.print_help()
        exit(-1)

    options.hostname = None
    # some COMMANDS needs a HOST as argument
    if options.command in ('connect', 'info'):
        if len(args) < 2:
            parser.print_help()
            exit(-1)
        options.hostname = args[1]
        validate_host(options.hostname)
        if options.port is not None:
            validate_port(options.port)
            options.port = int(options.port)

    return options


def keyboard_interrupt(fn):
    """ KeyboardInterrupt interceptor decorator. """
    def wrapper():
        logger = logging.getLogger('assh.advanced_ssh_config')
        try:
            fn()
        except KeyboardInterrupt:
            logging.fatal('Advanced SSH interrupted, bye.')
            sys.exit(1)
    return wrapper


def assh_from_options(options):
    return AdvancedSshConfig(
        hostname=options.hostname,
        port=options.port,
        verbose=options.verbose,
        dry_run=options.dry_run,
        ssh_config_file=options.config_file,
        force=options.force
    )


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
        if options.command == 'build':
            ssh = assh_from_options(options)
            ssh.write_sshconfig()
            print("The file '{}' has been rebuilt".format(options.config_file))
            # FIXME: print some stats
            # FIXME: print diff

        elif options.command == 'connect':
            ssh = assh_from_options(options)
            routing = ssh.get_routing()
            ssh.connect(routing)

        elif options.command == 'info':
            ssh = assh_from_options(options)
            routing = ssh.get_routing()
            for key, value in routing.items():
                if value is None:
                    continue
                if isinstance(value, list):
                    if len(value):
                        print('{}:'.format(key))
                        print('\n'.join(['- {}'.format(line) for line in value]))
                else:
                    print('{}: {}'.format(key, value))

        elif options.command == 'generate-etc-hosts':
            assh_to_etchosts()

        elif options.command == 'stats':
            ssh = assh_from_options(options)
            hosts = ssh.prepare_sshconfig()
            print('{} hosts'.format(len(hosts.keys())))


        elif options.command == 'init':
            ssh_config_file = os.path.expanduser('~/.ssh/config')
            with open(ssh_config_file, 'r') as file_descriptor:
                config = parse_ssh_config(file_descriptor)
                for host, config in config.iteritems():

                    # Escape regex
                    if host == '*':
                        escaped_host = 'default'
                    else:
                        escaped_host = host.replace('.', '\\.')
                        escaped_host = escaped_host.replace('*', '.*')
                        # escaped_host = '^{0}$'.format(escaped_host)

                    print('[{0}]'.format(escaped_host))
                    print('# non-escaped host: {}'.format(host))
                    for key, value in config.iteritems():
                        if isinstance(value, list):
                            for entry in value:
                                print('{0} = {1}'.format(key, entry))
                            pass
                        else:
                            print('{0} = {1}'.format(key, value))
                    print('')


        else:
            raise NotImplementedError(
                "Command '{}' not yet implemented".format(options.command)
            )

    except ConfigError as err:
        sys.stderr.write(err.message)

    except Exception as err:
        sys.stderr.write(err.__str__().strip() + '\n')


def assh_to_etchosts():
    """ assh-to-etchosts entry-point. """
    print('## Automatically generated with assh-to-etchosts')
    print('')
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

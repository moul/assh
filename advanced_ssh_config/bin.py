# -*- coding: utf-8 -*-

import sys
import optparse
import logging

from .utils import LOGGING_LEVELS, validate_host, validate_port
from .exceptions import ConfigError
from .advanced_ssh_config import AdvancedSshConfig


def parse_options():
    parser = optparse.OptionParser(usage='%prog [-v] -h hostname -p port',
                                   version='%prog 1.0')

    parser.add_option('-H', '--hostname',
                      dest='hostname',
                      help='Host')

    parser.add_option('-p', '--port',
                      dest='port',
                      default=None)

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

    validate_host(options.hostname)
    if not options.port is None:
        validate_port(options.port)
        options.port = int(options.port)

    return options


def main():
    options = parse_options()

    # Setup logging
    logging_level = LOGGING_LEVELS.get(options.log_level, logging.ERROR)
    if options.verbose and logging_level == logging.ERROR:
        logging_level = logging.DEBUG
    logging.basicConfig(level=logging_level,
                        filename=None,
                        format='%(asctime)s %(levelname)s: %(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S')

    try:
        ssh = AdvancedSshConfig(hostname=options.hostname,
                                port=options.port,
                                verbose=options.verbose,
                                dry_run=options.dry_run)
        if options.update_sshconfig:
            ssh.update_sshconfig()

        if ssh.hostname:
            routing = ssh.get_routing()
            ssh.connect(routing)
        elif not options.update_sshconfig:
            print 'Must specify a host!\n'

    except KeyboardInterrupt:
        logging.error('Advanced SSH Interrupted, bye.')
        sys.exit(1)

    except ConfigError as err:
        sys.stderr.write(err.message)

    except Exception as err:
        sys.stderr.write(err.__str__())

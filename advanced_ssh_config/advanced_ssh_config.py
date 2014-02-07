# -*- coding: utf-8 -*-

import subprocess
import os
import sys
import re
import logging
from collections import OrderedDict

from .config import Config
from .utils import (safe_makedirs, value_interpolate, construct_proxy_commands,
                    shellquotemultiple)


class AdvancedSshConfig(object):

    def __init__(self, hostname=None, port=None, configfiles=None,
                 verbose=False, dry_run=False, proxy_type='nc',
                 timeout=180, use_python_socket=True):

        self.verbose, self.dry_run = verbose, dry_run
        self.hostname, self.port = hostname, port
        self.proxy_type, self.timeout = proxy_type, timeout
        self.user_python_socket = use_python_socket

        self.log = logging.getLogger('')

        # Initializes the Config object
        if not configfiles:
            configfiles = [
                '/etc/ssh/config.advanced',
                '~/.ssh/config.advanced',
                ]
        self.config = Config(configfiles=configfiles)

    @property
    def controlpath_dir(self):
        controlpath = self.config.get('controlpath',
                                      'default',
                                      '/tmp/advssh_cm/')
        directory = os.path.dirname(os.path.expanduser(controlpath))
        directory = os.path.join(directory, self.hostname)
        directory = os.path.dirname(directory)
        return directory

    def get_routing(self):
        routing = {}
        safe_makedirs(self.controlpath_dir)

        section = None
        for sect in self.config.parser.sections():
            if re.match(sect, self.hostname):
                section = sect

        logging.debug('section "{}" '.format(section))

        # Parse special routing
        path = self.hostname.split('/')

        args = {}
        options = {
            'p': 'Port',
            'l': 'User',
            'h': 'Hostname',
            'i': 'IdentityFile'
            }
        default_options = {
            'p': None,
            'h': path[0]
            }
        if self.port:
            default_options['p'] = self.port
        updated = False
        for key in options:
            cfval = self.config.get(options[key],
                                    path[0],
                                    default_options.get(key))
            value = value_interpolate(cfval)
            if cfval != value:
                updated = True
                self.config.parser.set(section, options[key], value)
                args[key] = value

            logging.debug('get (-%-1s) %-12s : %s' % (key,
                                                      options[key],
                                                      value))
            if value:
                args[key] = value

        # If we interpolated any keys
        if updated:
            self.write_sshconfig()
            self.log.debug('Config updated. Need to restart SSH!?')

        logging.debug('args: {}'.format(args))

        routing['verbose'] = self.verbose
        routing['proxy_type'] = self.proxy_type
        for special_key in ('comment', 'password', 'gateways',
                            'reallocalcommand'):
            routing[special_key] = self.config.get(special_key, path[-1], None)
        if not routing['gateways']:
            routing['gateways'] = ['direct']
        else:
            routing['gateways'] = routing['gateways'].split(' ')

        routing['gateway_route'] = path[1:]
        routing['hostname'] = args['h']
        #routing['args'] = args
        routing['port'] = self.port
        if not routing['port'] and 'p' in args:
            routing['port'] = int(args['p'])
        if not routing['port']:
            routing['port'] = 22
        routing['proxy_commands'] = construct_proxy_commands(routing)

        logging.debug('Routing:')
        for key, value in routing.iteritems():
            logging.debug('  {0}: {1}'.format(key, value))

        return routing

    def connect(self, routing):
        for gateway in routing['gateways']:
            if gateway != 'direct':
                routing['gateway_route'] += [gateway]
                logging.info('Using gateway: '
                             '{}'.format(routing['gateway_route']))
            else:
                logging.info('Direct connection')
            cmd = []
            if len(routing['gateway_route']):
                cmd += ['ssh', '/'.join(routing['gateway_route'])]
                cmd.append(shellquotemultiple(routing['proxy_commands']))
                logging.info('cmd: {}'.format(cmd))
            else:
                cmd = routing['proxy_commands'][0]

            logging.info('Connection command {}'.format(map(str, cmd)))

            if not self.dry_run:
                self.connect_once(routing, cmd)

    def connect_once(self, routing, cmd):
        comment = routing.get('comment', None)
        if comment:
            sys.stderr.write('{}\n'.format('\n'.join(comment)))

        password = routing.get('password', None)
        if password:
            sys.stderr.write('password: {}\n'.format(password))

        rlc_process = None
        if routing['reallocalcommand']:
            logging.info('Executing localcommand: '
                         '{}'.format(routing['reallocalcommand']))
            rlc_cmd = ['/bin/sh', '-c', routing['reallocalcommand']]
            rlc_process = subprocess.Popen(rlc_cmd,
                                           stdout=sys.stderr,
                                           stderr=sys.stderr)

        if self.user_python_socket \
                and not len(routing['gateway_route']):
            logging.info('Using Python socket')
            from .network import Socket
            socket = Socket(routing['hostname'], routing['port'])
            socket.run()
        else:
            logging.info('Using ProxyCommand')
            proxy_process = subprocess.Popen(map(str, cmd))
            if proxy_process.wait() != 0:
                self.log.critical('There were some errors')

        if rlc_process is not None:
            print(rlc_process)
            rlc_process.kill()

    def write_sshconfig(self, filename='~/.ssh/config'):
        config = self.build_sshconfig()
        fhandle = open(os.path.expanduser(filename), 'w+')
        fhandle.write('\n'.join(config))
        fhandle.close()

    def build_sshconfig(self):
        config = []

        hosts = self.prepare_sshconfig()
        od = OrderedDict(sorted(hosts.items()))
        for entry in od.values():
            if entry.host == '*':
                continue
            else:
                config += entry.build_sshconfig()

        if '*' in hosts:
            config += build_entry(hosts['*'])

        return config

    def prepare_sshconfig(self):
        hosts = {}
        for host in self.config.full.values():
            host.resolve()
            hosts[host.host] = host
        return hosts

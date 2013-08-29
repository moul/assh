#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys, os, ConfigParser, re, socket, subprocess, optparse, logging


LOGGING_LEVELS = {
    'crit': logging.CRITICAL,
    'critical': logging.CRITICAL,
    'err': logging.ERROR,
    'error': logging.ERROR,
    'warn': logging.WARNING,
    'warning': logging.WARNING,
    'info': logging.INFO,
    'debug': logging.DEBUG
    }


class ConfigError(Exception):
    '''Config exceptions.'''


class AdvancedSshConfig():
    def __init__(self, hostname=None, port=22, configfile=None, verbose=False, update_sshconfig=False):
        self.verbose, self.hostname, self.port = verbose, hostname, port

        self.log = logging.getLogger('')

        self.configfiles = ['/etc/ssh/config.advanced', os.path.expanduser('~/.ssh/config.advanced')]
        if configfile:
            self.configfiles += configfile
        self.parser = ConfigParser.ConfigParser()
        self.parser.SECTCRE = re.compile(
            r'\['
            r'(?P<header>.+)'
            r'\]'
            )

        errors = 0
        self.parser.read(self.configfiles)
        includes = self.conf_get('includes', 'default', '').strip()
        for include in includes.split():
            incpath = os.path.expanduser(include)
            if not incpath in self.configfiles and os.path.exists(incpath):
                self.parser.read(incpath)
            else:
                self.log.error('\'%s\' include not found' % incpath)
                errors += 1

        if 0 == errors:
            self.debug()
            self.debug('configfiles : %s' % self.configfiles)
            self.debug('================')
        else:
            raise ConfigError('Errors found in config')

        if update_sshconfig:
            self._update_sshconfig()

    def debug(self, str=None, force=False):
        self.log.debug(str and str or '')

    def conf_get(self, key, host, default=None, vardct=None):
        for section in self.parser.sections():
            if re.match(section, host):
                if self.parser.has_option(section, key):
                    return self.parser.get(section, key, False, vardct)
        if self.parser.has_option('default', key):
            return self.parser.get('default', key)
        return default

    def connect(self):
        # Handle special settings
        controlpath = self.conf_get('controlpath', 'default', '/tmp')
        mkdir_path = os.path.dirname(os.path.join(os.path.dirname(os.path.expanduser(controlpath)), self.hostname))
        try:
            os.makedirs(mkdir_path)
        except:
            pass

        section = None
        sectdct = None
        for sect in self.parser.sections():
            if re.match(sect, self.hostname):
                section = sect
                sectdct = self.parser.items(sect, True)

        #if not (section and sectdct):
        #    raise ConfigError(''%s' section not found!' % self.hostname)
        self.log.debug('section \'%s\' ' % section)

        # Parse special routing
        path = self.hostname.split('/')

        args = {}
        options = {
            'p': 'Port',
            'l': 'User',
            'h': 'Hostname',
            'i': 'IdentityFile'
            }
        matches = None
        updated = False
        for key in options:
            cfval = self.conf_get(options[key], path[0], False, {'hostname': self.hostname, 'port': self.port})
            value = self._interpolate(cfval)
            if cfval != value:
                updated = True
                self.parser.set(section, options[key], value)
                args[key] = value

            self.debug('get (-%-1s) %-12s : %s' % (key, options[key], value))
            if value:
                args[key] = value

        # If we interpolated any keys
        if updated:
            self._update_sshconfig()
            self.log.debug('Config updated. Need to restart SSH!?')

        if not 'h' in args:
            args['h'] = path[0]
        self.debug('args: %s' % args)
        self.debug()

        self.debug('hostname    : %s' % self.hostname)
        self.debug('port        : %s' % self.port)
        self.debug('path        : %s' % path)
        self.debug('path[0]     : %s' % path[0])
        self.debug('path[1:]    : %s' % path[1:])
        self.debug('args        : %s' % args)

        self.debug()
        gateways = self.conf_get('Gateways', path[-1], 'direct').strip().split(' ')
        reallocalcommand = self.conf_get('RealLocalCommand', path[-1], '').strip().split(' ')
        self.debug('reallocalcommand: %s' % reallocalcommand)
        for gateway in gateways:
            right_path = path[1:]
            if gateway != 'direct':
                right_path += [gateway]
            cmd = []
            if len(right_path):
                cmd += ['ssh', '/'.join(right_path)]

            cmd += ['nc', args['h'], args['p']]

            self.debug('cmd         : %s' % cmd)
            self.debug('================')
            self.debug()
            ssh_process = subprocess.Popen(cmd)
            reallocalcommand_process = None
            if len(reallocalcommand[0]):
                reallocalcommand_process = subprocess.Popen(reallocalcommand)
            if ssh_process.wait() != 0:
                self.log.critical('There were some errors')
            if reallocalcommand_process is not None:
                reallocalcommand_process.kill()

    def _update_sshconfig(self, write=True):
        config = []

        for section in self.parser.sections():
            if section != 'default':
                host = section
                host = re.sub('\.\*', '*', host)
                host = re.sub('\\\.', '.', host)
                config += ['Host %s' % host]
                for key, value in self.parser.items(section, False, {'Hostname': host}):
                    if key not in ('hostname', 'gateways', 'reallocalcommand', 'remotecommand'):
                        if key == 'alias':
                            key = 'hostname'
                        config += ['  %s %s' % (key, value)]
                config += ['']

        config += ['Host *']
        for key, value in self.parser.items('default'):
            if key not in ['hostname', 'gateways', 'includes']:
                config += ['  %s %s' % (key, value)]

        if write:
            file = open(os.path.expanduser('~/.ssh/config'), 'w+')
            file.write('\n'.join(config))
            file.close()
        else:
            print '\n'.join(config)

    def _interpolate(self, value):
        matches = value and re.match('\$(\w+)', value) or None
        if matches:
            var = matches.group(1)
            val = os.environ.get(var)
            if val:
                self.log.debug('\'%s\' => \'%s\'' % (value, val))
                return self._interpolate(re.sub('\$%s' % var, val, value))

        return value


def main():
    parser = optparse.OptionParser(usage='%prog [-v] -h hostname -p port', version='%prog 1.0')
    parser.add_option('-H', '--hostname', dest='hostname', help='Host')
    parser.add_option('-p', '--port', dest='port')
    parser.add_option('-v', '--verbose', dest='verbose', action='store_true')
    parser.add_option('-l', '--log_level', dest='log_level')
    parser.add_option('-u', '--update-sshconfig', dest='update_sshconfig', action='store_true')
    (options, args) = parser.parse_args()

    logging_level = LOGGING_LEVELS.get(options.log_level, logging.ERROR)
    if options.verbose and logging_level == logging.ERROR:
        logging_level = logging.DEBUG
    logging.basicConfig(level=logging_level,
                        filename=None,
                        format='%(asctime)s %(levelname)s: %(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S')
    log = logging.getLogger('')

    try:
        ssh = AdvancedSshConfig(hostname=options.hostname,
                                port=options.port,
                                verbose=options.verbose,
                                update_sshconfig=options.update_sshconfig)
        if ssh.hostname is None:
            print 'Must specify a host!\n'
        else:
            ssh.connect()
    except ConfigError as e:
        sys.stderr.write(e.message)
    except Exception as e:
        log.debug(e.__str__())

if __name__ == '__main__':
    main()

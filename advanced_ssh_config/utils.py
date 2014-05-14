# -*- coding: utf-8 -*-

import logging
import re
import os
import errno

from .exceptions import ConfigError


def shellquote_arg(before):
    before = str(before)
    after = before
    after = after.replace('\\', '\\\\')
    after = after.replace('\'', '\\\'')
    if after != before:
        return "'{}'".format(after)
    else:
        return before


def shellquote(cmd):
    if type(cmd) != list:
        raise ValueError('`cmd` must be a list')
    return ' '.join(map(shellquote_arg, cmd))


def shellquotemultiple(cmds):
    if type(cmds) != list:
        raise ValueError('`cmds` must be a list of list')
    for cmd in cmds:
        if type(cmd) != list:
            raise ValueError('`cmd` in `cmds` must be lists')
    if len(cmds) > 1:
        return '({})'.format(' 2>/dev/null || '.join(map(shellquote, cmds)))
    else:
        return cmds[0]


def construct_proxy_commands(config):
    cmds = []
    proxy_type = config.get('proxy_type', 'nc')
    verbose = config.get('verbose', False)
    timeout = config.get('timeout', 180)
    connection_timeout = config.get('connection_timeout', 5)
    if 'hostname' not in config or 'port' not in config:
        raise ValueError('hostname and port must be configured')
    hostname, port = config['hostname'], config['port']
    if proxy_type in ('nc', 'ncat', 'netcat'):
        cmd = []  # cmd with options
        cmd.append(proxy_type)
        if verbose:
            cmd.append('-v')
        if timeout:
            cmd.append('-w')
            cmd.append(timeout)
        if connection_timeout:
            cmd.append('-G')
            cmd.append(connection_timeout)
        cmd.append(hostname)
        cmd.append(port)
        cmds.append(cmd)
        cmd = []  # cmd without options
        cmd.append(proxy_type)
        cmd.append(hostname)
        cmd.append(port)
        cmds.append(cmd)
    elif proxy_type == 'socat':
        cmd = []
        cmd.append('socat')
        cmd.append('STDIN')
        cmd.append('TCP:{}:{}'.format(hostname, port))
        cmds.append(cmd)
    elif proxy_type == 'socat_http_proxy':
        cmd = []
        cmd.append('socat')
        cmd.append('STDIN')
        args = [
            config.get('http_proxy_host', '127.0.0.1'),
            hostname,
            port,
            config.get('http_proxy_port', 3128),
            ]
        if config.get('http_proxy_auth', None):
            args.append(config.get('http_proxy_auth'))
            cmd.append('PROXY:{}:{}:{},'
                       'proxyport={},'
                       'proxyauth={}'.format(*args))
        else:
            cmd.append('PROXY:{}:{}:{},proxyport={}'.format(*args))
        cmds.append(cmd)
    elif proxy_type == 'socat_socks':
        cmd = []
        cmd.append('socat')
        cmd.append('STDIN')
        args = [
            config.get('socks_host', '127.0.0.1'),
            hostname,
            port,
            config.get('socks_port', 1080),
            ]
        cmd.append('SOCKS:{}:{}:{},socksport={}'.format(*args))
        cmds.append(cmd)
    else:
        raise ValueError('proxy_type `{}` is not handled'.format(proxy_type))
    return cmds


def validate_host(host):
    if not type(host).__name__ in ('str', 'unicode'):
        raise ValueError('host must be a string')
    if len(host) == 0:
        raise ValueError('host is empty')


def validate_port(port):
    if type(port).__name__ == 'str':
        try:
            port = int(port)
        except ValueError:
            raise ValueError('port must be a number')
    if type(port).__name__ != 'int':
        raise ValueError('port must be an integer')
    if port < 1 or port > 65535:
        raise ValueError('port must be between 1-65535')


def safe_makedirs(directory):
    try:
        os.makedirs(directory)
    except OSError as exception:
        if exception.errno != errno.EEXIST:
            raise exception


def value_interpolate(value, already_interpolated=None):
    if not already_interpolated:
        already_interpolated = []
    if type(value) == int:
        return value
    matches = value and re.match(r'\$(\w+)', value) or None
    if matches:
        var = matches.group(1)
        if var in already_interpolated:
            raise ConfigError('Interpolation loop')
        val = os.environ.get(var)
        if val:
            logging.getLogger('').debug('\'{}\' => \'{}\''.format(value, val))
            new_value = re.sub(r'\${}'.format(var), val, value)
            return value_interpolate(new_value, already_interpolated + [var])

    return value


def parent_ssh_process_info():
    try:
        import psutil
        return psutil.Process(os.getppid())
    except:
        return None


def setup_logging(options, parent=None):

    def get_logging_level():
        if options.verbose:
            return logging.DEBUG

        if parent:
            parent_level = 0
            for arg in parent.cmdline:
                if arg == '-v':
                    parent_level += 1
                elif arg == '-vv':
                    parent_level += 2
                elif arg == '-vvv':
                    parent_level += 3
            if parent_level:
                parent_level = min(parent_level + 2, 4)
                return LOGGING_LEVELS.get(parent_level, logging.INFO)

        environ_log_level = os.environ.get('ASSH_LOG_LEVEL', None)
        if environ_log_level:
            return LOGGING_LEVELS.get(environ_log_level, logging.ERROR)

        return LOGGING_LEVELS.get(options.log_level, logging.ERROR)

    logging_level = get_logging_level()
    logging.basicConfig(level=logging_level,
                        filename=None,
                        format='%(asctime)s %(levelname)s: %(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S')


LOGGING_LEVELS = {
    'crit':     logging.CRITICAL,
    'critical': logging.CRITICAL,
    0:          logging.CRITICAL,
    'err':      logging.ERROR,
    'error':    logging.ERROR,
    1:          logging.ERROR,
    'warn':     logging.WARNING,
    'warning':  logging.WARNING,
    2:          logging.WARNING,
    'info':     logging.INFO,
    3:          logging.INFO,
    'debug':    logging.DEBUG,
    4:          logging.DEBUG,
    }

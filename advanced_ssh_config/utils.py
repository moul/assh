# -*- coding: utf-8 -*-

import logging
import re
import os
import errno

from .exceptions import ConfigError


def construct_proxy_command(config):
    cmd = []
    proxy_type = config.get('proxy_type', 'nc')
    verbose = config.get('verbose', False)
    timeout = config.get('timeout', 180)
    connection_timeout = config.get('connection_timeout', 5)
    if not 'hostname' in config or not 'port' in config:
        raise ValueError('hostname and port must be configured')
    hostname, port = config['hostname'], config['port']
    if proxy_type in ('nc', 'ncat', 'netcat'):
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
    elif proxy_type == 'socat':
        cmd.append('socat')
        cmd.append('STDIN')
        cmd.append('TCP:{}:{}'.format(hostname, port))
    elif proxy_type == 'socat_http_proxy':
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
    elif proxy_type == 'socat_socks':
        cmd.append('socat')
        cmd.append('STDIN')
        args = [
            config.get('socks_host', '127.0.0.1'),
            hostname,
            port,
            config.get('socks_port', 1080),
            ]
        cmd.append('SOCKS:{}:{}:{},socksport={}'.format(*args))
    else:
        raise ValueError('proxy_type `{}` is not handled'.format(proxy_type))
    return cmd


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


LOGGING_LEVELS = {
    'crit':     logging.CRITICAL,
    'critical': logging.CRITICAL,
    'err':      logging.ERROR,
    'error':    logging.ERROR,
    'warn':     logging.WARNING,
    'warning':  logging.WARNING,
    'info':     logging.INFO,
    'debug':    logging.DEBUG
    }

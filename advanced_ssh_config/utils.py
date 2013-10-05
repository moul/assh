# -*- coding: utf-8 -*-

import logging
import re
import os
import errno

from .exceptions import ConfigError


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


def safe_makedirs(dir):
    try:
        os.makedirs(dir)
    except OSError as exception:
        if exception.errno != errno.EEXIST:
            raise exception


def value_interpolate(value, already_interpolated=[]):
    matches = value and re.match(r'\$(\w+)', value) or None
    if matches:
        var = matches.group(1)
        if var in already_interpolated:
            raise ConfigError('Interpolation loop')
        val = os.environ.get(var)
        if val:
            logging.getLogger('').debug('\'{}\' => \'{}\''.format(value, val))
            return value_interpolate(re.sub(r'\${}'.format(var), val, value), already_interpolated + [var])

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

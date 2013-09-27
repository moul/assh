# -*- coding: utf-8 -*-

import logging


def validate_host(host):
    if not type(host).__name__ in ('str', 'unicode'):
        raise ValueError('host must be a string')
    if len(host) == 0:
        raise ValueError('host is empty')


def validate_port(port):
    if type(port).__name__ != 'int':
        raise ValueError('port must be an integer')
    if port < 1 or port > 65535:
        raise ValueError('port must be between 1-65535')


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

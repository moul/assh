# -*- coding: utf-8 -*-

import os
import logging
import ConfigParser
import re


class Config(object):
    def __init__(self, configfiles):

        self.configfiles = map(os.path.expanduser, configfiles)

        self.log = logging.getLogger('')

        self.parser = ConfigParser.ConfigParser()
        self.parser.SECTCRE = re.compile(
            r'\['
            r'(?P<header>.+)'
            r'\]'
            )

    def debug(self, string=None):
        self.log.debug(string and string or '')

    def read(self):
        errors = 0
        self.parser.read(self.configfiles)
        includes = self.get('includes', 'default', '').strip()
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

    def get(self, key, host, default=None, vardct=None):
        for section in self.parser.sections():
            if re.match(section, host):
                if self.parser.has_option(section, key):
                    return self.parser.get(section, key, False, vardct)
        if self.parser.has_option('default', key):
            return self.parser.get('default', key)
        return default

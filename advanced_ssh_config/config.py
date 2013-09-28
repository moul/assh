# -*- coding: utf-8 -*-

import os
import logging
import ConfigParser
import re

from .exceptions import ConfigError


class Config(object):
    def __init__(self, configfiles):

        self.configfiles = map(os.path.expanduser, configfiles)
        self.loaded_files = []

        self.log = logging.getLogger('')

        self.parser = ConfigParser.ConfigParser()
        self.parser.SECTCRE = re.compile(
            r'\['
            r'(?P<header>.+)'
            r'\]'
            )
        self._read()

        for section in self.sections:
            if re.sub(r'[^a-zA-Z0-9\\\.\*_-]', '', section) != section:
                raise ConfigError('Invalid characters used in section {}'.format(section))

    def debug(self, string=None):
        self.log.debug(string and string or '')

    def _load_file(self, filename):
        if filename in self.loaded_files:
            return
        self.parser.read(filename)
        self.loaded_files.append(filename)

    def _read(self):
        for configfile in self.configfiles:
            self._load_file(configfile)

        # Load sub files
        includes = self.get('includes', 'default', '').strip()
        for include in includes.split():
            incpath = os.path.expanduser(include)
            if os.path.exists(incpath):
                self._load_file(incpath)
            else:
                raise ConfigError('\'{}\' include not found'.format(incpath))

    @property
    def sections(self):
        return self.parser.sections()

    def get(self, key, host, default=None, vardct=None):
        for section in self.sections:
            if re.match(section, host):
                if self.parser.has_option(section, key):
                    return self.parser.get(section, key, False, vardct)
        if self.parser.has_option('default', key):
            return self.parser.get('default', key)
        return default

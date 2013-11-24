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
                raise ConfigError('Invalid characters used in '
                                  'section {}'.format(section))

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
        includes = str(self.get('includes', 'default', '')).strip()
        for include in includes.split():
            incpath = os.path.expanduser(include)
            if os.path.exists(incpath):
                self._load_file(incpath)
            else:
                raise ConfigError('\'{}\' include not found'.format(incpath))

    @property
    def sections(self):
        return self.parser.sections()

    def get_in_section(self, section, key, raw=False, vardct=None):
        if not self.parser.has_option(section, key):
            return False
        var = self.parser.get(section, key, raw, vardct)
        if key in ('identityfile', 'localforward', 'remoteforward'):
            var = var.split('\n')
            var = map(str.strip, var)
        return var

    def get(self, key, host, default=None, vardct=None):
        for section in self.sections:
            if re.match(section, host):
                val = self.get_in_section(section, key, vardct=vardct)
                if val:
                    return val
        val = self.get_in_section('default', key)
        return val or default

# -*- coding: utf-8 -*-

import os
import logging
import ConfigParser
import re
from collections import OrderedDict

from .exceptions import ConfigError


class ConfigHost(object):
    special_keys = (
        'hostname',
        'gateways',
        'reallocalcommand',
        'remotecommand',
        'includes',
        'inherits',
        #'proxycommand',
        )

    key_translation = {
        'alias': 'hostname',
        }

    def __init__(self, c, host, config=None, extra=None, inherited_config=None,
                 inherited_extra=None):
        self.c = c
        self.host = host
        self.config = config or {}
        self.extra = extra or {}
        self.inherited = None
        self.resolved = False

    @classmethod
    def prepare_hostname(cls, host):
        host = re.sub(r'\.\*', '*', host)
        host = re.sub(r'\\\.', '.', host)
        return host

    @classmethod
    def from_config_file(cls, c, host, entry):
        config = []
        extra_config = []
        for key, value in entry:
            if key in ConfigHost.key_translation:
                key = ConfigHost.key_translation.get(key)
            if key in ('identityfile', 'localforward', 'remoteforward'):
                values = value.split('\n')
                values = map(str.strip, values)
            else:
                values = [value]
            for line in values:
                if key in ConfigHost.special_keys:
                    extra_config.append((key, line))
                else:
                    config.append((key, line))
        return cls(c, host, config=config, extra=extra_config)

    def config_keys(self):
        return [entry[0] for entry in self.config]

    @property
    def config_dict(self):
        if not self.resolved:
            self.resolve()
        config = {}
        for entry in self.config:
            config[entry[0]] = entry[1]
        return config

    @property
    def extra_dict(self, sort=True):
        extra = {}
        for entry in self.extra:
            extra[entry[0]] = entry[1]
        if sort:
            return OrderedDict(sorted(extra.items()))
        return extra

    @property
    def clean_config(self, sort=True):
        config = self.config_dict
        if self.inherited:
            config = dict(self.inherited.items() + config.items())
        if sort:
            return OrderedDict(sorted(config.items()))
        return config

    def resolve(self, rec=10):
        if not rec:
            raise ConfigError('Maximum recursion deptch exceeded')
        if self.resolved:
            return
        for key, value in self.extra:
            if key == 'inherits':
                if value in self.c.full:
                    parent = self.c.full[value]
                    parent.resolve(rec - 1)
                    self.inherited = parent.clean_config
                else:
                    raise ConfigError('Inheriting an unkonwn '
                                      'host: `{}`'.format(value))
        self.resolved = True

    def get_prep_value(self):
        return {
            'config': self.config,
            'extra': self.extra,
            'inherited': self.inherited,
            }

    def __repr__(self):
        max_len = 50
        dict_string = ', '.join([
            '%s=%s' % (key, str(val)[:max_len - 2] + '..'
                       if len(str(val)) > max_len else str(val))
            for key, val in self.get_prep_value().items()
            ])
        return '<{}.{} at {} - {}>'.format(self.__class__.__module__,
                                           self.__class__.__name__,
                                           hex(id(self)),
                                           dict_string)

    def build_sshconfig(self):
        sub_config = []

        if self.host == 'default':
            host = '*'
        else:
            host = self.host

        sub_config.append('Host {}'.format(host))

        attrs = OrderedDict(sorted(self.clean_config.items()))
        for key, value in attrs.iteritems():
            sub_config.append('  {} {}'.format(key, value))

        for key, value in self.extra_dict.iteritems():
            sub_config.append('  # {} {}'.format(key, value))

        sub_config.append('')
        return sub_config


class Config(object):
    def __init__(self, configfiles):

        self.full_cache = None
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

    @property
    def full(self):
        if not self.full_cache:
            self.full_cache = {}
            for section in self.parser.sections():
                host = ConfigHost.prepare_hostname(section)
                config_file_entry = self.parser.items(section,
                                                      False,
                                                      {'Hostname': host})
                conf = ConfigHost.from_config_file(self, host,
                                                   config_file_entry)
                self.full_cache[host] = conf
        return self.full_cache

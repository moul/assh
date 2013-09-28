# -*- coding: utf-8 -*-

import unittest
import os

from advanced_ssh_config.config import Config
from advanced_ssh_config.exceptions import ConfigError


PREFIX = '/tmp/test-asc-config'
DEFAULT_CONFIG = os.path.join(PREFIX, 'config.advanced')


def write_config(contents, name='config.advanced'):
    with open(os.path.join(PREFIX, name), 'w') as f:
        f.write(contents)


class TestConfig(unittest.TestCase):

    def setUp(self):
        os.system('rm -rf {}'.format(PREFIX))
        os.makedirs(PREFIX)
        write_config('')

    def test_initialize_config(self):
        config = Config([DEFAULT_CONFIG])
        self.assertIsInstance(config, Config)

    def test_include_existing_files(self):
        write_config('', name='include-1')
        write_config('', name='include-2')
        contents = """
[default]
Includes = {0}/include-1 {0}/include-2
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        self.assertEquals(config.loaded_files, [
                DEFAULT_CONFIG,
                '{}/include-1'.format(PREFIX),
                '{}/include-2'.format(PREFIX),
                ])

    def test_include_not_exists(self):
        contents = """
[default]
Includes = {0}/include-1 {0}/include-2
""".format(PREFIX)
        write_config(contents)
        self.assertRaises(ConfigError, Config, [DEFAULT_CONFIG])

    def test_include_same_file(self):
        write_config('', name='include-1')
        contents = """
[default]
Includes = {0}/include-1 {0}/include-1
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        self.assertEquals(config.loaded_files, [
                DEFAULT_CONFIG,
                '{}/include-1'.format(PREFIX),
                ])

    def test_sections_simple(self):
        contents = """
[hosta]
[default]
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        self.assertEquals(config.sections, ['hosta', 'default'])

    def test_sections_with_double(self):
        contents = """
[hosta]
[hosta]
[default]
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        self.assertEquals(config.sections, ['hosta', 'default'])

    def test_sections_with_case(self):
        contents = """
[hosta]
[hostA]
[default]
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        self.assertEquals(config.sections, ['hosta', 'hostA', 'default'])

    def test_sections_with_regex(self):
        contents = """
[hosta]
[host.*]
[default]
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        self.assertEquals(config.sections, ['hosta', 'host.*', 'default'])

    def test_get_simple(self):
        contents = """
[hosta]
hostmame = 1.2.3.4
port = 23

[default]
port = 22
""".format(PREFIX)
        write_config(contents)
        config = Config([DEFAULT_CONFIG])
        print()
        print('-' * 80)
        print(config.get('hostname', 'hosta', 'def'))
        print('-' * 80)

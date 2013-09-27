# -*- coding: utf-8 -*-

import unittest
import sys

from advanced_ssh_config.bin import parse_options


def set_argv(*args):
    sys.argv = [''] + list(args)


class TestParseOptions(unittest.TestCase):

    def test_with_args_without_options(self):
        set_argv('localhost')
        self.assertRaises(ValueError, parse_options)

    def test_with_args_with_options(self):
        set_argv('--port=22', 'localhost')
        self.assertRaises(ValueError, parse_options)

        set_argv('localhost', '--port=22')
        self.assertRaises(ValueError, parse_options)

    def test_unkown_options(self):
        set_argv('--toto=titi')
        self.assertRaises(SystemExit, parse_options)

    def test_default(self):
        set_argv('--host=1.2.3.4')
        options = parse_options()
        should_be = {
            'log_level': None,
            'verbose': None,
            'dry_run': None,
            'hostname': '1.2.3.4',
            'update_sshconfig': None,
            'port': 22
            }
        self.assertEqual(options, should_be)

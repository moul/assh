# -*- coding: utf-8 -*-

import unittest
import sys
import os

from advanced_ssh_config.bin import advanced_ssh_config_parse_options


class RedirectStdStreams(object):
    def __init__(self, stdout=None, stderr=None):
        devnull = open(os.devnull, 'w')
        if stdout == 'devnull':
            stdout = devnull
        if stderr == 'devnull':
            stderr = devnull

        self._stdout = stdout or sys.stdout
        self._stderr = stderr or sys.stderr

    def __enter__(self):
        self.old_stdout, self.old_stderr = sys.stdout, sys.stderr
        self.old_stdout.flush(); self.old_stderr.flush()
        sys.stdout, sys.stderr = self._stdout, self._stderr

    def __exit__(self, exc_type, exc_value, traceback):
        self._stdout.flush(); self._stderr.flush()
        sys.stdout = self.old_stdout
        sys.stderr = self.old_stderr


def _advanced_ssh_config_parse_options(*args):
    sys.argv = [''] + list(map(str, args))
    return advanced_ssh_config_parse_options()


class TestParseOptions(unittest.TestCase):

    def test_with_args_without_options(self):
        self.assertRaises(ValueError, _advanced_ssh_config_parse_options, 'localhost')

    def test_with_args_with_options(self):
        self.assertRaises(ValueError, _advanced_ssh_config_parse_options, '--port=22', 'localhost')
        self.assertRaises(ValueError, _advanced_ssh_config_parse_options, 'localhost', '--port=22')

    def test_unkown_options(self):
        with RedirectStdStreams(stdout='devnull', stderr='devnull'):
            self.assertRaises(SystemExit, _advanced_ssh_config_parse_options, '--toto=titi')

    def test_default(self):
        options = _advanced_ssh_config_parse_options('--host=1.2.3.4')
        should_be = {
            'log_level': None,
            'verbose': None,
            'dry_run': None,
            'hostname': '1.2.3.4',
            'update_sshconfig': None,
            'port': None
            }
        self.assertEqual(options, should_be)

    def test_option_valid_host(self):
        for host in ('1.2.3.4',
                     'localhost',
                     'google.com',
                     'localhost/localhost',
                     '1.2.3.4.5',
                     '::1',
                     'fe80::1%lo0'):
            options = _advanced_ssh_config_parse_options('--host={}'.format(host))
            self.assertEqual(options.hostname, host)

            options = _advanced_ssh_config_parse_options('-H', host)
            self.assertEqual(options.hostname, host)

    def test_option_invalid_host(self):
        with RedirectStdStreams(stdout='devnull', stderr='devnull'):
            self.assertRaises(SystemExit, _advanced_ssh_config_parse_options, '-H')
            self.assertRaises(SystemExit, _advanced_ssh_config_parse_options, '--host')

        for host in ('',):
            self.assertRaises(ValueError, _advanced_ssh_config_parse_options, '--host={}'.format(host))
            self.assertRaises(ValueError, _advanced_ssh_config_parse_options, '-H', host)

    def test_option_valid_port(self):
        for port in (22, 1, 65535):
            options = _advanced_ssh_config_parse_options('--host=localhost', '--port={}'.format(port))
            self.assertEqual(options.port, port)

            options = _advanced_ssh_config_parse_options('--host=localhost', '-p', port)
            self.assertEqual(options.port, port)

    def test_option_no_port(self):
        options = _advanced_ssh_config_parse_options('--host=localhost')
        self.assertEqual(options.port, None)

    def test_option_invalid_port(self):
        with RedirectStdStreams(stdout='devnull', stderr='devnull'):
            self.assertRaises(SystemExit, _advanced_ssh_config_parse_options, '--host=localhost', '--port')
            self.assertRaises(SystemExit, _advanced_ssh_config_parse_options, '--host=localhost', '-p')

        for port in (-1, 65536, 'test', '', '1.0'):
            self.assertRaises(ValueError, _advanced_ssh_config_parse_options, '--host=localhost', '--port={}'.format(port))
            self.assertRaises(ValueError, _advanced_ssh_config_parse_options, '--host=localhost', '-p', port)

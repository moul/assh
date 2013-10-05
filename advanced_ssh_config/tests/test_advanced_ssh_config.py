# -*- coding: utf-8 -*-

import unittest

from advanced_ssh_config.advanced_ssh_config import AdvancedSshConfig


class TestAdvancedSshConfig(unittest.TestCase):

    def test_load_advanced_ssh_config(self):
        advssh = AdvancedSshConfig()
        self.assertIsInstance(advssh, AdvancedSshConfig)

    def test_routing_simple(self):
        advssh = AdvancedSshConfig(hostname='test',
                                   port=23,
                                   verbose=True,
                                   dry_run=True)
        routing = advssh.get_routing()
        self.assertEqual(routing['port'], 23)
        self.assertEqual(routing['hostname'], 'test')
        self.assertEqual(routing['reallocalcommand'], [''])
        self.assertEqual(routing['right_path'], [])
        self.assertEqual(routing['gateways'], ['direct'])

    # FIXME: test_routing_with_config
    # FIXME: test_routing_override_config
    # FIXME: test_connect
    # FIXME: test_dryrun
    # FIXME: test_verbose

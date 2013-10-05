# -*- coding: utf-8 -*-

import unittest

from advanced_ssh_config.advanced_ssh_config import AdvancedSshConfig


class TestAdvancedSshConfig(unittest.TestCase):

    def test_load_advanced_ssh_config(self):
        advssh = AdvancedSshConfig()
        self.assertIsInstance(advssh, AdvancedSshConfig)

    def test_routing(self):
        advssh = AdvancedSshConfig(hostname='test',
                                   port=23,
                                   verbose=True,
                                   dry_run=True)
        routing = advssh.get_routing()
        from pprint import pprint
        print()
        print('-' * 80)
        pprint(routing)
        print('-' * 80)

    # FIXME: test_dryrun
    # FIXME: test_verbose

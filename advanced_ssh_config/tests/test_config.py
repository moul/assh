# -*- coding: utf-8 -*-

import unittest
import os

from advanced_ssh_config.config import Config


class TestConfig(unittest.TestCase):

    def setUp(self):
        self.prefix = '/tmp/test-asc-config'
        os.system('rm -rf {}'.format(self.prefix))
        os.makedirs(self.prefix)
        with open('{}/config.advanced'.format(self.prefix), 'w') as f:
            f.write('')

    def test_initialize_config(self):
        config = Config(['{}/config.advanced'.format(self.prefix)])
        self.assertIsInstance(config, Config)

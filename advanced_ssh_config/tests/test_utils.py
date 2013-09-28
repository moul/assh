# -*- coding: utf-8 -*-

import unittest
import os

from advanced_ssh_config.utils import safe_makedirs


class TestSafeMakedirs(unittest.TestCase):

    def setUp(self):
        self.prefix = '/tmp/test-safe-makedirs'
        if os.path.exists(self.prefix):
            os.system('rm -rf {}'.format(self.prefix))
        os.makedirs(self.prefix)

    def test_already_exists(self):
        safe_makedirs('{}/dir'.format(self.prefix))
        safe_makedirs('{}/dir'.format(self.prefix))

    def test_invalid(self):
        for path in ('/dev/null/test',):
            self.assertRaises(OSError, safe_makedirs, path)

    def test_makedirs_on_file(self):
        open('{}/file'.format(self.prefix), 'w').write('hello')
        self.assertRaises(OSError, safe_makedirs, '{}/file/dir'.format(self.prefix))

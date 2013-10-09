# -*- coding: utf-8 -*-

import os


PREFIX = '/tmp/test-asc'
DEFAULT_CONFIG = os.path.join(PREFIX, 'config.advanced')


def write_config(contents, name='config.advanced'):
    with open(os.path.join(PREFIX, name), 'w') as f:
        f.write(contents)

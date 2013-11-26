# -*- coding: utf-8 -*-

import os

from advanced_ssh_config.config import Config


PREFIX = '/tmp/test-asc'
DEFAULT_CONFIG = os.path.join(PREFIX, 'config.advanced')


def write_config(contents, name='config.advanced'):
    with open(os.path.join(PREFIX, name), 'w') as f:
        f.write(contents)


def prepare_config():
    os.system('rm -rf {}'.format(PREFIX))
    os.makedirs(PREFIX)
    #write_config('')


def set_config(contents, load=True):
    contents = contents.format(PREFIX)
    write_config(contents)
    if load:
        config = Config([DEFAULT_CONFIG])
        return config
    return True

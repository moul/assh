#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import imp

from setuptools import setup, find_packages


sys.path.append(os.path.abspath(os.path.dirname(__file__)))

MODULE_NAME = 'advanced_ssh_config'
MODULE = imp.load_module(MODULE_NAME, *imp.find_module(MODULE_NAME))


setup(
    name=MODULE_NAME.replace('_', '-'),
    version=MODULE.__version__,
    author='Manfred Touron',
    author_email='m@42.am',
    url='https://github.com/moul/advanced-ssh-config',
    packages=find_packages(),
    classifiers=[
        # As from http://pypi.python.org/pypi?%3Aaction=list_classifiers
        'Development Status :: 3 - Alpha',
        'Intended Audience :: System Administrators',
        'License :: OSI Approved :: MIT License',
        'Operating System :: POSIX',
        'Operating System :: MacOS',
        'Operating System :: Unix',
        'Programming Language :: Python',
        'Topic :: Software Development :: Libraries :: Python Modules',
        'Topic :: Internet',
        'Topic :: System :: Systems Administration',
        'Topic :: System :: Shells',
    ],
    license='MIT',
    entry_points={
        'console_scripts': [
            'advanced-ssh-config = advanced_ssh_config.asc:main',
            ],
    },
)

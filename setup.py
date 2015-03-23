#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import imp

from setuptools import setup, find_packages


sys.path.append(
    os.path.abspath(os.path.dirname(__file__))
)


MODULE_NAME = 'advanced_ssh_config'
MODULE = imp.load_module(MODULE_NAME, *imp.find_module(MODULE_NAME))

DEPENDENCIES = [
]

TEST_DEPENDENCIES = [
    'coverage',
    'mock',
    'nose',
    'pep8',
]

description = "Add some magic to SSH and .ssh/config"
long_description = description
if os.path.exists('README.rst'):
    long_description = open('README.rst').read()


setup(
    name=MODULE_NAME.replace('_', '-'),
    description=description,
    long_description=long_description,
    version=MODULE.__version__,
    author='Manfred Touron',
    author_email='m@42.am',
    url='https://github.com/moul/advanced-ssh-config',
    download_url='https://github.com/moul/advanced-ssh-config/archive/master.zip',
    packages=find_packages(),
    package_dir={'advanced_ssh_config': 'advanced_ssh_config'},
    install_requires=DEPENDENCIES,
    tests_require=TEST_DEPENDENCIES,
    extras_require={
        'process_inspection': ['psutil'],
        'password': ['pexpect'],
        'release': ['PyInstaller'],
    },
    test_suite='{}.tests'.format(MODULE_NAME),
    classifiers=[
        # As from http://pypi.python.org/pypi?%3Aaction=list_classifiers
        'Development Status :: 3 - Alpha',
        'Intended Audience :: System Administrators',
        'License :: OSI Approved :: MIT License',
        'Operating System :: MacOS',
        'Operating System :: POSIX',
        'Operating System :: Unix',
        'Programming Language :: Python',
        'Topic :: Internet',
        'Topic :: Software Development :: Libraries :: Python Modules',
        'Topic :: System :: Shells',
        'Topic :: System :: Systems Administration',
    ],
    license='MIT',
    entry_points={
        'console_scripts': [
            'advanced-ssh-config = advanced_ssh_config.bin:advanced_ssh_config',
            'assh = advanced_ssh_config.bin:advanced_ssh_config',
            ],
    },
)

FROM python:2.7
MAINTAINER Manfred Touron "m@42.am"

RUN apt-get update && \
    apt-get -qq install openssh-client netcat-openbsd && \
    apt-get clean

RUN mkdir /advanced_ssh_config
WORKDIR /advanced_ssh_config
VOLUME /advanced_ssh_config

RUN echo '#!/bin/bash \n pep8 advanced_ssh_config | grep -v tests \n python setup.py test' > /test.sh && \
    chmod +x /test.sh

ADD . /advanced_ssh_config
RUN python setup.py install >/dev/null
RUN python setup.py develop >/dev/null
# Install test dependencies (probably a better way to do this)
RUN python setup.py test >/dev/null 2>/dev/null || true

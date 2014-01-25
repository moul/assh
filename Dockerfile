FROM mattias/python-dev
MAINTAINER Manfred Touron "m@42.am"

RUN apt-get -qq install openssh-client netcat-openbsd && \
    apt-get clean

RUN virtualenv /venv
RUN /venv/bin/pip install pep8 mock nose coverage -q

RUN mkdir /advanced_ssh_config
WORKDIR /advanced_ssh_config
VOLUME /advanced_ssh_config

CMD /bin/bash --rcfile /venv/bin/activate

ENV VIRTUAL_ENV /venv
ENV PATH /venv/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ADD . /advanced_ssh_config
RUN /venv/bin/python setup.py install >/dev/null
RUN /venv/bin/python setup.py develop >/dev/null
# Install test dependencies (probably a better way to do this)
RUN /venv/bin/python setup.py test >/dev/null 2>/dev/null || true

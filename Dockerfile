FROM mattias/python-dev
MAINTAINER Manfred Touron "m@42.am"

RUN apt-get -qq install openssh-client netcat-openbsd
RUN virtualenv /venv
RUN mkdir /advanced_ssh_config
WORKDIR /advanced_ssh_config
RUN /venv/bin/pip install pep8 mock nose coverage -q
VOLUME /advanced_ssh_config
CMD /bin/bash --rcfile /venv/bin/activate

ADD . /advanced_ssh_config
RUN /venv/bin/python setup.py install >/dev/null
RUN /venv/bin/python setup.py develop >/dev/null
# Install test dependencies (probably a better way to do this)
RUN /venv/bin/python setup.py test >/dev/null 2>/dev/null || true

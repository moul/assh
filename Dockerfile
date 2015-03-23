FROM python:2.7
MAINTAINER Manfred Touron "m@42.am"

RUN apt-get update && \
    apt-get -qq install openssh-client netcat-openbsd && \
    apt-get clean

RUN mkdir /code/
WORKDIR /code/

ADD setup.py /code/
ADD advanced_ssh_config/__init__.py /code/advanced_ssh_config/
RUN pip install .
RUN pip install -e '.[process_inspection]'
RUN pip install -e '.[release]'
# RUN pip install -e '.[tests]'
ADD . /code/
RUN chmod 777 /code/ && chmod -R 777 /code/build

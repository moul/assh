.PHONY: develop test release release_pypi release_binary fclean

all:	develop


develop:
	pip install -e .[release]


test:
	python2.7 setup.py test
	# FIXME: handle python 3


release: test release_binary release_pypi README.rst README.md


README.txt README.rst:	README.md
	-pandoc -o $@ $<


release_pypi:
	python2.7 setup.py register
	python2.7 setup.py sdist upload
	python2.6 setup.py bdist_egg upload
	python2.7 setup.py bdist_egg upload
	python2.6 setup.py bdist_wheel --python-tag=py26 upload
	python2.7 setup.py bdist_wheel --python-tag=py27 upload
	# FIXME: handle python 3


release_binary: dist/advanced-ssh-config-$(shell uname -s)-$(shell uname -m)


dist/advanced-ssh-config-$(shell uname -s)-$(shell uname -m): dist/advanced-ssh-config
	cp $< $@

dist/advanced-ssh-config: bin/advanced-ssh-config
	pyinstaller -F $<


fclean:
	-rm -rf ./dist/


install:
	pip install -e .

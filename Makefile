all:	develop


develop:
	pip install -e .


test:
	python2.7 setup.py test
	# FIXME: handle python 2.6 and 3.4


release: test dist/advanced-ssh-config release/pypi


release_pypi:
	python2.7 setup.py register
	python2.7 setup.py sdist upload
	python2.6 setup.py bdist_egg upload
	python2.7 setup.py bdist_egg upload
	python2.6 setup.py bdist_wheel --python-tag=py26 upload
	python2.7 setup.py bdist_wheel --python-tag=py27 upload
	# FIXME: handle python 3.4


dist/advanced-ssh-config: bin/advanced-ssh-config
	pyinstaller -F $<


fclean:
	-rm -rf ./dist/

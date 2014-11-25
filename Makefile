develop:
	pip install -e .

test:
	python2.7 setup.py test
	# FIXME: handle python 2.6 and 3.4

release:
	python2.7 setup.py register
	python2.7 setup.py sdist upload
	python2.6 setup.py bdist_egg upload
	python2.7 setup.py bdist_egg upload
	python2.6 setup.py bdist_wheel --python-tag=py26 upload
	python2.7 setup.py bdist_wheel --python-tag=py27 upload
	# FIXME: handle python 3.4

fclean:
	-rm -rf ./dist/

develop:
	pip install -e .

test:
	python setup.py test

release:
	python setup.py sdist register upload
	python2.6 setup.py bdist_egg register upload
	python2.7 setup.py bdist_egg register upload

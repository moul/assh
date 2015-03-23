.PHONY: all develop test test_local release release_pypi release_binary fclean
.PHONY: install shell


all:	develop


develop:
	pip install -e .[release]


test:	test_local


test_local:
	python2.7 setup.py test


# FIXME: add test_docker


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


release_binary: dist/assh-Darwin-x86_64 dist/assh-Linux-x86_64


dist/assh-Darwin-x86_64:	bin/assh setup.py advanced_ssh_config/__init__.py
	rm -rf venv dist/assh
	virtualenv venv
	venv/bin/pip install .
	venv/bin/pip install -e '.[process_inspection]'
	venv/bin/pip install -e '.[release]'
	venv/bin/pyinstaller -F bin/assh
	./dist/assh --version
	mv dist/assh $@


dist/assh-Linux-x86_64:		bin/assh setup.py advanced_ssh_config/__init__.py
	rm -f dist/assh
	mkdir -p $(PWD)/dist
	chmod 777 $(PWD)/dist
	docker build -t assh .
	docker run --rm \
	  -v $(PWD)/dist:/code/dist \
	  -u nobody \
	  --entrypoint pyinstaller \
	  assh \
	  -F bin/assh
	mv dist/assh $@
	docker run --rm \
	  -v $(PWD)/dist:/code/dist \
	  -u nobody \
	  --entrypoint $@ \
	  assh \
	  --version


shell:
	docker build -t assh .
	docker run -v $(PWD):/code -it --rm --entrypoint bash assh


fclean:
	-rm -rf ./dist/ venv/


install:
	pip install -e .

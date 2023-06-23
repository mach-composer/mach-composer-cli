.PHONY: all build clean install lint isort flake8 requirements
VERSION=$(shell git describe --tags 2> /dev/null || echo latest)
NAME=mach:$(VERSION)
# REPO_NAME=mach
# REPOSITORY=


all: clean install

clean:
	find . -name '*.pyc' -delete
	find . -name '__pycache__' -delete
	find . -name '*.egg-info' | xargs rm -rf

install: requirements
	pip install pip-tools
	pip-sync requirements_dev.txt requirements.txt

upgrade:
	pip install pip-tools pur
	pur -r requirements.in
	pip-compile -v --upgrade

requirements: requirements_dev.txt requirements.txt

requirements_dev.txt: requirements_dev.in requirements.txt
	pip install pip-tools
	pip-compile requirements_dev.in

requirements.txt: requirements.in
	pip install pip-tools
	pip-compile requirements.in

schema:
	python generate_schema.py

release-pypi:
	pip install twine wheel
	rm -rf build/* dist/*
	python setup.py sdist bdist_wheel
	twine upload dist/*

#
# Testing
#
test:
	py.test tests/

retest:
	py.test --lf -vvs tests/

coverage:
	py.test tests/ --cov=mach --cov-report=term-missing

mypy:
	mypy --config-file=mypy.ini src/mach

#
# Lint targets
#
format:
	isort src tests
	black src/ tests/

lint: flake8 isort mypy


isort:
	isort --check-only src tests

flake8:
	flake8 src/ tests/

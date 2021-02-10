#!/usr/bin/env python
import re

from setuptools import find_packages, setup

try:
    with open("README.md", "r") as fh:
        readme = re.sub(
            "^.. start-no-pypi.*^.. end-no-pypi", "", fh.read(), flags=re.M | re.S
        )
except IOError:
    readme = ""

install_requires = [
    "arrow==0.17.0",
    "attrs==20.3.0",
    "binaryornot==0.4.4",
    "certifi==2020.6.20",
    "chardet==3.0.4",
    "click==7.1.2",
    "colorama==0.4.4",
    "cookiecutter==1.7.2",
    "dataclasses-json==0.5.2",
    "dataclasses-jsonschema==2.13.0",
    "filelock==3.0.12",
    "idna==2.10",
    "jinja2-time==0.2.0",
    "jinja2==2.11.1",
    "jsonschema==3.2.0",
    "markupsafe==1.1.1",
    "marshmallow-enum==1.5.1",
    "marshmallow==3.8.0",
    "mypy-extensions==0.4.3",
    "poyo==0.5.0",
    "pyrsistent==0.17.3",
    "python-dateutil==2.8.1",
    "python-slugify==4.0.1",
    "pyyaml==5.3",
    "requests-file==1.5.1",
    "requests==2.23.0",
    "six==1.15.0",
    "stringcase==1.2.0",
    "text-unidecode==1.3",
    "tldextract==3.1.0",
    "typing-extensions==3.7.4.3",
    "typing-inspect==0.6.0",
    "urllib3==1.25.10",
]

setup(
    name="mach-composer",
    version="1.0.0-rc.7",
    author="Lab Digital B.V.",
    author_email="info@labdigital.nl",
    url="https://github.com/labd/mach-composer",
    description="MACH composer",
    long_description=readme,
    long_description_content_type='text/markdown',
    zip_safe=False,
    install_requires=install_requires,
    extras_require={},
    package_dir={"": "src"},
    packages=find_packages("src"),
    include_package_data=True,
    entry_points="""
        [console_scripts]
        mach=mach.commands:mach
    """,
    license="Proprietary",
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: Implementation :: CPython",
    ],
    project_urls={
        'Documentation': 'https://docs.machcomposer.io',
        'Source': 'https://github.com/labd/mach-composer',
    },
)

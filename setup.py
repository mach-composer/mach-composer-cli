#!/usr/bin/env python
import re

import pkg_resources
from setuptools import find_packages, setup

install_requires = []
with open("requirements.txt") as r:
    result = pkg_resources.parse_requirements(r)
    install_requires = [
        str(requirement) for requirement in pkg_resources.parse_requirements(r)
    ]

try:
    with open("README.md", "r") as fh:
        readme = re.sub(
            "^.. start-no-pypi.*^.. end-no-pypi", "", fh.read(), flags=re.M | re.S
        )
except IOError:
    readme = ""

setup(
    name="mach-composer",
    version="1.0.0-rc.10",
    author="Lab Digital B.V.",
    author_email="info@labdigital.nl",
    url="https://github.com/labd/mach-composer",
    description="MACH composer",
    long_description=readme,
    long_description_content_type="text/markdown",
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
        "Documentation": "https://docs.machcomposer.io",
        "Source": "https://github.com/labd/mach-composer",
    },
)

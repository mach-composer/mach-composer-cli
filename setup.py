#!/usr/bin/env python
from setuptools import find_packages, setup

try:
    with open("README.md", "r") as fh:
        readme = "\n".join(fh.readlines())
except IOError:
    readme = ""

setup(
    name="mach",
    version=None,
    author="Lab Digital B.V.",
    author_email="info@labdigital.nl",
    url="https://git.labdigital.nl/mach/mach/",
    description="",
    long_description=readme,
    zip_safe=False,
    install_requires=[],
    extras_require={},
    package_dir={"": "src"},
    packages=find_packages("src"),
    include_package_data=True,
    entry_points="""
        [console_scripts]
        mach=mach.commands:mach
    """,
    license="Proprietary",
    classifiers=["Private :: Do Not Upload", "License :: Other/Proprietary License"],
)

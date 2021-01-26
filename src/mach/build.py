import os.path
import subprocess
import sys
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from mach.types import MachConfig


def build_packages(config: "MachConfig"):
    for component in config.components:
        if not component.package_script:
            continue

        run_package_script(component.package_script)
        component.package_filename = os.path.abspath(component.package_filename)
        if not os.path.exists(component.package_filename):
            raise ValueError(f"The package_filename on {component.name} doesn't exist")


def run_package_script(package_script: str):
    p = subprocess.run(package_script, stdout=sys.stdout, stderr=sys.stderr, shell=True)
    p.check_returncode()

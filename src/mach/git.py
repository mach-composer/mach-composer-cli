import subprocess
from typing import List

import click
from mach import exceptions


def commit(message: str):
    result = _run(["git", "status", "--short"])
    if not result:
        click.echo("No changes detected, won't commit anything")
        return

    _run(["git", "commit", "-m", message])


def add(file: str):
    _run(["git", "add", file])


def _run(cmd: List) -> str:
    try:
        return subprocess.check_output(cmd)
    except subprocess.CalledProcessError as e:
        raise exceptions.MachError(f"Could not perform command: {e}")

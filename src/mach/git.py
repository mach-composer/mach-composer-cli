import json
import os
import subprocess
from dataclasses import dataclass
from typing import List

import click
from mach import exceptions

PRETTY_FMT = """format:{%n  "commit": "%H",%n  "author": "%aN <%aE>",%n  "date": "%ad",%n  "message": "%s"%n},"""  # noqa


@dataclass
class Commit:
    id: str
    msg: str


def commit(message: str):
    result = _run(["git", "status", "--short"])
    if not result:
        click.echo("No changes detected, won't commit anything")
        return

    _run(["git", "commit", "-m", message])


def add(file: str):
    _run(["git", "add", file])


def ensure_local(repo: str, dest: str, *, reference: str = ""):
    """Ensure the repository is present on the given dest."""
    if os.path.exists(dest):
        _run(["git", "pull"], cwd=dest)
    else:
        clone(repo, dest)

    if reference:
        _run(["git", "reset", "--hard", reference], cwd=dest)


def clone(repo: str, dest: str):
    _run(["git", "clone", repo, dest])


def history(dir: str, from_ref: str, *, branch: str = "") -> List[Commit]:
    if branch:
        _run(["git", "checkout", branch], cwd=dir)

    cmd = ["git", "log", f"--pretty={PRETTY_FMT}"]
    if from_ref:
        cmd.append(f"{from_ref}..{branch or ''}")

    result = _run(cmd, cwd=dir).decode("utf-8").rstrip(",")
    result = f"[{result}]"
    data = json.loads(result)
    return [
        Commit(
            id=_clean_commit_id(line["commit"]), msg=_clean_commit_msg(line["message"])
        )
        for line in data
    ]


def _clean_commit_msg(msg: str) -> str:
    return msg


def _clean_commit_id(commit_id: str) -> str:
    """Get the correct commit ID for this commit.

    It will trim the short_id since mach and the components are using a
    different commit id format (7 chars long).
    """
    return commit_id[:7]


def _run(cmd: List, *args, **kwargs) -> str:
    kwargs["stderr"] = subprocess.DEVNULL

    try:
        return subprocess.check_output(cmd, *args, **kwargs)
    except subprocess.CalledProcessError as e:
        raise exceptions.MachError(f"Could not perform command: {e}")

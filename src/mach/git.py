import os
import re
import subprocess
from dataclasses import dataclass
from email.utils import parseaddr
from typing import List, Optional

import click
from mach import exceptions

PRETTY_FMT = {
    "commit": "%H",
    "author": "%aN <%aE>",
    "date": "%ad",
    "message": "%s",
}

PRETTY_FMT_STR = "format:" + "|".join([fmt for fmt in PRETTY_FMT.values()])


class GitError(exceptions.MachError):
    pass


@dataclass
class Commit:
    id: str
    msg: str
    author: str


def commit(message: str):
    result = _run(["git", "status", "--short"])
    if not result:
        click.echo("No changes detected, won't commit anything")
        return

    _run(["git", "commit", "-m", message])


def add(file: str):
    _run(["git", "add", file])


def ensure_local(repo: str, dest: str):
    """Ensure the repository is present on the given dest."""
    reference = ""
    reference_match = re.match(r"(.*)(?:(?:@)(\w+))$", repo)
    if reference_match:
        repo, reference = reference_match.groups()

    if os.path.exists(dest):
        _run(["git", "pull"], cwd=dest)
    else:
        clone(repo, dest)

    if reference:
        try:
            _run(["git", "reset", "--hard", reference], cwd=dest)
        except GitError as e:
            raise GitError(f"Unable to swtich to reference {reference}: {e}")


def clone(repo: str, dest: str):
    _run(["git", "clone", repo, dest])


def history(dir: str, from_ref: str, *, branch: Optional[str] = "") -> List[Commit]:
    if branch:
        _run(["git", "checkout", branch], cwd=dir)

    cmd = ["git", "log", f"--pretty={PRETTY_FMT_STR}"]
    if from_ref:
        cmd.append(f"{from_ref}..{branch or ''}")

    lines = _run(cmd, cwd=dir).decode("utf-8").splitlines()
    commits = []
    for line in lines:
        commit_id, author, date, message = line.split("|", 3)
        author, _ = parseaddr(author)
        commits.append(
            Commit(
                id=_clean_commit_id(commit_id),
                author=author,
                msg=_clean_commit_msg(message),
            )
        )

    return commits


def _clean_commit_msg(msg: str) -> str:
    return msg


def _clean_commit_id(commit_id: str) -> str:
    """Get the correct commit ID for this commit.

    It will trim the short_id since mach and the components are using a
    different commit id format (7 chars long).
    """
    return commit_id[:7]


def _run(cmd: List, *args, **kwargs) -> bytes:
    kwargs["stderr"] = subprocess.STDOUT
    try:
        return subprocess.check_output(cmd, *args, **kwargs)
    except subprocess.CalledProcessError as e:
        raise GitError(e.output.decode() if e.output else str(e)) from e

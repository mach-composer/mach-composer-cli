import io
import os
import re
import tempfile
from contextlib import contextmanager
from pathlib import PurePath
from typing import Union

import click
import requests
import yaml
import yamlinclude
from mach import exceptions, git

EXTERNAL_RE = re.compile(r"^(git::)?(http|https)://")


class YamlIncludeConstructor(yamlinclude.YamlIncludeConstructor):
    def load(
        self,
        loader: Union[yaml.Loader, yaml.FullLoader],
        pathname: str,
        recursive: bool = False,
        encoding: str = "",
        reader: str = "",
    ):
        with resolve_file(pathname) as file:
            return super().load(loader, file, recursive, encoding, reader)


@contextmanager
def resolve_file(path):
    if not EXTERNAL_RE.match(path):
        yield path
        return

    git_match = re.match(r"^git::(.*)", path)
    if git_match:
        git_path = git_match.group(1)
        with tempfile.TemporaryDirectory() as tmpdir:
            yield resolve_git(git_path, tmpdir)
    elif path.startswith("http"):
        suffix = "." + path.rsplit(".", 1)[1]
        with tempfile.NamedTemporaryFile(suffix=suffix) as tmpfile:
            yield resolve_http(path, tmpfile)
    else:
        raise Exception(f"External path {path} not supported.")


def resolve_http(path: str, tmpfile: io.FileIO):
    resp = requests.get(path)
    resp.raise_for_status()
    tmpfile.write(resp.content)
    tmpfile.seek(0)
    return tmpfile.name


def resolve_git(path: str, tmpdir: str):
    """Resolve file served on Git.

    Return path of local file
    """
    match = re.match(r"(.*\/.*)\/\/(.*)$", path)
    if match:
        repo, file = match.groups()
    else:
        raise Exception(
            f"Missing file component in include {path}: can't be just the Git repository"
        )

    click.echo(f"Resolving {file} from {repo}")
    dest = os.path.join(tmpdir, "repo")

    try:
        git.ensure_local(repo, dest)
    except git.GitError as e:
        raise exceptions.MachError(f"Could not fetch {repo}: {e}") from e

    return os.path.join(dest, file)

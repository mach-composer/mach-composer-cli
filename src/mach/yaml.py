import io
import os
import re
import subprocess
import tempfile
from contextlib import contextmanager
from os.path import abspath, dirname
from typing import Tuple, Union

import click
import requests
import yaml
import yamlinclude
from mach import exceptions, git
from yaml.error import YAMLError  # noqa

EXTERNAL_RE = re.compile(r"^(git::)?(http|https)://")
INCLUDE_RE = re.compile(r"^(.*)\${include\((.*)\)}\s*$")

Encrypted = bool

__all__ = ["YAMLError", "load"]


class YamlFileIO(io.TextIOWrapper):
    """YAML IO wrapper to modify the yaml input stream before parsing it.

    We support the !include tag using the yamlinclude package.
    However, MACH will provide this option in the same fashion as variables can be defined
    using the `${...}` syntax: `${include(file.yml)}`.
    An added bonus here is that SOPS will not strip this syntax as it does with !include
    """

    def read(self, size=-1):
        lines = super().read(size)
        return "\n".join([self.parse_line(line) for line in lines.split("\n")])

    def parse_line(self, line: str) -> str:
        match = INCLUDE_RE.match(line)
        if match:
            prefix, include_part = match.groups()
            line = f"{prefix}!include {include_part}"
        return line


def load(file: str) -> Tuple[dict, Encrypted]:
    YamlIncludeConstructor.add_to_loader_class(
        loader_class=yaml.FullLoader,
        base_dir=abspath(dirname(file)),
    )
    encrypted = False

    with open(file, "r+b") as fh:
        data = _yaml_load(fh)

    if "sops" in data:
        click.echo("Detected SOPS encryption; decrypting...")
        data = _yaml_load(_sops_stream(file))
        encrypted = True

    return data, encrypted


def _yaml_load(iostream: io.IOBase):
    yaml_io = YamlFileIO(iostream)
    return yaml.full_load(yaml_io)


def _sops_stream(file: str, *args, **kwargs) -> bytes:
    kwargs["stderr"] = subprocess.STDOUT
    cmd = ["sops", "-d", file, "--output-type=yaml"]
    try:
        return io.BytesIO(subprocess.check_output(cmd, *args, **kwargs))
    except subprocess.CalledProcessError as e:
        raise exceptions.MachError(e.output.decode() if e.output else str(e)) from e


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

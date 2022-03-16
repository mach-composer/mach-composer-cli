import concurrent.futures
import os
import re
import subprocess
from abc import ABC
from dataclasses import dataclass
from typing import List, Optional, Tuple, Union

import click
from mach import cache, exceptions, git, parse
from mach.types import ComponentConfig, MachConfig
from mach.variables import ignore_variable_not_found

NAME_RE = re.compile(r".* name: [\"']?(.*)[\"']?")
VERSION_RE = re.compile(r"(\s*version: )([\"']?.*[\"']?)")

ROOT_BLOCK_START = re.compile(r"^\w+")

MAX_WORKERS = 4

Updates = List[Tuple[ComponentConfig, str]]


@dataclass
class UpdaterInput:
    data: Union[MachConfig, List[ComponentConfig]]
    file: str

    @property
    def is_mach_config(self):
        return isinstance(self.data, MachConfig)

    @property
    def file_encrypted(self):
        return self.is_mach_config and self.data.file_encrypted


def update_file(
    file: str,
    *,
    component_name: Optional[str],
    new_version: Optional[str],
    verbose=False,
    check_only=False,
):
    try:
        with ignore_variable_not_found():
            config = parse.parse_config_from_file(file)
        data = UpdaterInput(config, file)
    except exceptions.ParseError as e:
        # We might have a components yml as input, try to parse that
        try:
            components = parse.parse_components(file)
        except exceptions.ParseError:
            # Raise original error
            raise e

        data = UpdaterInput(components, file)

    if component_name and new_version:
        update_config_component(data, component_name, new_version)
    else:
        update_config_components(data, verbose=verbose, check_only=check_only)


def update_config_component(  # noqa: C901
    updater_input: UpdaterInput,
    component_name: str,
    new_version: str,
):
    config = updater_input.data

    component = config.get_component(component_name)
    if not component:
        raise exceptions.MachError(f"Could not find component {component_name}")

    if component.version == new_version:
        click.echo(f"Component {component_name} is already on version {new_version}.")

    click.echo(f"Updating {component_name} to version {new_version}...")

    updater_cls = updater_for(updater_input)
    updater_cls.apply_updates(updater_input, [(component, new_version)])


def update_config_components(  # noqa: C901
    updater_input: UpdaterInput,
    *,
    check_only=False,
    verbose=False,
):
    """
    Update a given MACH configuration file.

    :param config: The MACH configuration to update components for
    :param check_only: Only check for updates; don't update the file
    :param verbose: Enable verbose output
    """
    intro_msg = f"Checking updates for components in {updater_input.file}"
    click.echo(intro_msg)
    click.echo("-" * len(intro_msg))

    updates: Updates = _fetch_changes(updater_input)

    if not check_only:
        updater_cls = updater_for(updater_input)
        updater_cls.apply_updates(updater_input, updates)


def updater_for(updater_input: UpdaterInput):
    if updater_input.is_mach_config:
        return SopsUpdater if updater_input.file_encrypted else FileUpdater

    return ComponentsFileUpdater


def _fetch_changes(updater_input: UpdaterInput) -> Updates:
    cache_dir = cache.cache_dir_for(updater_input.file)

    if updater_input.is_mach_config:
        components = updater_input.data.components
    else:
        components = updater_input.data

    updates: Updates = []

    def _inner_update(component):
        outputs = [f"Updates for {component.name}..."]

        match = re.match(r"^git::(.*)", component.source)
        if not match:
            outputs.append(
                f"  Cannot check {component.name} component since it doesn't have a Git source defined"  # noqa
            )
            click.echo("\n".join(outputs))
            return

        component_dir = os.path.join(cache_dir, component.name)
        repo = match.group(1)
        match = re.match(r"(.*\/.*)(\/\/.*)$", repo)
        if match:
            repo = match.group(1)
        git.ensure_local(repo, component_dir)

        commits = git.history(component_dir, component.version, branch=component.branch)
        if not commits:
            outputs.append("  No updates\n")
            click.echo("\n".join(outputs))
            return

        for commit in commits:
            outputs.append(f"  {commit.id}: {commit.msg} <{commit.author}>")

        outputs.append("")
        click.echo("\n".join(outputs))
        updates.append((component, commits[0].id))

    with concurrent.futures.ThreadPoolExecutor(max_workers=MAX_WORKERS) as executor:
        executor.map(_inner_update, components)

    return updates


class BaseUpdater(ABC):
    """Updater which update component version in-place.

    We'll use a very basic search-and-replace based on regular expressions
    instead of the yaml parser to not mess with any formatting.
    """

    @classmethod
    def apply_updates(cls, updater_input: UpdaterInput, updates: Updates):
        """Apply given updates to the file."""
        instance = cls()
        instance.apply(updater_input, updates)

    def apply(self, updater_input: UpdaterInput, updates: Updates):
        raise NotImplementedError()


class SopsUpdater(BaseUpdater):
    def apply(self, updater_input: UpdaterInput, updates: Updates):
        assert updater_input.file_encrypted

        if not updater_input.is_mach_config:
            raise NotImplementedError(
                "SOPS support not supported yet for updating a components file"
            )

        component_indexes = {
            c.name: i for i, c in enumerate(updater_input.data.components)
        }

        for component, version in updates:

            index = component_indexes[component.name]
            cmd = [
                "sops",
                "--set",
                f'["components"][{index}]["version"] "{version}"',
                updater_input.file,
            ]

            try:
                subprocess.run(cmd, check=True)
            except subprocess.CalledProcessError as e:
                click.echo(f"Failed to update {component.name}: {e}")
            else:
                click.echo(f"Updated {component.name} to {version} using SOPS")


class BaseFileUpdater(BaseUpdater):
    """Updater which update component version in-place.

    We'll use a very basic search-and-replace based on regular expressions
    instead of the yaml parser to not mess with any formatting.
    """

    def apply(self, updater_input: UpdaterInput, updates: Updates):
        click.echo("Writing updated to file...")
        file = updater_input.file
        self.current_component: Optional[ComponentConfig] = None
        self.updates = {component.name: version for component, version in updates}
        self.applied = []
        self.component_map = {component.name: component for component, _ in updates}

        with open(file) as f:
            lines = f.readlines()

        newlines = [self.process_line(line) for line in lines]

        not_applied = set(self.updates.keys()) - set(self.applied)
        if not_applied:
            click.echo(
                click.style(
                    f"Unable to apply all updates to the components: {', '.join(not_applied)}",
                    fg="yellow",
                )
            )

        with open(file, mode="w") as f:
            for line in newlines:
                f.write(line)

    def process_line(self, line: str) -> str:
        raise NotImplementedError()

    def process_component_line(self, line: str):
        name_match = NAME_RE.match(line)
        if name_match:
            component_name = name_match.group(1)

            try:
                self.current_component = self.component_map[component_name]
            except KeyError:
                self.current_component = None

            return line

        if not self.current_component:
            return line

        match = VERSION_RE.match(line)
        if not match:
            return line

        assert self.current_component.version in match.group(2)

        try:
            new_version = self.updates[self.current_component.name]
        except KeyError:
            return line

        new_version = f'"{new_version}"'
        self.applied.append(self.current_component.name)
        return VERSION_RE.sub(rf"\g<1>{new_version}", line)


class FileUpdater(BaseFileUpdater):
    def __init__(self):
        self.in_components = False

    def process_line(self, line: str) -> str:
        if line.startswith("components:"):
            self.in_components = True
        elif ROOT_BLOCK_START.match(line):
            self.in_components = False
        elif self.in_components:
            return self.process_component_line(line)

        return line


class ComponentsFileUpdater(BaseFileUpdater):
    def process_line(self, line: str) -> str:
        return self.process_component_line(line)

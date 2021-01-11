import os
import re
import typing

import click
from mach import cache, exceptions, git
from mach.types import ComponentConfig, MachConfig

NAME_RE = re.compile(r".* name: [\"']?(.*)[\"']?")
VERSION_RE = re.compile(r"(\s*version: )([\"']?.*[\"']?)")

ROOT_BLOCK_START = re.compile(r"^\w+")


Updates = typing.List[typing.Tuple[ComponentConfig, str]]


def update_config_component(  # noqa: C901
    config: MachConfig,
    component_name: str,
    new_version: str,
):
    component = config.get_component(component_name)
    if not component:
        raise exceptions.MachError(f"Could not find component {component_name}")

    if component.version == new_version:
        click.echo(f"Component {component_name} is already on version {new_version}.")

    click.echo(f"Updating {component_name} to version {new_version}...")

    FileUpdater.apply_updates(config.file, [(component, new_version)])


def update_config_components(  # noqa: C901
    config: MachConfig,
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
    intro_msg = f"Checking updates for components in {config.file}"
    print(intro_msg)
    print("-" * len(intro_msg))

    updates: Updates = _fetch_changes(config)

    if not check_only:
        FileUpdater.apply_updates(config.file, updates)


def _fetch_changes(config: MachConfig) -> Updates:
    cache_dir = cache.cache_dir_for(config)
    updates: Updates = []

    for component in config.components:
        click.echo(f"Updates for {component.name}...")

        match = re.match(r"^git::(.*)", component.source)
        if not match:
            click.echo(
                f"  Cannot check {component.name} component since it doesn't have a Git source defined"  # noqa
            )
            continue

        component_dir = os.path.join(cache_dir, component.name)
        repo = match.group(1)
        match = re.match(r"(.*\/.*)(\/\/.*)$", repo)
        if match:
            repo = match.group(1)
        git.ensure_local(repo, component_dir)

        commits = git.history(component_dir, component.version, branch=component.branch)
        if not commits:
            click.echo("  No updates\n")
            continue

        for commit in commits:
            print(f"  {commit.id}: {commit.msg}")

        click.echo("")

        updates.append((component, commits[0].id))

    return updates


class FileUpdater:
    """Updater which update component version in-place.

    We'll use a very basic search-and-replace based on regular expressions
    instead of the yaml parser to not mess with any formatting.
    """

    @classmethod
    def apply_updates(cls, file: str, updates: Updates):
        """Apply given updates to the file."""
        instance = cls()
        instance._apply_updates(file, updates)

    def _apply_updates(self, file: str, updates: Updates):
        self.in_components = False
        self.current_component: typing.Optional[ComponentConfig] = None
        self.updates = {component.name: version for component, version in updates}
        self.component_map = {component.name: component for component, _ in updates}

        with open(file) as f:
            lines = f.readlines()

        newlines = []
        for line in lines:
            if line.startswith("components:"):
                self.in_components = True
            elif ROOT_BLOCK_START.match(line):
                self.in_components = False
            elif self.in_components:
                line = self.process_component_line(line)
            newlines.append(line)

        with open(file, mode="w") as f:
            for line in newlines:
                f.write(line)

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

        if new_version.isdigit():
            new_version = f'"{new_version}"'

        return VERSION_RE.sub(rf"\g<1>{new_version}", line)

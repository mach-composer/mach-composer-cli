import datetime
import os
import re
import typing
import urllib.parse

import requests
from mach import exceptions
from mach.types import ComponentConfig, MachConfig

NAME_RE = re.compile(r".*name: [\"']?(.*)[\"']?")
VERSION_RE = re.compile(r"(\s*version: )([\"']?.*[\"']?)")


Updates = typing.List[typing.Tuple[ComponentConfig, str]]


def update_config_components(  # noqa: C901
    config: MachConfig, check_only=False, verbose=False
):
    # TODO: It now only checks/outputs the available updates
    gitlab_access_key = os.environ.get("GITLAB_ACCESS_KEY")
    if not gitlab_access_key:
        raise exceptions.UpdateError(
            "Must set GITLAB_ACCESS_KEY environment variable to use this option"
        )

    intro_msg = f"Checking updates for components in {config.file}"
    print(intro_msg)
    print("-" * len(intro_msg))

    updates: Updates = []

    session = requests.Session()
    session.headers["Authorization"] = f"Bearer {gitlab_access_key}"

    for component in config.components:
        if not component.source.startswith("git::"):
            print(
                f"Cannot check {component.name} component since it doesn't have a Git source defined"  # noqa
            )
            continue

        match = re.match(r".*git.labdigital.nl/(.*)\.git.*", component.source)
        if not match:
            print(
                f"Cannot check {component.name} component: only supports git.labdigital.nl sources"
            )
            continue

        repo = match.group(1)
        repo_encoded = urllib.parse.quote(repo, safe="")

        try:
            resp = session.get(
                f"https://git.labdigital.nl/api/v4/projects/{repo_encoded}/repository/commits/",
                params={"ref_name": f"{component.version}..master"},
            )
            resp.raise_for_status()
        except EnvironmentError as e:
            print(f"Error checking {repo}: {e}")
            continue

        print(f"{component.name}:")

        resp_data = resp.json()
        if not resp_data:
            print("   no updates\n")
            continue

        for commit in resp_data:
            commit["commit_id"] = commit_id(commit)

            if verbose:
                commit["committed_date"] = datetime.datetime.fromisoformat(
                    commit["committed_date"]
                )
                print(
                    " * {commit_id}: {title} ({committed_date:%Y-%m-%d %H:%M} - {author_name})".format(  # noqa
                        **commit
                    )
                )
            else:
                print(" * {commit_id}: {title}".format(**commit))

        print("")

        updates.append((component, resp_data[0]["commit_id"]))

    if not check_only:
        FileUpdater.apply_updates(config.file, updates)


def commit_id(data: dict) -> str:
    """Get the correct commit ID for this commit.

    It will trim the short_id since mach and the components are using a
    different commit id format (7 chars long).
    """
    return data["short_id"][:7]


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

        with open(file, mode="w") as f:
            for line in lines:
                if line.startswith("components:"):
                    self.in_components = True
                elif self.in_components:
                    line = self.process_component_line(line)

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

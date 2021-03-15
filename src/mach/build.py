import os.path
import subprocess
import sys
from typing import TYPE_CHECKING, Optional

if TYPE_CHECKING:
    from mach.types import LocalArtifact, MachConfig


def build_packages(config: "MachConfig", restrict_components=None):
    for component in config.components:
        if restrict_components and component.name not in restrict_components:
            continue

        for artifact, artifact_cfg in component.artifacts.items():
            build_artifact(artifact_cfg)


def build_artifact(artifact: "LocalArtifact"):
    run_script(artifact.script, artifact.workdir)

    if artifact.workdir:
        artifact.filename = os.path.abspath(
            os.path.join(artifact.workdir, artifact.filename)
        )
    else:
        artifact.filename = os.path.abspath(artifact.filename)

    if not os.path.exists(artifact.filename):
        raise ValueError(f"The file {artifact.filename} doesn't exist")


def run_script(script: str, workdir: Optional[str]):
    p = subprocess.run(
        script, stdout=sys.stdout, stderr=sys.stderr, shell=True, cwd=workdir
    )
    p.check_returncode()

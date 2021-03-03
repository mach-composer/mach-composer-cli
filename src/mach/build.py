import os.path
import subprocess
import sys
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from mach.types import LocalArtifact, MachConfig


def build_packages(config: "MachConfig"):
    for component in config.components:
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


def run_script(script: str, workdir: str):
    p = subprocess.run(
        script, stdout=sys.stdout, stderr=sys.stderr, shell=True, cwd=workdir
    )
    p.check_returncode()

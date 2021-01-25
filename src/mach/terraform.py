import os
import re
import subprocess
import sys
from pathlib import Path
from typing import List, Optional, Union

import click
from mach.templates import setup_jinja
from mach.types import MachConfig, Site


def generate_terraform(config: MachConfig, *, site: str = None):
    """Generate Terraform file from template and reformat it."""
    env = setup_jinja()
    template = env.get_template("site.tf")
    sites = _filter_sites(config.sites, site)
    for site in sites:
        site_dir = config.deployment_path / Path(site.identifier)
        site_dir.mkdir(exist_ok=True)
        output_file = site_dir / Path("site.tf")
        content = _clean_tf(
            template.render(
                config=config,
                general_config=config.general_config,
                site=site,
            )
        )
        with open(output_file, "w+") as fh:
            fh.write(content)
        click.echo(f"Generated file {output_file}")

        run_terraform("fmt", cwd=site_dir)


def _clean_tf(content: str) -> str:
    """Clean the Terraform file.

    The pyhcl (used in testing for example) doesn't like empty objects with newlines
    for example. Let's get rid of those.
    """
    return re.sub(r"\{(\s*)\}", "{}", content)


def plan_terraform(
    config: MachConfig,
    *,
    site: str = None,
    components: List[str] = [],
    with_sp_login: bool = False,
    reuse=False,
):
    """Terraform init and plan for all generated sites."""
    sites = _filter_sites(config.sites, site)
    for site in sites:
        site_dir = config.deployment_path / Path(site.identifier)
        if not site_dir.is_dir():
            click.echo(f"Could not find site directory {site_dir}")
            continue

        click.echo(f"Terraform plan for {site_dir.name}")

        if not reuse:
            run_terraform("init", site_dir)

        if with_sp_login:
            azure_sp_login()

        cmd = ["plan"]
        for component in components:
            cmd.append(f"-target=module.{component}")

        run_terraform(cmd, site_dir)


def apply_terraform(
    config: MachConfig,
    *,
    site: str = None,
    components: List[str] = [],
    with_sp_login: bool = False,
    auto_approve: bool = False,
    reuse=False,
):
    """Terraform apply for all generated sites."""
    sites = _filter_sites(config.sites, site)
    for site in sites:
        site_dir = config.deployment_path / Path(site.identifier)
        if not site_dir.is_dir():
            click.echo(f"Could not find site directory {site_dir}")
            continue

        click.echo(f"Applying Terraform for {site.identifier}")

        if not reuse:
            run_terraform("init", site_dir)

        if with_sp_login:
            azure_sp_login()

        cmd = ["apply"]
        if auto_approve:
            cmd += ["-auto-approve"]

        for component in components:
            cmd.append(f"-target=module.{component}")
        run_terraform(cmd, site_dir)


def azure_sp_login():
    """Login the service principal with az cli."""
    p = subprocess.run(
        [
            "az",
            "login",
            "--service-principal",
            "-u",
            os.environ["ARM_CLIENT_ID"],
            "-p",
            os.environ["ARM_CLIENT_SECRET"],
            "--tenant",
            os.environ["ARM_TENANT_ID"],
        ],
        stdout=sys.stdout,
        stderr=sys.stderr,
    )
    p.check_returncode()


def run_terraform(command: Union[List[str], str], cwd):
    """Run any Terraform command."""
    if isinstance(command, str):
        command = [command]
    p = subprocess.run(
        ["terraform", *command], cwd=cwd, stdout=sys.stdout, stderr=sys.stderr
    )
    p.check_returncode()


def _filter_sites(sites: List[Site], site_identifier: Optional[str]):
    if not site_identifier:
        return sites

    return [s for s in sites if s.identifier == site_identifier]

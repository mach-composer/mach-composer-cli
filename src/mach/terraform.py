import os
import subprocess
import sys
from pathlib import Path
from typing import List, Union

from mach.templates import setup_jinja
from mach.types import MachConfig


def generate_terraform(config: MachConfig):
    """Generate Terraform file from template and reformat it."""
    env = setup_jinja()
    template = env.get_template("terraform_config.html")
    for site in config.sites:
        site_dir = config.deployment_path / Path(site.identifier)
        site_dir.mkdir(exist_ok=True)
        output_file = site_dir / Path("site.tf")
        with open(output_file, "w+") as fh:
            fh.write(
                template.render(
                    config=config, general_config=config.general_config, site=site
                )
            )
            print(f"Generated file {output_file}")

        run_terraform("fmt", cwd=site_dir)


def plan_terraform(output_dir: Path):
    """Terraform init and plan for all generated sites."""
    for site_dir in output_dir.iterdir():
        if site_dir.is_dir():
            print(f"Terraform plan for {site_dir.name}")
            run_terraform("init", site_dir)
            run_terraform("plan", site_dir)


def apply_terraform(output_dir: Path, with_sp_login: bool):
    """Terraform apply for all generated sites."""
    for site_dir in output_dir.iterdir():
        if site_dir.is_dir():
            print(f"Applying Terraform for {site_dir.name}")
            run_terraform("init", site_dir)
            if with_sp_login:
                azure_sp_login()
            run_terraform(["apply", "-auto-approve"], site_dir)


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

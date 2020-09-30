import glob
import subprocess
import sys
from functools import update_wrapper
from typing import List, Optional

import click
from mach import parse
from mach.exceptions import MachError
from mach.terraform import apply_terraform, generate_terraform, plan_terraform
from mach.update import update_config_components


@click.group()
def mach():
    pass


def terraform_command(f):
    @click.option(
        "-f",
        "--file",
        default=None,
        help="YAML file to parse. If not set parse all *.yml files.",
    )
    @click.option(
        "--with-sp-login",
        is_flag=True,
        default=False,
        help="If az login with service principal environment variables "
        "(ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID) should be done.",
    )
    @click.option(
        "--auto-approve",
        is_flag=True,
        default=False,
        help="",
    )
    @click.option(
        "--output-path",
        default="deployments",
        help="Output path, defaults to `cwd`/deployments.",
    )
    def new_func(file, with_sp_login: bool, auto_approve: bool, output_path: str):
        files = get_input_files(file)

        try:
            configs = parse.parse_configs(files, output_path)
        except MachError as e:
            raise click.ClickException(str(e)) from e

        try:
            result = f(
                file=file,
                with_sp_login=with_sp_login,
                auto_approve=auto_approve,
                configs=configs,
            )
        except subprocess.CalledProcessError as e:
            click.echo("Failed to run")
            sys.exit(e.returncode)
        else:
            click.echo("Done 👍")
            return result

    return update_wrapper(new_func, f)


@mach.command()
@terraform_command
def generate(file, configs, *args, **kwargs):
    for config in configs:
        generate_terraform(config)


@mach.command()
@terraform_command
def plan(file, configs, *args, **kwargs):
    for config in configs:
        generate_terraform(config)
        plan_terraform(config.deployment_path)


@mach.command()
@terraform_command
def apply(file, configs, with_sp_login, auto_approve, *args, **kwargs):
    for config in configs:
        generate_terraform(config)
        apply_terraform(
            config.deployment_path,
            with_sp_login=with_sp_login,
            auto_approve=auto_approve,
        )


@mach.command()
@click.option(
    "-f",
    "--file",
    default=None,
    help="YAML file to update. If not set update all *.yml files.",
)
@click.option("-v", "--verbose", is_flag=True, default=False, help="Verbose output.")
@click.option(
    "--check",
    default=False,
    is_flag=True,
    help="Only checks for updates, doesnt change files.",
)
def update(file: str, check: bool, verbose: bool):
    if not check:
        # TODO: Simply ignore the check option for now; won't update the config (yet)
        click.echo(
            "WARNING: Only supports checking/output of the available component updates"
        )

    files = get_input_files(file)
    configs = parse.parse_configs(files)
    try:
        for config in configs:
            update_config_components(config, verbose=verbose, check_only=check)
    except MachError as e:
        raise click.ClickException(str(e)) from e


def get_input_files(file: Optional[str]) -> List[str]:
    """Determine input files. If file is not specified use all *.yml files."""
    if file:
        files = [file]
    else:
        files = glob.glob("./*.yml")
    if not files:
        click.echo("No .yml files found")
        sys.exit(1)
    return files

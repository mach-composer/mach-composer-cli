import glob
import subprocess
import sys
from functools import update_wrapper
from os.path import splitext
from pathlib import Path
from typing import Dict, List, Optional

import click
import yaml
from mach.terraform import apply_terraform, generate_terraform, plan_terraform
from mach.types import ComponentConfig, MachConfig
from mach.update import UpdateError, update_config_components
from mach.validate import validate_config


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
        "--output-path",
        default="deployments",
        help="Output path, defaults to `cwd`/deployments.",
    )
    def new_func(file, with_sp_login: bool, output_path: str):
        files = get_input_files(file)
        configs = parse_configs(files, output_path)

        try:
            result = f(file=file, with_sp_login=with_sp_login, configs=configs)
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
def apply(file, configs, with_sp_login, *args, **kwargs):
    for config in configs:
        generate_terraform(config)
        apply_terraform(config.deployment_path, with_sp_login)


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
    configs = parse_configs(files)
    try:
        for config in configs:
            update_config_components(config, verbose=verbose, check_only=check)
    except UpdateError as e:
        raise click.ClickException(str(e)) from e


def parse_configs(files: List[str], output_path: str = None) -> List[MachConfig]:
    """Parse and validate configurations."""
    valid_configs = []
    for file in files:
        config = parse_config_from_file(file)
        config.file = file
        click.echo(f"Parsed {file} into config")

        validate_config(config)

        config = resolve_general_config(config)
        config = resolve_components(config)

        if output_path:
            full_output_path = Path(f"{output_path}/{splitext(file)[0]}")
            full_output_path.mkdir(exist_ok=True, parents=True)
            config.output_path = str(full_output_path)

        valid_configs.append(config)
    return valid_configs


def parse_config_from_file(file: str) -> MachConfig:
    """Parse file into MachConfig object."""
    click.echo(f"Got {file}")
    with open(file, "r") as fh:
        dictionary_config = yaml.full_load(fh)
    config = MachConfig.schema().load(dictionary_config)  # type: ignore
    return config


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


def resolve_components(config: MachConfig) -> MachConfig:
    """If no component info is specified, use global component settings."""
    component_info: Dict[str, ComponentConfig] = {
        component.name: component for component in config.components
    }
    for site in config.sites:
        if site.components:
            for component in site.components:
                if not component.version:
                    component.version = component_info[component.name].version
                if not component.source:
                    component.source = component_info[component.name].source
                if not component.short_name:
                    component.short_name = component_info[component.name].short_name
    return config


def resolve_general_config(config: MachConfig) -> MachConfig:
    """If no general config is specified, use global config settings."""
    if not config.general_config.azure:
        return config

    for site in config.sites:
        if (
            site.azure
            and config.general_config.azure.front_door
            and not site.azure.front_door
        ):
            site.azure.front_door = config.general_config.azure.front_door
    return config

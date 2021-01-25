import glob
import subprocess
import sys
from functools import update_wrapper
from typing import List, Optional

import click
from mach import bootstrap as _bootstrap
from mach import git, parse, updater
from mach.build import build_packages
from mach.exceptions import MachError
from mach.terraform import apply_terraform, generate_terraform, plan_terraform


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
        "-s",
        "--site",
        default=None,
        help="Site to parse. If not set parse all sites.",
    )
    @click.option(
        "--output-path",
        default="deployments",
        help="Output path, defaults to `cwd`/deployments.",
    )
    def new_func(file, site, output_path: str, **kwargs):
        files = get_input_files(file)

        try:
            configs = parse.parse_configs(files, output_path)
        except MachError as e:
            raise click.ClickException(str(e)) from e

        try:
            result = f(
                file=file,
                site=site,
                configs=configs,
                **kwargs,
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
def generate(file, site, configs, *args, **kwargs):
    """Generate the Terraform files."""
    for config in configs:
        generate_terraform(config, site=site)


@mach.command()
@click.option(
    "--with-sp-login",
    is_flag=True,
    default=False,
    help="If az login with service principal environment variables "
    "(ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID) should be done.",
)
@click.option(
    "-c",
    "--component",
    multiple=True,
    default=False,
    help="",
)
@click.option(
    "--reuse",
    default=False,
    is_flag=True,
    help="Supress a terraform init for improved speed (not recommended for production usage)",
)
@terraform_command
def plan(file, site, configs, with_sp_login, component, reuse, *args, **kwargs):
    """Output the deploy plan."""
    for config in configs:
        generate_terraform(config, site=site)
        plan_terraform(
            config,
            site=site,
            components=component,
            with_sp_login=with_sp_login,
            reuse=reuse,
        )


@mach.command()
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
    "-c",
    "--component",
    multiple=True,
    default=False,
    help="",
)
@click.option(
    "--reuse",
    is_flag=True,
    default=False,
    help="Supress a terraform init for improved speed (not recommended for production usage)",
)
@terraform_command
def apply(
    file,
    site,
    configs,
    with_sp_login,
    auto_approve,
    component,
    reuse,
    *args,
    **kwargs,
):
    """Apply the configuration."""
    for config in configs:
        build_packages(config)
        generate_terraform(config, site=site)
        apply_terraform(
            config,
            site=site,
            components=component,
            with_sp_login=with_sp_login,
            auto_approve=auto_approve,
            reuse=reuse,
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
@click.option(
    "-c",
    "--commit",
    default=False,
    is_flag=True,
    help="Automatically commits the change.",
)
@click.argument("component", required=False)
@click.argument("version", required=False)
def update(
    file: str,
    check: bool,
    verbose: bool,
    component: str,
    version: str,
    commit: bool,
):
    """Update all (or a given) component.

    When no component and version is given, it will check the git repositories for any updates.
    This command can also be used to manually update a single component by specifying a component
    and version.
    """
    if check and commit:
        raise click.ClickException(
            "check_only is not possible when create_commit is enabled."
        )

    if component and not version:
        raise click.ClickException(
            f"When specifying a component ({component}) you should specify a version as well"
        )

    files = get_input_files(file)
    configs = parse.parse_configs(files)
    try:
        for config in configs:
            if component and version:
                updater.update_config_component(config, component, version)
            else:
                updater.update_config_components(
                    config, verbose=verbose, check_only=check
                )
            if commit:
                git.add(config.file)

        if commit:
            if component:
                commit_msg = f"Updated {component} component"
            else:
                commit_msg = "Updated components"

            git.commit(commit_msg)

    except MachError as e:
        raise click.ClickException(str(e)) from e


@mach.command()
@click.option(
    "-f",
    "--file",
    default=None,
    help="YAML file to read. If not set read all *.yml files.",
)
def components(file: str):
    """List all components."""
    files = get_input_files(file)
    configs = parse.parse_configs(files)
    for config in configs:
        click.echo(f"{config.file}:")
        for component in config.components:
            click.echo(f" - {component.name}")
            click.echo(f"   version: {component.version}")

        click.echo("")


@mach.command()
@click.option(
    "-f",
    "--file",
    default=None,
    help="YAML file to read. If not set read all *.yml files.",
)
def sites(file: str):
    """List all sites."""
    files = get_input_files(file)
    configs = parse.parse_configs(files)
    for config in configs:
        click.echo(f"{config.file}:")
        for site in config.sites:
            click.echo(f" - {site.identifier}")
            click.echo("   components:")
            for component in site.components:
                click.echo(f"     {component.name}")

        click.echo("")


@mach.command()
@click.option(
    "-o",
    "--output",
    help="Output file or directory.",
)
@click.option(
    "-c",
    "--cookiecutter",
    default="",
    help="cookiecutter repository to generate from.",
)
@click.argument("type_", required=True, type=click.Choice(["config", "component"]))
def bootstrap(output: str, type_: str, cookiecutter: str):
    """Bootstraps a configuration or component."""
    if type_ == "config":
        _bootstrap.create_configuration(output or "main.yml")
    if type_ == "component":
        _bootstrap.create_component(output, cookiecutter)


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

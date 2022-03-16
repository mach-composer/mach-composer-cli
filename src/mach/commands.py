import glob
import subprocess
import sys
import textwrap
from functools import update_wrapper
from typing import List, Optional

import click
from mach import bootstrap as _bootstrap
from mach import git, parse, updater
from mach.__version__ import __version__
from mach.build import build_packages
from mach.exceptions import MachError, ParseError
from mach.terraform import (
    apply_terraform,
    generate_terraform,
    init_terraform,
    plan_terraform,
)
from mach.variables import ignore_variable_not_found


class Command(click.Command):
    def invoke(self, ctx):
        try:
            return super().invoke(ctx)
        except MachError as e:
            raise click.ClickException(str(e)) from e


class Group(click.Group):
    def command(self, *args, **kwargs):
        kwargs.setdefault("cls", Command)
        return super().command(*args, **kwargs)


@click.group(cls=Group)
@click.version_option(__version__)
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
        "--ignore-version",
        is_flag=True,
        default=False,
        help="Skip MACH composer version check",
    )
    @click.option(
        "--output-path",
        default="deployments",
        help="Output path, defaults to `cwd`/deployments.",
    )
    @click.option(
        "--var-file",
        help="Use a variable file to parse the configuration with",
    )
    def new_func(
        file, site, output_path: str, ignore_version: bool, var_file, **kwargs
    ):
        files = get_input_files(file, var_file=var_file)

        configs = parse.parse_configs(
            files, output_path, ignore_version=ignore_version, var_file=var_file
        )

        if not configs:
            raise click.ClickException("No valid MACH configurations found")

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
@terraform_command
def init(file, site, configs, *args, **kwargs):
    """Initialize site directories Terraform files."""
    for config in configs:
        generate_terraform(config, site=site)
        init_terraform(config, site=site)


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
    default=[],
    help="",
)
@click.option(
    "--reuse",
    default=False,
    is_flag=True,
    help="Supress a terraform init for improved speed (not recommended for production usage)",
)
@click.option(
    "--destroy",
    default=False,
    is_flag=True,
    help="Destroy option is a convenient way to destroy all remote objects managed by this mach config",  # noqa
)
@terraform_command
def plan(
    file, site, configs, with_sp_login, component, reuse, destroy, *args, **kwargs
):
    """Output the deploy plan."""
    for config in configs:
        generate_terraform(config, site=site)
        plan_terraform(
            config,
            site=site,
            components=component,
            with_sp_login=with_sp_login,
            reuse=reuse,
            destroy=destroy,
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
    default=[],
    help="",
)
@click.option(
    "--reuse",
    is_flag=True,
    default=False,
    help="Supress a terraform init for improved speed (not recommended for production usage)",
)
@click.option(
    "--destroy",
    default=False,
    is_flag=True,
    help="Destroy option is a convenient way to destroy all remote objects managed by this mach config",  # noqa
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
    destroy,
    *args,
    **kwargs,
):
    """Apply the configuration."""
    for config in configs:
        build_packages(config, restrict_components=component)
        generate_terraform(config, site=site)
        apply_terraform(
            config,
            site=site,
            components=component,
            with_sp_login=with_sp_login,
            auto_approve=auto_approve,
            reuse=reuse,
            destroy=destroy,
        )


# flake8: noqa: C901
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
@click.option(
    "-m", "--commit-message", default=None, help="Use a custom message for the commit."
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
    commit_message: str,
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
        if "@" not in component:
            raise click.ClickException(
                f"When specifying a component ({component}) you should specify a version as well"
            )
        component, version = component.split("@")

    for file in get_input_files(file):
        try:
            updater.update_file(
                file,
                component_name=component,
                new_version=version,
                verbose=verbose,
                check_only=check,
            )
        except ParseError as e:
            click.echo(textwrap.indent(str(e), "  "))
            continue

        if commit:
            git.add(file)

    if commit:
        if commit_message:
            commit_msg = commit_message
        elif component:
            commit_msg = f"Updated {component} component"
        else:
            commit_msg = "Updated components"

        git.commit(commit_msg)


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

    with ignore_variable_not_found():
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

    with ignore_variable_not_found():
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


@mach.command()
@click.argument("site", required=True)
@click.option(
    "-f",
    "--file",
    default=None,
    help="YAML file to read. If not set read all *.yml files.",
)
def path(file: str, site: str):
    """Return output path for given site."""
    # Suppress click output, we only want to output the output path
    click.echo = lambda *args, **kwargs: None
    files = get_input_files(file)
    configs = parse.parse_configs(files, "deployments")
    for config in configs:
        click.echo(f"{config.file}:")
        for site_ in config.sites:
            if site_.identifier == site:
                sys.stdout.write(f"{config.deployment_path_for(site_)}\n")
                return
    sys.exit(f"Could not find site definition {site}\n")


def get_input_files(file: Optional[str], *, var_file: str = None) -> List[str]:
    """Determine input files. If file is not specified use all *.yml files."""
    if file:
        files = [file]
    else:
        files = glob.glob("./*.yml")
    if not files:
        click.echo("No .yml files found")
        sys.exit(1)

    # If a var-file is given, strip it from the list of files to parse a MACH configurations.
    # This is mainly a convenience for when you have the following files;
    # - main.yml
    # - variables.yml
    # and to run `mach apply --var-file variables.yml`
    # instead of `mach apply -f main.yml --var-file variables.yml`
    if var_file:
        files = filter(lambda f: f.lstrip("./") != var_file, files)
    return files

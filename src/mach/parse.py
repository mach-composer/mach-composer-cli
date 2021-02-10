from collections import defaultdict
from os.path import abspath, basename, splitext
from pathlib import Path
from typing import Dict, List

import click
import yaml
from mach import exceptions
from mach.types import (
    CloudOption,
    ComponentConfig,
    Endpoint,
    MachConfig,
    SentryDsn,
    Site,
    SiteAzureSettings,
)
from mach.validate import validate_config
from marshmallow.exceptions import ValidationError


def parse_configs(files: List[str], output_path: str = None) -> List[MachConfig]:
    """Parse and validate configurations."""
    valid_configs = []
    for file in files:
        config = parse_config_from_file(file)
        config.file = file
        click.echo(f"Parsed {file} into config")
        validate_config(config)

        if output_path:
            full_output_path = Path(f"{output_path}/{splitext(basename(file))[0]}")
            full_output_path.mkdir(exist_ok=True, parents=True)
            config.output_path = str(full_output_path)

        valid_configs.append(config)
    return valid_configs


def parse_config_from_file(file: str) -> MachConfig:
    """Parse file into MachConfig object."""
    click.echo(f"Parsing {file}...")
    with open(file, "r") as fh:
        dictionary_config = yaml.full_load(fh)

    try:
        config = MachConfig.schema(infer_missing=True).load(dictionary_config)  # type: ignore
    except KeyError as e:
        # Most probably a missing value in the configuration.
        # dataclasses_json doesn't really give a proper Exception for this.
        # TODO: See if we can improve this / make it more robust. Either by improving
        # dataclassess_json (with a PR) or by extending it (if possible)
        raise exceptions.ParseError(f"Required attribute {e} missing") from e
    except ValidationError as e:
        # TODO: We don't have any path here, so not the best of error messages
        raise exceptions.ParseError(
            "Configuration file could not be validated", details=e.normalized_messages()
        ) from e

    return parse_config(config)


def parse_config(config: MachConfig) -> MachConfig:
    config = resolve_component_definitions(config)
    config = resolve_site_configs(config)
    return config


# flake8: noqa: C901
def resolve_site_configs(config: MachConfig) -> MachConfig:
    """Use and merge site-specific configurations with general config."""
    for site in config.sites:
        if config.general_config.cloud == CloudOption.AZURE:
            if site.azure:
                site.azure.merge(config.general_config.azure)
            else:
                site.azure = SiteAzureSettings.from_config(config.general_config.azure)

            if site.azure.resource_group:
                click.echo(
                    click.style(
                        (
                            f"WARNING: resource_group on {site.identifier} "
                            f"is used ({site.azure.resource_group}). "
                        ),
                        fg="red",
                        bold=True,
                    )
                )
                click.echo(
                    click.style(
                        (
                            "   Make sure it wasn't managed by MACH before otherwise "
                            "the resource group will get deleted."
                        ),
                        fg="red",
                    )
                )

        # Merge Contentful settings
        if config.general_config.contentful:
            for site in config.sites:
                if site.contentful:
                    site.contentful.merge(config.general_config.contentful)

        # Merge Amplience settings
        if config.general_config.amplience:
            for site in config.sites:
                if site.amplience:
                    site.amplience.merge(config.general_config.amplience)

        if config.general_config.sentry:
            if not site.sentry:
                site.sentry = SentryDsn.from_config(config.general_config.sentry)
            else:
                site.sentry.merge(config.general_config.sentry)

    config = resolve_site_components(config)

    for site in config.sites:
        resolve_endpoint_components(site)

    return config


def resolve_endpoint_components(site: Site):
    endpoint_components = defaultdict(list)
    for c in site.components:
        for endpoint in c.endpoints.values():
            endpoint_components[endpoint].append(c)

    site_endpoint_keys = {e.key for e in site.endpoints}
    # If one of the components has a 'default' endpoint defined,
    # we'll include it to our site endpoints.
    # A 'default' endpoint is one without a custom domain, so no further
    # Route53 or DNS zone settings required.
    if "default" in endpoint_components and not "default" in site_endpoint_keys:
        click.echo(
            click.style(
                (
                    "WARNING: 'default' endpoint used but not defined in the site endpoints.\n"
                    "MACH will create a default endpoint without any custom domain attached to it.\n"
                    "More info: https://docs.machcomposer.io/reference/syntax/sites.html#endpoints"
                ),
                fg="yellow",
            )
        )
        site.endpoints.append(
            Endpoint(
                url="",
                key="default",
            )
        )

    for endpoint in site.endpoints:
        endpoint.components = endpoint_components.get(endpoint.key, [])


def resolve_site_components(config: MachConfig) -> MachConfig:
    """If no component info is specified, use global component settings."""
    component_info: Dict[str, ComponentConfig] = {
        component.name: component for component in config.components
    }
    for site in config.sites:
        if not site.components:
            continue

        for component in site.components:
            try:
                info = component_info[component.name]
            except KeyError:
                raise exceptions.ParseError(
                    f"Component {component.name} does not exist in global components."
                )
            component.definition = info

            if not component.short_name:
                component.short_name = info.short_name

            if site.sentry:
                if not component.sentry:
                    component.sentry = site.sentry
                else:
                    component.sentry.merge(site.sentry)

    return config


def resolve_component_definitions(config: MachConfig) -> MachConfig:
    for comp in config.components:

        # Terraform needs absolute paths to modules
        if comp.source.startswith("."):
            comp.source = abspath(comp.source)

        if comp.integrations:
            continue

        # If no integrations are given, set the Cloud integrations as default
        if config.general_config.cloud == CloudOption.AWS:
            comp.integrations = ["aws"]
        elif config.general_config.cloud == CloudOption.AZURE:
            comp.integrations = ["azure"]

    return config

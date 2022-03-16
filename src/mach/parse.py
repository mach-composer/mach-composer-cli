import re
import textwrap
import warnings
from collections import defaultdict
from dataclasses import dataclass
from os.path import abspath, basename, splitext
from pathlib import Path
from typing import Any, Dict, List, Tuple

import click
from mach import exceptions, yaml
from mach.types import (
    CloudOption,
    ComponentAzureConfig,
    ComponentConfig,
    Endpoint,
    MachConfig,
    SentryDsn,
    ServicePlan,
    Site,
    SiteAzureSettings,
    TerraformReference,
)
from mach.validate import validate_config
from mach.variables import resolve_env_variable, resolve_variable
from marshmallow.exceptions import ValidationError

VARIABLE_RE = re.compile(r"^\${(var|env)\.(.*)}$")


def parse_components(file: str):
    yaml_data, _ = yaml.load(file)

    try:
        with warnings.catch_warnings():
            # Suppress a 'Unknown type ForwardRef('Component')' warning from dataclasses_json
            warnings.simplefilter("ignore")
            components = ComponentConfig.schema(infer_missing=True, many=True).load(
                yaml_data
            )
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

    return components


@dataclass
class VariablesFile:
    path: str
    vars: dict
    encrypted: bool

    @classmethod
    def from_file(cls, path: str):
        try:
            vars, vars_encrypted = yaml.load(path)
        except yaml.YAMLError as e:
            raise exceptions.ParseError(f"Could not parse variables file:\n{e}")

        return cls(path=path, vars=vars, encrypted=vars_encrypted)


def parse_configs(
    files: List[str],
    output_path: str = None,
    *,
    ignore_version=True,
    var_file: str = "",
) -> List[MachConfig]:
    """Parse and validate configurations."""
    vars = VariablesFile.from_file(var_file) if var_file else None

    configs = []
    for file in files:
        try:
            config = parse_and_validate(
                file,
                output_path,
                ignore_version=ignore_version,
                vars=vars,
            )
        except exceptions.ParseError as e:
            click.echo(textwrap.indent(str(e), "  "))
            continue

        configs.append(config)

    return configs


def parse_and_validate(
    file: str,
    output_path: str = None,
    *,
    ignore_version=True,
    vars: VariablesFile = None,
    vars_encrypted=False,
) -> MachConfig:
    """Parse and validate configuration."""
    yaml_content, encrypted = yaml.load(file)
    if not vars:
        vars = get_config_vars(yaml_content)

    config = parse_config_from_yaml(yaml_content, vars=vars)
    config.file_encrypted = encrypted
    config.file = file

    validate_config(config, ignore_version=ignore_version)

    if output_path:
        full_output_path = Path(f"{output_path}/{splitext(basename(file))[0]}")
        full_output_path.mkdir(exist_ok=True, parents=True)
        config.output_path = str(full_output_path)

    return config


def get_config_vars(yaml_content: dict) -> VariablesFile:
    var_file = yaml_content.get("mach_composer", {}).get("variables_file", None)
    if not var_file:
        return None

    return VariablesFile.from_file(var_file)


def parse_config_from_file(file: str, *, vars: VariablesFile = None) -> MachConfig:
    click.echo(f"Parsing {file}...")
    yaml_content, encrypted = yaml.load(file)
    config = parse_config_from_yaml(yaml_content, vars=vars)
    config.file_encrypted = encrypted
    config.file = file
    return config


def parse_config_from_yaml(
    yaml_content: dict, *, vars: VariablesFile = None
) -> MachConfig:
    """Parse file into MachConfig object."""
    try:
        with warnings.catch_warnings():
            # Suppress a 'Unknown type ForwardRef('Component')' warning from dataclasses_json
            warnings.simplefilter("ignore")
            config = MachConfig.schema(infer_missing=True).load(yaml_content)  # type: ignore

        if vars:
            config.variables = vars.vars
            config.variables_encrypted = vars.encrypted
            config.variables_path = vars.path

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
    resolve_variables(config, config.variables, config.variables_encrypted)
    parse_global_config(config)
    resolve_component_definitions(config)
    resolve_site_configs(config)
    return config


def resolve_variables(obj: Any, vars: dict, vars_encrypted: bool = False):
    """Resolve variables in the configuration.

    Only look for ${var.} variables since these can and must be rendered
    during parsing phase.

    This method will loop over objects that are either;
    - a dataclass (all configuration types are dataclasses)
    - a list of values
    - a dictionary

    ${component.} variables for example must be rendered in the Jinja template.
    """
    if isinstance(obj, str):
        var_m = VARIABLE_RE.match(obj.strip())
        if not var_m:
            # Return as is.
            return obj

        type_, var_name = var_m.groups()

        if type_ == "var":
            # We'll resolve the variable.
            # In case of encrypted vars we won't return the value as is,
            # but still the function will raise an error in case the variable cannot be found.
            # Better to catch it here then during terraform apply.
            var_value = resolve_variable(var_name, vars)

            if vars_encrypted:
                return TerraformReference(
                    f'data.sops_external.variables.data["{var_name}"]'
                )
            return var_value
        elif type_ == "env":
            return resolve_env_variable(var_name)
        else:
            raise exceptions.MachError(f"Unsupported variables type '{type_}': {obj}")

    annotations = getattr(obj, "__annotations__", {})
    if annotations:
        for field in annotations.keys():
            value = getattr(obj, field, None)
            if not value:
                continue
            value = resolve_variables(value, vars, vars_encrypted)
            setattr(obj, field, value)
    elif isinstance(obj, list):
        return [resolve_variables(v, vars, vars_encrypted) for v in obj]
    elif isinstance(obj, dict):
        return {k: resolve_variables(v, vars, vars_encrypted) for k, v in obj.items()}

    return obj


def parse_global_config(config: MachConfig):
    if config.general_config.cloud == CloudOption.AZURE:
        assert config.general_config.azure
        if "default" not in config.general_config.azure.service_plans:
            config.general_config.azure.service_plans["default"] = ServicePlan(
                kind="FunctionApp", tier="Dynamic", size="Y1"
            )


# flake8: noqa: C901
def resolve_site_configs(config: MachConfig):
    """Use and merge site-specific configurations with general config."""
    for site in config.sites:
        if config.general_config.cloud == CloudOption.AZURE:
            assert config.general_config.azure

            if site.azure:
                site.azure.merge(config.general_config.azure)
            else:
                site.azure = SiteAzureSettings.from_config(config.general_config.azure)

            assert site.azure

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
        for site in config.sites:
            if not site.sentry:
                site.sentry = SentryDsn.from_config(config.general_config.sentry)
            else:
                site.sentry.merge(config.general_config.sentry)

    resolve_site_components(config)

    for site in config.sites:
        resolve_endpoint_components(site)
        if site.azure:
            resolve_used_service_plans(site)


def resolve_used_service_plans(site: Site):
    "Azure-specific method to find out which service plans are actually being used by components"
    assert site.azure

    used = [
        c.azure.service_plan
        for c in site.components
        if c.azure and c.azure.service_plan
    ]
    site.azure.service_plans = {
        key: sp for key, sp in site.azure.service_plans.items() if key in used
    }


def resolve_endpoint_components(site: Site):
    endpoint_components = defaultdict(list)
    for c in site.components:
        for endpoint_key in c.endpoints.values():
            endpoint_components[endpoint_key].append(c)

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

            if site.sentry:
                if not component.sentry:
                    component.sentry = SentryDsn.from_config(site.sentry)
                else:
                    component.sentry.merge(site.sentry)

            if site.azure:
                if not component.azure:
                    component.azure = info.azure
                elif info.azure:
                    component.azure.merge(info.azure)

    return config


def resolve_component_definitions(config: MachConfig):
    for comp in config.components:
        # Terraform needs absolute paths to modules
        if comp.source.startswith("."):
            comp.source = abspath(comp.source)

        if not comp.integrations:
            # If no integrations are given, set the Cloud integrations as default
            if config.general_config.cloud == CloudOption.AWS:
                comp.integrations = ["aws"]
            elif config.general_config.cloud == CloudOption.AZURE:
                comp.integrations = ["azure"]

        if config.general_config.cloud == CloudOption.AZURE:
            if not comp.azure:
                comp.azure = ComponentAzureConfig()
            if not comp.azure.short_name:
                comp.azure.short_name = comp.name

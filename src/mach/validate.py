import re
from itertools import chain
from typing import Any, List

import click
from mach import types
from mach.__version__ import __version__
from mach.exceptions import ValidationError
from mach.variables import EmptyVar

STORE_KEY_RE = re.compile(r"^[\w_-]*$")


def check_variable(val: Any) -> bool:
    if type(val) == EmptyVar:
        return True
    return bool(val)


# Short-hand for readability
c = check_variable


def validate_config(config: types.MachConfig, *, ignore_version=True):
    """Check the config for invalid configuration."""
    if not ignore_version and config.mach_composer.version != __version__:
        raise ValidationError(
            f"MACH composer version defined in configuration ({config.mach_composer.version}) "
            f"does not match current version {__version__}"
        )

    validate_general_config(config.general_config)
    validate_components(config)

    for site in config.sites:
        validate_site(site, config=config)


def validate_general_config(config: types.GlobalConfig):
    if config.cloud == types.CloudOption.AZURE:
        if not config.azure:
            raise ValidationError("Missing azure configuration")
        if not config.terraform_config.azure_remote_state:
            raise ValidationError("Missing azure_remote_state configuration")
        if config.terraform_config.aws_remote_state:
            raise ValidationError(
                "Found aws_remote_state configuration, while cloud is set to 'azure'"
            )
    elif config.cloud == types.CloudOption.AWS:
        if not config.terraform_config.aws_remote_state:
            raise ValidationError("Missing aws_remote_state configuration")
        if config.terraform_config.azure_remote_state:
            raise ValidationError(
                "Found azure_remote_state configuration, while cloud is set to 'aws'"
            )

    if config.terraform_config:
        validate_terraform_config(config.terraform_config)

    if config.sentry:
        validate_sentry_config(config.sentry)


def validate_terraform_config(config: types.TerraformConfig):
    if config.providers:
        click.secho("Terraform provider versions", bold=True, fg="yellow")
        click.secho(
            "\n".join(
                [
                    "You are using custom Terraform provider versions.",
                    "Please be aware that some unexpected changes might occur compared to the MACH defaults.",  # noqa
                ]
            ),
            fg="yellow",
        )


def validate_site(site: types.Site, *, config: types.MachConfig):
    if config.general_config.cloud == types.CloudOption.AWS and not site.aws:
        raise ValidationError(f"Site {site.identifier} is missing an aws configuration")

    if config.general_config.cloud == types.CloudOption.AZURE:
        if site.endpoints and not site.azure.frontdoor:
            raise ValidationError(
                f"Site {site.identifier} has endpoints defined but frontdoor is not configured. "
                "Make sure your azure.frontdoor block is correcty configured in either the global "
                "or site configuration."
            )

    validate_endpoints(site, config.general_config.cloud)
    validate_commercetools(site)

    if site.components:
        validate_site_components(site.components, site=site)


def validate_endpoints(site: types.Site, cloud: types.CloudOption):
    # Construct lookup dictionary of all endpoints with the components that use them
    expected_endpoint_names = set(
        chain.from_iterable(c.endpoints.values() for c in site.components)
    )
    endpoint_names = {e.key for e in site.endpoints}

    missing = expected_endpoint_names - endpoint_names
    if missing:
        raise ValidationError(f"Missing required endpoints {', '.join(missing)}")


def validate_site_components(components: List[types.Component], *, site: types.Site):
    """Sanity checks on component configuration per site."""
    defined_stores = (
        [s.key for s in site.commercetools.stores] if site.commercetools else []
    )

    for component in components:
        if component.health_check_path and not component.health_check_path.startswith(
            "/"
        ):
            raise ValidationError(
                f"Component health check {component.health_check_path} does "
                "not start with '/'."
            )

        for store in component.store_variables.keys():
            if store not in defined_stores:
                raise ValidationError(
                    f"Store {store} is not defined in your commercetools stores definition"
                )

        if site.azure:
            service_plans = site.azure.service_plans
            if (
                component.azure
                and component.azure.service_plan
                and component.azure.service_plan not in service_plans
            ):
                raise ValidationError(
                    f"Component {component.name} requires service plan "
                    f"{component.azure.service_plan} which is not defined in the "
                    "Azure configuration."
                )


def validate_commercetools(site: types.Site):
    if site.commercetools:
        validate_store_keys(site.commercetools)
        validate_taxes(site.commercetools)


def validate_store_keys(ct_settings: types.CommercetoolsSettings):
    """Sanity checks on store values."""
    if ct_settings.stores:
        store_keys = [store.key for store in ct_settings.stores]
        for key in store_keys:
            if len(key) < 2:
                raise ValidationError(
                    f"Store key {key} should be minimum two characters."
                )
            if store_keys.count(key) != 1:
                raise ValidationError(f"Store key {key} must be unique.")

            if not STORE_KEY_RE.match(key):
                raise ValidationError(
                    f"Store key {key} may only contain alphanumeric characters or underscores"
                )


def validate_taxes(ct_settings: types.CommercetoolsSettings):
    """Ensure either `taxes` or `tax_categories` are used."""
    if ct_settings.taxes and ct_settings.tax_categories:
        raise ValidationError("`taxes` and `tax_categories` are mutually exclusive.")


def validate_components(config: types.MachConfig):
    """Validate global component data is valid."""
    if config.general_config.cloud == types.CloudOption.AWS:
        validate_aws_components(config)
    elif config.general_config.cloud == types.CloudOption.AZURE:
        validate_azure_components(config)


def validate_sentry_config(config: types.SentryConfig):
    if not any([c(config.dsn), c(config.auth_token)]):
        raise ValidationError("sentry: Either dsn or auth_token should be set")

    if all([c(config.dsn), c(config.auth_token)]):
        raise ValidationError("sentry: Only a dsn or auth_token should be defined")

    if config.auth_token and not any([c(config.project), c(config.organization)]):
        raise ValidationError(
            "sentry: A project and organization should be defined when using an auth_token"
        )


def validate_aws_components(config: types.MachConfig):
    """Validate components specifically for AWS usage."""
    pass


def validate_azure_components(config: types.MachConfig):
    """Validate components specifically for Azure usage.

    Only requirements for now is that a correct short_name should be set.
    Otherwise problems will arise when creating the Azure resources since
    for example Storage Accounts names have a limited length.
    """
    for comp in config.components:
        if comp.artifacts:
            raise ValidationError(
                f"The artifacts options on the '{comp.name}' component "
                "are not supported when targetting Azure"
            )

        if "azure" not in comp.integrations:
            continue

        # azure naming length is limited, so verify it doesn't get too long.
        assert comp.azure and comp.azure.short_name  # Should have been set by parser
        if len(comp.azure.short_name) > 10:
            raise ValidationError(
                f"Component {comp.name} short_name '{comp.azure.short_name}' "
                "cannot be more than 10 characters."
            )

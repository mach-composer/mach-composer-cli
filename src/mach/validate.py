from typing import List

from mach.types import CommercetoolsSettings, MachConfig, Site


def validate_config(config: MachConfig):
    """Check the config for invalid configuration."""
    validate_components(config)
    component_names = [component.name for component in config.components]

    for site in config.sites:
        validate_commercetools(site)
        validate_site_components(component_names, site)


def validate_site_components(component_names: List[str], site: Site):
    """Sanity checks on component configuration per site."""
    if site.components:
        for component in site.components:
            if component.name not in component_names:
                raise ValueError(
                    f"Component {component.name} does not exist in global components."
                )
            if (
                component.health_check_path
                and not component.health_check_path.startswith("/")
            ):
                raise ValueError(
                    f"Component health check {component.health_check_path} does "
                    f"not start with '/'."
                )


def validate_commercetools(site: Site):
    if site.commercetools:
        validate_store_keys(site.commercetools)


def validate_store_keys(ct_settings: CommercetoolsSettings):
    """Sanity checks on store values."""
    if ct_settings.stores:
        store_keys = [store.key for store in ct_settings.stores]
        for key in store_keys:
            if len(key) < 2:
                raise ValueError(f"Store key {key} should be minimum two characters.")
            if store_keys.count(key) != 1:
                raise ValueError(f"Store key {key} must be unique.")


def validate_components(config: MachConfig):
    """Validate global component data is valid."""
    for comp in config.components:
        # ignore product type modules
        # TODO: Make this somehow more generic / configurable
        if comp.short_name.startswith("product-types-"):
            continue

        # azure naming length is limited, so verify it doesn't get too long.
        if len(comp.short_name) > 10:
            raise ValueError(
                f"Component ({comp.name}) short name '{comp.short_name}' "
                f"cannot be more than 10 characters."
            )

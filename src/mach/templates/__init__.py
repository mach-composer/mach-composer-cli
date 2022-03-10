import os
import re
from typing import List

from jinja2 import Environment, FileSystemLoader
from jinja2.filters import do_mark_safe
from mach import utils
from mach.exceptions import MachError
from mach.types.sites import Component, Endpoint
from mach.types.values import TerraformReference

VARIABLE_RE = re.compile(r"^\${(component|var)\.(.*)}$")


def setup_jinja() -> Environment:
    templates_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)))
    env = Environment(
        loader=FileSystemLoader(templates_dir), trim_blocks=True, lstrip_blocks=True
    )
    load_filters(env)
    return env


def load_filters(env: Environment):
    env.filters.update(
        {
            "tfvalue": render_tfvalue,
            # Alias:
            "tf": render_tfvalue,
            "azure_region_long": azure_region_long,
            "azure_region_short": azure_region_short,
            "azure_frontend_endpoint_name": azure_frontend_endpoint_name,
            "zone_name": zone_name,
            "slugify": utils.slugify,
            "service_plan_resource_name": service_plan_resource_name,
            "render_commercetools_scopes": render_commercetools_scopes,
            "component_endpoint_name": component_endpoint_name,
        }
    )


def render_tfvalue(value):
    parsed_variable = parse_config_variable(value)
    if parsed_variable:
        return do_mark_safe(parsed_variable)

    if isinstance(value, TerraformReference):
        return value
    if isinstance(value, bool):
        return "true" if value else "false"
    if isinstance(value, (int, float)):
        return value
    if isinstance(value, list):
        values = ",".join([render_tfvalue(val) for val in value])
        return f"[{values}]"
    if isinstance(value, dict):
        values = "\n".join(
            [f"{key} = {render_tfvalue(val)}" for key, val in value.items()]
        )
        return f"{{\n{values}\n}}"
    return do_mark_safe(f'"{value}"')


def parse_config_variable(value: any):
    var_m = VARIABLE_RE.match(str(value).strip())
    if not var_m:
        return None

    type_, var = var_m.groups()
    if type_ == "component":
        parts = var.split(".", maxsplit=1)
        if len(parts) < 2:
            raise MachError(
                f"Invalid variable '{value}'; "
                "When using a ${component...} variable it has to consist of 2 parts; "
                "component-name.output-name"
            )
        return f"module.{parts[0]}.{parts[1]}"
    if type_ == "var":
        # This should not happen
        raise MachError(
            "${var...} variables should have been parsed already by MACH composer. "
            "Something went wrong."
        )

    return None


AZURE_REGION_DISPLAY_MAP_LONG = {
    "eastasia": "East Asia",
    "southeastasia": "Southeast Asia",
    "centralus": "Central US",
    "eastus": "East US",
    "eastus2": "East US 2",
    "westus": "West US",
    "northcentralus": "North Central US",
    "southcentralus": "South Central US",
    "northeurope": "North Europe",
    "westeurope": "West Europe",
    "japanwest": "Japan West",
    "japaneast": "Japan East",
    "brazilsouth": "Brazil South",
    "australiaeast": "Australia East",
    "australiasoutheast": "Australia Southeast",
    "southindia": "South India",
    "centralindia": "Central India",
    "westindia": "West India",
    "canadacentral": "Canada Central",
    "canadaeast": "Canada East",
    "uksouth": "UK South",
    "ukwest": "UK West",
    "westcentralus": "West Central US",
    "westus2": "West US 2",
    "koreacentral": "Korea Central",
    "koreasouth": "Korea South",
    "francecentral": "France Central",
    "francesouth": "France South",
    "australiacentral": "Australia Central",
    "australiacentral2": "Australia Central 2",
    "southafricanorth": "South Africa North",
    "southafricawest": "South Africa West",
}

AZURE_REGION_DISPLAY_MAP_SHORT = {
    "eastasia": "ea",
    "southeastasia": "sea",
    "centralus": "cus",
    "eastus": "eus",
    "eastus2": "eus2",
    "westus": "wus",
    "northcentralus": "ncus",
    "southcentralus": "scus",
    "northeurope": "ne",
    "westeurope": "we",
    "japanwest": "jw",
    "japaneast": "je",
    "brazilsouth": "bs",
    "australiaeast": "ae",
    "australiasoutheast": "ase",
    "southindia": "si",
    "centralindia": "ci",
    "westindia": "wi",
    "canadacentral": "cc",
    "canadaeast": "ce",
    "uksouth": "us",
    "ukwest": "uw",
    "westcentralus": "wc",
    "westus2": "wus2",
    "koreacentral": "kc",
    "koreasouth": "ks",
    "francecentral": "fc",
    "francesouth": "fs",
    "australiacentral": "ac",
    "australiacentral2": "ac2",
    "southafricanorth": "san",
    "southafricawest": "saw",
}


def azure_region_long(value):
    try:
        return AZURE_REGION_DISPLAY_MAP_LONG[value]
    except KeyError:
        raise


def azure_region_short(value):
    try:
        return AZURE_REGION_DISPLAY_MAP_SHORT[value]
    except KeyError:
        raise


def azure_frontend_endpoint_name(value: Endpoint):
    if value.azure and value.azure.internal_name:
        return render_tfvalue(value.azure.internal_name)
    return render_tfvalue(value.key)


def zone_name(value: str) -> str:
    value = utils.strip_protocol(value)
    return ".".join(value.split(".")[1:])


# Azure specific filters
def service_plan_resource_name(value: str) -> str:
    """Retreive the resource name for a Azure app service plan.

    The reason to make this conditional is because of backwards compatability;
    existing environments already have a `functionapp` resource. We want to keep that intact.
    """
    if value == "default":
        return "functionapps"
    return f"functionapps_{value}"


STORE_SUPPORTED_SCOPES = [
    "manage_orders",
    "manage_my_orders",
    "view_orders",
    "manage_customers",
    "view_customers",
    "manage_my_profile",
]


def render_commercetools_scopes(
    value: List[str], project_key: str, store_key: str = ""
):
    scopes = []
    for scope in value:
        scopes.append(f'"{scope}:{project_key}",')
        if store_key and scope in STORE_SUPPORTED_SCOPES:
            scopes.append(f'"{scope}:{project_key}:{store_key}",')

    return "[\n" + "".join(scopes) + "\n]"


def component_endpoint_name(
    component: Component,
    endpoint: Endpoint,
) -> str:
    """Take an component and an site-endpoint, and return a Terraform reference to an output.

    The endpoint might have a different name in the component itself based on the mappings
    """
    for component_key, ep_key in component.endpoints.items():
        if ep_key == endpoint.key:
            return component_key

    raise ValueError(f"Endpoint {endpoint.key} not found on {component.name}")

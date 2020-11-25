import os

from jinja2 import Environment, FileSystemLoader
from jinja2.filters import do_mark_safe
from mach import utils


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
            "component_value": render_value,
            "azure_region_long": azure_region_long,
            "azure_region_short": azure_region_short,
            "zone_name": zone_name,
            "slugify": utils.slugify,
        }
    )


def render_value(value):
    if isinstance(value, bool):
        return "true" if value else "false"
    if isinstance(value, (int, float)):
        return value
    return do_mark_safe(f'"{value}"')


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


def zone_name(value: str) -> str:
    value = utils.strip_protocol(value)
    return ".".join(value.split(".")[1:])

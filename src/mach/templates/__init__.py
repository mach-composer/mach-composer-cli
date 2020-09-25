import os

from jinja2 import Environment, FileSystemLoader
from jinja2.filters import do_mark_safe


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
        }
    )


def render_value(value):
    if isinstance(value, bool):
        return "true" if value else "false"
    if isinstance(value, (int, float)):
        return value
    return do_mark_safe(f'"{value}"')


AZURE_REGION_DISPLAY_MAP_LONG = {
    "westeurope": "West Europe",
    "northeurope": "North Europe",
}

AZURE_REGION_DISPLAY_MAP_SHORT = {
    "westeurope": "we",
    "northeurope": "ne",
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

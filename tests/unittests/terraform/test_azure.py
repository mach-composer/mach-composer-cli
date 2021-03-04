"""Test Azure specific configurations."""
from mach import parse, types

from tests.unittests.terraform import utils as tf


def test_generate_azure_service_plans(azure_config: types.MachConfig, tf_mock):
    config = azure_config
    config.components[0].azure = None
    config.sites[0].components.append(
        types.Component(
            name="payment",
        )
    )
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_resource_group.main",
    ]

    config.sites[0].components[0].azure.service_plan = "default"
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_resource_group.main",
    ]

    config.general_config.azure.service_plans["premium"] = types.ServicePlan(
        kind="Linux",
        tier="PremiumV2",
        size="P2v2",
    )
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_resource_group.main",
    ]

    config.sites[0].components[0].azure.service_plan = "premium"
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps_premium",
        "azurerm_resource_group.main",
    ]

    config.sites[0].components[1].azure.service_plan = "default"
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_app_service_plan.functionapps_premium",
        "azurerm_resource_group.main",
    ]

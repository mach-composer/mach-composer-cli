"""Test Azure specific configurations."""
from mach import parse, types

from tests.unittests.terraform import utils as tf


def test_generate_azure_service_plans(azure_config: types.MachConfig, tf_mock):
    config = azure_config
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_resource_group.main",
    ]

    config.components[0].azure = types.ComponentAzureConfig(service_plan="premium")
    config.general_config.azure.service_plans["premium"] = types.ServicePlan(
        kind="Linux",
        tier="PremiumV2",
        size="P2v2",
    )
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_app_service_plan.functionapps_premium",
        "azurerm_resource_group.main",
    ]

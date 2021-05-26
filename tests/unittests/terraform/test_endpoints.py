from mach import parse, types

from tests.unittests.terraform import utils as tf


def test_generate_aws_w_endpoints(config: types.MachConfig, tf_mock):
    config.sites[0].endpoints = [
        types.Endpoint(key="public", url="api.mach-example.com")
    ]
    data = tf.generate(parse.parse_config(config))

    # 'public' endpoint not used in component yet; no resources created
    assert "resource" not in data

    config.components[0].endpoints = {
        "main": "public",
    }
    data = tf.generate(parse.parse_config(config))

    # API gateway items need to be created since a component now uses it
    expected_resources = [
        "aws_acm_certificate.public",
        "aws_apigatewayv2_api.public_gateway",
        "aws_apigatewayv2_api_mapping.public",
        "aws_apigatewayv2_domain_name.public",
        "aws_apigatewayv2_route.public_application",
        "aws_apigatewayv2_stage.public_default",
        "aws_route53_record.public",
        "aws_route53_record.public_acm_validation",
    ]
    expected_data_sources = [
        "aws_route53_zone.mach_examplecom",
    ]
    assert tf.get_resource_ids(data) == expected_resources
    assert tf.get_data_ids(data) == expected_data_sources

    config.sites[0].endpoints.append(
        types.Endpoint(key="private", url="private-api.mach-services.io")
    )
    data = tf.generate(parse.parse_config(config))

    # We've added an extra endpoint definition, but hasn't been used.
    # List of resources should be the same as previous check
    assert tf.get_resource_ids(data) == expected_resources
    assert tf.get_data_ids(data) == expected_data_sources

    config.components.append(
        types.ComponentConfig(
            name="logger",
            source="some-source//terraform",
            version="1.0",
            endpoints={
                "main": "private",
            },
        )
    )
    config.sites[0].components.append(
        types.Component(
            name="logger",
        )
    )
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "aws_acm_certificate.private",
        "aws_acm_certificate.public",
        "aws_apigatewayv2_api.private_gateway",
        "aws_apigatewayv2_api.public_gateway",
        "aws_apigatewayv2_api_mapping.private",
        "aws_apigatewayv2_api_mapping.public",
        "aws_apigatewayv2_domain_name.private",
        "aws_apigatewayv2_domain_name.public",
        "aws_apigatewayv2_route.private_application",
        "aws_apigatewayv2_route.public_application",
        "aws_apigatewayv2_stage.private_default",
        "aws_apigatewayv2_stage.public_default",
        "aws_route53_record.private",
        "aws_route53_record.private_acm_validation",
        "aws_route53_record.public",
        "aws_route53_record.public_acm_validation",
    ]
    assert tf.get_data_ids(data) == [
        "aws_route53_zone.mach_examplecom",
        "aws_route53_zone.mach_servicesio",
    ]


def test_generate_azure_w_endpoints(azure_config: types.MachConfig, tf_mock):
    config = azure_config
    config.sites[0].endpoints = [
        types.Endpoint(key="public", url="api.mach-example.com")
    ]
    data = tf.generate(parse.parse_config(config))

    # 'public' endpoint not used in component yet; no resources created
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_resource_group.main",
    ]

    config.components[0].endpoints = {
        "main": "public",
    }
    data = tf.generate(parse.parse_config(config))

    # Frontdoor instance need to be created since a component now uses it
    expected_resources = [
        "azurerm_app_service_plan.functionapps",
        "azurerm_dns_cname_record.public",
        "azurerm_frontdoor.app-service",
        "azurerm_frontdoor_custom_https_configuration.public",
        "azurerm_resource_group.main",
    ]
    assert tf.get_resource_ids(data) == expected_resources

    config.sites[0].endpoints.append(
        types.Endpoint(key="private", url="private-api.mach-example.com")
    )
    data = tf.generate(parse.parse_config(config))

    # We've added an extra endpoint definition, but hasn't been used.
    # List of resources should be the same as previous check
    assert tf.get_resource_ids(data) == expected_resources

    config.components.append(
        types.ComponentConfig(
            name="logger",
            source="some-source//terraform",
            version="1.0",
            endpoints={
                "main": "private",
            },
        )
    )
    config.sites[0].components.append(
        types.Component(
            name="logger",
        )
    )
    data = tf.generate(parse.parse_config(config))
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_dns_cname_record.private",
        "azurerm_dns_cname_record.public",
        "azurerm_frontdoor.app-service",
        "azurerm_frontdoor_custom_https_configuration.private",
        "azurerm_frontdoor_custom_https_configuration.public",
        "azurerm_resource_group.main",
    ]


def test_generate_aws_w_default_endpoint(config: types.MachConfig, tf_mock):
    """When endpoint 'default' is used, no custom domain has to be set."""
    data = tf.generate(parse.parse_config(config))

    # 'public' endpoint not used in component yet; no resources created
    assert "resource" not in data

    config.components[0].endpoints = {
        "main": "default",
    }
    data = tf.generate(parse.parse_config(config))

    # API gateway items need to be created since a component now uses it
    expected_resources = [
        "aws_apigatewayv2_api.default_gateway",
        "aws_apigatewayv2_route.default_application",
        "aws_apigatewayv2_stage.default_default",
    ]
    assert tf.get_resource_ids(data) == expected_resources


def test_generate_azure_w_default_endpoint(azure_config: types.MachConfig, tf_mock):
    """When endpoint 'default' is used, no custom domain has to be set."""
    config = azure_config
    data = tf.generate(parse.parse_config(config))

    # 'public' endpoint not used in component yet; no resources created
    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_resource_group.main",
    ]

    config.components[0].endpoints = {
        "main": "default",
    }
    data = tf.generate(parse.parse_config(config))

    assert tf.get_resource_ids(data) == [
        "azurerm_app_service_plan.functionapps",
        "azurerm_frontdoor.app-service",
        "azurerm_resource_group.main",
    ]

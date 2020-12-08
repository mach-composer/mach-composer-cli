import json
import os
from typing import List

import pytest
from mach import parse, terraform, types

from tests.utils import HclWrapper, load_hcl


@pytest.fixture
def tf_mock(mocker):
    return mocker.patch("mach.terraform.run_terraform")


def _generate(config: types.MachConfig) -> HclWrapper:
    terraform.generate_terraform(config)
    file_path = os.path.join(config.output_path, config.sites[0].identifier, "site.tf")
    assert os.path.exists(file_path), "No site.tf file found"
    return load_hcl(file_path)


def test_generate_terraform(parsed_config: types.MachConfig, tf_mock):
    data = _generate(parsed_config)
    tf_mock.assert_called_once()

    assert data == {
        "terraform": [
            {
                "backend": [
                    {
                        "s3": {
                            "bucket": ["unittest"],
                            "key": ["test/unittest-nl"],
                            "region": ["eu-west-1"],
                            "encrypt": [True],
                        }
                    }
                ]
            },
            {"required_providers": [{}]},
        ],
        "provider": [{"aws": {"region": ["eu-central-1"], "version": ["~> 3.20.0"]}}],
        "module": [
            {
                "api-extensions": {
                    "source": ["some-source//terraform"],
                    "component_version": ["1.0"],
                    "environment": ["test"],
                    "site": ["unittest-nl"],
                    "variables": [{}],
                    "secrets": [{}],
                    "providers": [{"aws": "${aws}"}],
                    "depends_on": [[]],
                }
            }
        ],
    }


def test_generate_w_sentry(parsed_config: types.MachConfig, tf_mock):
    data = _generate(parsed_config)
    data_str = json.dumps(data)
    assert "sentry" not in data_str

    parsed_config.components[0].integrations = ["aws", "sentry"]
    data = _generate(parsed_config)

    assert "sentry_dsn" in data.module["api-extensions"]
    assert "sentry_key" not in data.get("resource", {})

    parsed_config.general_config.sentry = types.SentryConfig(
        auth_token="12345",
        organization="labd",
        project="unittest",
    )
    data = _generate(parsed_config)
    assert "sentry_dsn" in data.module["api-extensions"]
    assert "sentry_key" in data.resource
    assert "api-extensions" in data.resource.sentry_key
    sentry_data = data.resource.sentry_key["api-extensions"]
    assert "rate_limit_window" not in sentry_data
    assert "rate_limit_count" not in sentry_data

    comp_sentry = parsed_config.sites[0].components[0].sentry
    comp_sentry.rate_limit_window = 21600
    comp_sentry.rate_limit_count = 100

    data = _generate(parsed_config)
    sentry_data = data.resource.sentry_key["api-extensions"]
    assert sentry_data["rate_limit_window"] == 21600
    assert sentry_data["rate_limit_count"] == 100


def test_generate_w_endpoints(config: types.MachConfig, tf_mock):
    config.sites[0].endpoints = [
        types.Endpoint(key="public", url="api.mach-example.com")
    ]
    data = _generate(parse.parse_config(config))

    # 'public' endpoint not used in component yet; no resources created
    assert "resource" not in data

    config.components[0].endpoint = "public"
    data = _generate(parse.parse_config(config))

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
    assert _get_resource_ids(data) == expected_resources

    config.sites[0].endpoints.append(
        types.Endpoint(key="private", url="private-api.mach-example.com")
    )
    data = _generate(parse.parse_config(config))

    # We've added an extra endpoint definition, but hasn't been used.
    # List of resources should be the same as previous check
    assert _get_resource_ids(data) == expected_resources

    config.components.append(
        types.ComponentConfig(
            name="logger",
            source="some-source//terraform",
            version="1.0",
            endpoint="private",
        )
    )
    config.sites[0].components.append(
        types.Component(
            name="logger",
        )
    )
    data = _generate(parse.parse_config(config))
    assert _get_resource_ids(data) == [
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


def test_generate_w_stores(config: types.MachConfig, tf_mock):
    config.sites[0].commercetools = types.CommercetoolsSettings(
        project_key="ct-unit-test",
        client_id="a96e59be-24da-4f41-a6cf-d61d7b6e1766",
        client_secret="98c32de8-1a6c-45a9-a718-d3cce5201799",
        scopes="manage_project:ct-unit-test",
        languages=["nl-NL"],
        countries=["NL"],
        currencies=["EUR"],
        stores=[
            types.Store(
                name={
                    "en-GB": "Default store",
                },
                key="main-store",
            ),
            types.Store(
                name={
                    "en-GB": "Some other store",
                },
                key="other-store",
            ),
            types.Store(
                name={
                    "en-GB": "Forgotten store",
                },
                key="forgotten-store",
            ),
        ],
    )
    config.components[0].integrations = ["aws", "commercetools"]
    data = _generate(parse.parse_config(config))

    assert len(data.resource.commercetools_store) == 3
    assert "main-store" in data.resource.commercetools_store
    assert "other-store" in data.resource.commercetools_store
    assert "forgotten-store" in data.resource.commercetools_store

    assert len(data.module["api-extensions"].ct_stores) == 3

    for store_key, store in data.module["api-extensions"].ct_stores.items():
        assert store["key"] == store_key
        assert not store["variables"]
        assert not store["secrets"]

    config.sites[0].components[0].store_variables = {
        "main-store": {
            "FOO": "BAR",
            "EXTRA": "VALUES",
        },
        "other-store": {
            "FOO": "SOMETHING ELSE",
        },
    }
    config.sites[0].components[0].store_secrets = {
        "main-store": {
            "PAYMENT_KEY": "TLrlDf6XhKkXFGGHeQGY",
        },
    }

    data = _generate(parse.parse_config(config))
    main_store = data.module["api-extensions"].ct_stores["main-store"]
    other_store = data.module["api-extensions"].ct_stores["other-store"]
    assert len(main_store.variables) == 2
    assert len(other_store.variables) == 1
    assert len(main_store.secrets) == 1
    assert not other_store.secrets


def _get_resource_ids(data: HclWrapper) -> List[str]:
    """Get all resource ids in <resource-type>.<name> format."""
    result = []
    for type_, resources in data.resource.items():
        for key, resource in resources.items():
            result.append(f"{type_}.{key}")
    return sorted(result)

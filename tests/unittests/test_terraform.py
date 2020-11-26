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
        "provider": [{"aws": {"region": ["eu-central-1"], "version": ["~> 3.8.0"]}}],
        "module": [
            {
                "api-extensions": {
                    "source": ["some-source//terraform"],
                    "component_version": ["1.0"],
                    "environment": ["test"],
                    "site": ["unittest-nl"],
                    "variables": [{}],
                    "environment_variables": [{}],
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


def test_generate_w_endpoints(parsed_config: types.MachConfig, tf_mock):
    config = parsed_config
    config.sites[0].endpoints = {
        "public": "api.mach-example.com",
    }
    data = _generate(config)

    # 'public' endpoint not used in component yet; no resources created
    assert "resource" not in data

    config.components[0].endpoint = "public"
    data = _generate(config)

    # API gateway items need to be created since a component now uses it
    expected_resources = [
        "aws_acm_certificate.public",
        "aws_apigatewayv2_api.public_gateway",
        "aws_apigatewayv2_api_mapping.public",
        "aws_apigatewayv2_deployment.public_default",
        "aws_apigatewayv2_domain_name.public",
        "aws_apigatewayv2_route.public_application",
        "aws_apigatewayv2_stage.public_default",
        "aws_route53_record.public",
        "aws_route53_record.public_acm_validation",
    ]
    assert _get_resource_ids(data) == expected_resources

    config.sites[0].endpoints["private"] = "private-api.mach-example.com"
    data = _generate(config)

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
        "aws_apigatewayv2_deployment.private_default",
        "aws_apigatewayv2_deployment.public_default",
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


def _get_resource_ids(data: HclWrapper) -> List[str]:
    """Get all resource ids in <resource-type>.<name> format."""
    result = []
    for type_, resources in data.resource.items():
        for key, resource in resources.items():
            result.append(f"{type_}.{key}")
    return sorted(result)

import json
import os

import pytest
from mach import terraform, types

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

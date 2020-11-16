import json
import os

import hcl
import pytest
from mach import terraform, types


@pytest.fixture
def tf_mock(mocker):
    return mocker.patch("mach.terraform.run_terraform")


def _generate(config: types.MachConfig) -> dict:
    terraform.generate_terraform(config)
    file_path = os.path.join(config.output_path, config.sites[0].identifier, "site.tf")
    os.path.exists(file_path)
    with open(file_path) as f:
        return hcl.load(f)


def test_generate_terraform(parsed_config: types.MachConfig, tf_mock):
    data = _generate(parsed_config)
    tf_mock.assert_called_once()

    assert data == {
        "terraform": {
            "backend": {
                "s3": {
                    "bucket": "unittest",
                    "encrypt": True,
                    "key": "test/unittest-nl",
                    "region": "eu-west-1",
                }
            },
            "required_providers": {},
        },
        "module": {
            "api-extensions": {
                "depends_on": [],
                "providers": {},
                "source": "some-source//terraform",
            }
        },
    }


def test_generate_w_sentry(parsed_config: types.MachConfig, tf_mock):
    data = _generate(parsed_config)
    data_str = json.dumps(data)
    assert "sentry" not in data_str

    parsed_config.components[0].integrations = ["aws", "sentry"]
    data = _generate(parsed_config)
    assert "sentry_dsn" in data["module"]["api-extensions"]
    assert "sentry_key" not in data.get("resource", {})

    parsed_config.general_config.sentry = types.SentryConfig(
        auth_token="12345",
        organization="labd",
        project="unittest",
    )
    data = _generate(parsed_config)
    assert "sentry_dsn" in data["module"]["api-extensions"]
    assert "sentry_key" in data["resource"]
    assert "api-extensions" in data["resource"]["sentry_key"]
    sentry_data = data["resource"]["sentry_key"]["api-extensions"]
    assert "rate_limit_window" not in sentry_data
    assert "rate_limit_count" not in sentry_data

    comp_sentry = parsed_config.sites[0].components[0].sentry
    comp_sentry.rate_limit_window = 21600
    comp_sentry.rate_limit_count = 100

    data = _generate(parsed_config)
    sentry_data = data["resource"]["sentry_key"]["api-extensions"]
    assert sentry_data["rate_limit_window"] == 21600
    assert sentry_data["rate_limit_count"] == 100

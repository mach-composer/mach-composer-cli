import os
import re
from typing import Tuple

import pytest
from mach import terraform, types


@pytest.fixture
def tf_mock(mocker):
    return mocker.patch("mach.terraform.run_terraform")


def _generate(config: types.MachConfig) -> str:
    terraform.generate_terraform(config)
    file_path = os.path.join(config.output_path, config.sites[0].identifier, "site.tf")
    os.path.exists(file_path)
    with open(file_path) as f:
        return f.read()


def test_generate_terraform(parsed_config: types.MachConfig, tf_mock):
    content = _generate(parsed_config)
    tf_mock.assert_called_once()

    # Perform some very basic checks on the file
    # Later we'll convert this to more intelligent parsing and checking
    assert 'module "api-extensions"' in content
    assert 'backend "s3"' in content


def test_generate_w_sentry(parsed_config: types.MachConfig, tf_mock):
    content = _generate(parsed_config)
    assert "sentry" not in content

    parsed_config.components[0].integrations = ["aws", "sentry"]
    content = _generate(parsed_config)
    assert "sentry_dsn" in content
    assert "sentry_key" not in content

    parsed_config.general_config.sentry = types.SentryConfig(
        auth_token="12345",
        organization="labd",
        project="unittest",
    )
    content = _generate(parsed_config)
    assert "sentry_dsn" in content
    assert 'resource "sentry_key" "api-extensions"' in content
    assert "rate_limit_window" not in content
    assert "rate_limit_count" not in content

    comp_sentry = parsed_config.sites[0].components[0].sentry
    comp_sentry.rate_limit_window = 21600
    comp_sentry.rate_limit_count = 100

    content = _generate(parsed_config)
    assert _fetch_attr_line(content, "rate_limit_window") == "21600"
    assert _fetch_attr_line(content, "rate_limit_count") == "100"


def _fetch_attr_line(content: str, attr_name: str) -> Tuple[str, str]:
    result = re.search(rf"{attr_name}\s+=\s+(.*)", content)
    assert result, f"Attribute {attr_name} not found in content"
    return result.group(1)

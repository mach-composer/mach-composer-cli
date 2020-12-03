import os

import pytest
import yaml
from mach.commands import apply

from tests.utils import get_file


@pytest.fixture
def tf_mock(mocker):
    return mocker.patch("mach.terraform.run_terraform")


def test_endpoint_redeploys(click_runner, click_dir, tf_mock, tmp_path):
    orig_config = get_file("aws_config1.yml")

    with open(orig_config, "r") as fh:
        config_data = yaml.full_load(fh)

    result = click_runner.invoke(apply, ["-f", orig_config])
    assert result.exit_code == 0, result.stdout_bytes

    assert call_args(tf_mock) == [
        "fmt",  # mach-site-eu
        "fmt",  # mach-site-us
        "init",  # mach-site-eu
        ["apply"],  # mach-site-eu
        "init",  # mach-site-us
        ["apply"],  # mach-site-us
    ]
    tf_mock.reset_mock()

    new_config = os.path.join(tmp_path, "aws_config2.yml")
    config_data["sites"][0]["endpoints"] = {
        "main": {
            "url": "https://api.eu-tst.mach-example.net",
            "redeploy": True,
        }
    }
    content = yaml.dump(config_data, indent=2, explicit_start=True, sort_keys=False)
    with open(new_config, "w") as f:
        f.write(content)

    result = click_runner.invoke(apply, ["-f", new_config])
    assert result.exit_code == 0, result.stdout_bytes

    assert call_args(tf_mock) == [
        "fmt",  # mach-site-eu
        "fmt",  # mach-site-us
        "init",  # mach-site-eu
        ["taint", "aws_apigatewayv2_deployment.main_default"],  # mach-site-eu
        ["apply"],  # mach-site-eu
        "init",  # mach-site-us
        ["apply"],  # mach-site-us
    ]


def test_endpoint_redeploy_via_cmd(click_runner, click_dir, tf_mock):
    orig_config = get_file("aws_config1.yml")

    result = click_runner.invoke(apply, ["-f", orig_config])
    assert result.exit_code == 0, result.stdout_bytes

    assert call_args(tf_mock) == [
        "fmt",  # mach-site-eu
        "fmt",  # mach-site-us
        "init",  # mach-site-eu
        ["apply"],  # mach-site-eu
        "init",  # mach-site-us
        ["apply"],  # mach-site-us
    ]
    tf_mock.reset_mock()

    result = click_runner.invoke(
        apply, ["-f", orig_config, "--endpoint-redeploy", "main"]
    )
    assert result.exit_code == 0, result.stdout_bytes

    assert call_args(tf_mock) == [
        "fmt",  # mach-site-eu
        "fmt",  # mach-site-us
        "init",  # mach-site-eu
        ["taint", "aws_apigatewayv2_deployment.main_default"],  # mach-site-eu
        ["apply"],  # mach-site-eu
        "init",  # mach-site-us
        ["taint", "aws_apigatewayv2_deployment.main_default"],  # mach-site-us
        ["apply"],  # mach-site-us
    ]


def call_args(mock):
    result = []
    for arg in mock.call_args_list:
        result.append(arg[0][0])
    return result

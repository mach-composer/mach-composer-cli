import pytest
from mach.commands import apply

from tests.utils import get_file


@pytest.fixture
def tf_mock(mocker):
    return mocker.patch("mach.terraform.run_terraform")


def test_apply(click_runner, click_dir, tf_mock):
    orig_config = get_file("aws_config1.yml")

    result = click_runner.invoke(apply, ["-f", orig_config, "--ignore-version"])
    assert result.exit_code == 0, result.stdout_bytes
    assert call_args(tf_mock) == [
        "fmt",  # mach-site-eu
        "fmt",  # mach-site-us
        "init",  # mach-site-eu
        ["apply"],  # mach-site-eu
        "init",  # mach-site-us
        ["apply"],  # mach-site-us
    ]


def call_args(mock):
    result = []
    for arg in mock.call_args_list:
        result.append(arg[0][0])
    return result

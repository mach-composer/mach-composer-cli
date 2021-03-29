import os
from io import StringIO

import pytest
import yaml as yaml_
from mach.yaml import YamlIncludeConstructor

from tests.utils import get_file_content


@pytest.fixture
def yaml():
    YamlIncludeConstructor.add_to_loader_class()
    return yaml_


def test_include_local_resolve(yaml):
    data = yaml.load(StringIO("!include tests/files/components.yml"))
    assert data == [
        {
            "integrations": [""],
            "name": "commercetools-config",
            "source": "git::https://github.com/some-organisation/mach-component-commercetools.git//terraform",  # noqa
            "version": "1aa9215",
        },
        {
            "endpoints": {"public": "main"},
            "integrations": ["aws", "commercetools"],
            "name": "payment",
            "source": "git::https://github.com/some-organisation/mach-component-payment.git//terraform",  # noqa
            "version": "0a9a0b5",
        },
        {
            "endpoints": {"public": "default"},
            "integrations": ["aws", "commercetools"],
            "name": "us-payment",
            "source": "git::https://github.com/some-organisation/mach-component-payment.git//terraform",  # noqa
            "version": "0a9a0b5",
        },
        {
            "integrations": ["aws", "commercetools", "sentry"],
            "name": "api-extensions",
            "source": "git::https://github.com/some-organisation/mach-component-api-extensions.git//terraform",  # noqa
            "version": "a4bbb28",
        },
    ]


def test_include_http_resolve(yaml, requests_mock):
    mock = requests_mock.get(
        "https://machcomposer.io/components.yml",
        text=get_file_content("components.yml"),
    )
    data = yaml.load(StringIO("!include https://machcomposer.io/components.yml"))
    assert data == [
        {
            "integrations": [""],
            "name": "commercetools-config",
            "source": "git::https://github.com/some-organisation/mach-component-commercetools.git//terraform",  # noqa
            "version": "1aa9215",
        },
        {
            "endpoints": {"public": "main"},
            "integrations": ["aws", "commercetools"],
            "name": "payment",
            "source": "git::https://github.com/some-organisation/mach-component-payment.git//terraform",  # noqa
            "version": "0a9a0b5",
        },
        {
            "endpoints": {"public": "default"},
            "integrations": ["aws", "commercetools"],
            "name": "us-payment",
            "source": "git::https://github.com/some-organisation/mach-component-payment.git//terraform",  # noqa
            "version": "0a9a0b5",
        },
        {
            "integrations": ["aws", "commercetools", "sentry"],
            "name": "api-extensions",
            "source": "git::https://github.com/some-organisation/mach-component-api-extensions.git//terraform",  # noqa
            "version": "a4bbb28",
        },
    ]

    assert mock.call_count == 1


def test_include_git_resolve(yaml, mocker):
    def ensure_(repo, dest):
        os.mkdir(dest)
        with open(os.path.join(dest, "components.yml"), "w+") as f:
            f.write(get_file_content("components.yml"))

    mock = mocker.patch("mach.git.ensure_local", side_effect=ensure_)
    data = yaml.load(
        StringIO(
            "!include git::https://github.com/labd/mach-configurations.git@9f42fe2//components.yml"
        )
    )
    assert data == [
        {
            "integrations": [""],
            "name": "commercetools-config",
            "source": "git::https://github.com/some-organisation/mach-component-commercetools.git//terraform",  # noqa
            "version": "1aa9215",
        },
        {
            "endpoints": {"public": "main"},
            "integrations": ["aws", "commercetools"],
            "name": "payment",
            "source": "git::https://github.com/some-organisation/mach-component-payment.git//terraform",  # noqa
            "version": "0a9a0b5",
        },
        {
            "endpoints": {"public": "default"},
            "integrations": ["aws", "commercetools"],
            "name": "us-payment",
            "source": "git::https://github.com/some-organisation/mach-component-payment.git//terraform",  # noqa
            "version": "0a9a0b5",
        },
        {
            "integrations": ["aws", "commercetools", "sentry"],
            "name": "api-extensions",
            "source": "git::https://github.com/some-organisation/mach-component-api-extensions.git//terraform",  # noqa
            "version": "a4bbb28",
        },
    ]

    mock.assert_called_once()

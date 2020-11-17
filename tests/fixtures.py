import tempfile

import pytest
from click.testing import CliRunner
from mach import parse, types


@pytest.fixture
def config():
    return types.MachConfig(
        general_config=types.GeneralConfig(
            environment="test",
            terraform_config=types.TerraformConfig(
                aws_remote_state=types.AWSTFState(
                    bucket="unittest",
                    key_prefix="test",
                )
            ),
            cloud="aws",
            sentry=types.SentryConfig(dsn="sentry-dsn"),
        ),
        sites=[
            types.Site(
                identifier="unittest-nl",
                components=[
                    types.Component(
                        name="api-extensions",
                    )
                ],
            ),
        ],
        components=[
            types.ComponentConfig(
                name="api-extensions",
                source="some-source//terraform",
                version="1.0",
            )
        ],
        output_path=tempfile.gettempdir(),
    )


@pytest.fixture
def parsed_config(config):
    return parse.parse_config(config)


@pytest.fixture
def click_runner():
    return CliRunner()


@pytest.fixture
def click_dir(click_runner):
    with click_runner.isolated_filesystem() as f:
        yield f


class CookiecutterCache:
    _result = None

    def determine_repo_dir(self, *args, **kwargs):
        if self._result:
            return self._result

        from cookiecutter.repository import determine_repo_dir as _determine_repo_dir

        self._result = _determine_repo_dir(*args, **kwargs)
        return self._result


cc_cache = CookiecutterCache()


@pytest.fixture()
def cookiecutter(mocker):
    return mocker.patch(
        "cookiecutter.main.determine_repo_dir", side_effect=cc_cache.determine_repo_dir
    )

import pytest
import requests_mock as rmock
from click.testing import CliRunner


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


@pytest.fixture()
def requests_mock():
    with rmock.Mocker() as m:
        yield m

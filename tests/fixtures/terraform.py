import pytest


@pytest.fixture
def tf_mock(mocker):
    return mocker.patch("mach.terraform.run_terraform")

import pytest
from mach import types


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
    )

import tempfile

import pytest
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
                    region="eu-central-1",
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
                aws=types.SiteAWSSettings(
                    account_id=1234567890,
                    region="eu-central-1",
                ),
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
def azure_config():
    return types.MachConfig(
        general_config=types.GeneralConfig(
            environment="test",
            terraform_config=types.TerraformConfig(
                azure_remote_state=types.AzureTFState(
                    resource_group="shared-rg",
                    storage_account="machsaterra",
                    container_name="tfstate",
                    state_folder="test",
                )
            ),
            azure=types.AzureConfig(
                tenant_id="6f10659d-4227-43e6-95ab-80d12a18acf9",
                subscription_id="5f34d95d-4dd8-40b3-9d18-f9007e2ce6ac",
                region="westeurope",
            ),
            cloud="azure",
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
                short_name="apiexts",
                version="1.0",
            ),
            types.ComponentConfig(
                name="product-types",
                source="product-types//terraform",
                version="v0.1.0",
                integrations=[""],
            )
        ],
        output_path=tempfile.gettempdir(),
    )


@pytest.fixture
def parsed_config(config):
    return parse.parse_config(config)


@pytest.fixture
def parsed_azure_config(azure_config):
    return parse.parse_config(azure_config)

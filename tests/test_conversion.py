from mach.types import MachConfig


def test_short_name_is_set():
    example_config = {
        "general_config": {
            "environment": "test",
            "terraform_config": {
                "azure_remote_state": {
                    "resource_group_name": "my-shared-rg",
                    "storage_account_name": "mysharedsaterra",
                    "container_name": "tfstate",
                    "state_folder": "dev",
                }
            },
            "cloud": "azure",
            "sentry": {"dsn": "some sentry url"},
        },
        "sites": [],
        "components": [
            {"name": "example", "source": "some source", "version": "HEAD"},
            {
                "name": "api_extensions",
                "short_name": "apiexts",
                "source": "some source",
                "version": "HEAD",
            },
        ],
    }

    config = MachConfig.schema().load(example_config)

    assert config.components[0].short_name == "example"
    assert config.components[1].short_name == "apiexts"

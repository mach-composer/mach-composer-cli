import json

from mach import parse, types

from tests.unittests.terraform import utils as tf


def test_generate_terraform(parsed_config: types.MachConfig, tf_mock):
    data = tf.generate(parsed_config)
    tf_mock.assert_called_once()

    assert data == {
        "terraform": [
            {
                "backend": [
                    {
                        "s3": {
                            "bucket": ["unittest"],
                            "key": ["test/unittest-nl"],
                            "region": ["eu-central-1"],
                            "encrypt": [True],
                        }
                    }
                ]
            },
            {"required_providers": [{"aws": [{"version": "~> 3.74.1"}]}]},
        ],
        "provider": [{"aws": {"region": ["eu-central-1"]}}],
        "locals": [{"tags": [{"Environment": "test", "Site": "unittest-nl"}]}],
        "module": [
            {
                "api-extensions": {
                    "source": ["some-source//terraform"],
                    "component_version": ["1.0"],
                    "environment": ["test"],
                    "site": ["unittest-nl"],
                    "tags": ["${local.tags}"],
                    "variables": [{}],
                    "secrets": [{}],
                    "providers": [{"aws": "${aws}"}],
                    "depends_on": [[]],
                }
            }
        ],
    }


def test_generate_w_sentry(parsed_config: types.MachConfig, tf_mock):
    data = tf.generate(parsed_config)
    data_str = json.dumps(data)
    assert "sentry" not in data_str

    parsed_config.components[0].integrations = ["aws", "sentry"]
    data = tf.generate(parsed_config)

    assert "sentry_dsn" in data.module["api-extensions"]
    assert "sentry_key" not in data.get("resource", {})

    parsed_config.general_config.sentry = types.SentryConfig(
        auth_token="12345",
        organization="labd",
        project="unittest",
    )
    data = tf.generate(parsed_config)
    assert "sentry_dsn" in data.module["api-extensions"]
    assert "sentry_key" in data.resource
    assert "api-extensions" in data.resource.sentry_key
    sentry_data = data.resource.sentry_key["api-extensions"]
    assert "rate_limit_window" not in sentry_data
    assert "rate_limit_count" not in sentry_data

    comp_sentry = parsed_config.sites[0].components[0].sentry
    comp_sentry.rate_limit_window = 21600
    comp_sentry.rate_limit_count = 100

    data = tf.generate(parsed_config)
    sentry_data = data.resource.sentry_key["api-extensions"]
    assert sentry_data["rate_limit_window"] == 21600
    assert sentry_data["rate_limit_count"] == 100


def test_generate_w_stores(config: types.MachConfig, tf_mock):
    config.sites[0].commercetools = types.CommercetoolsSettings(
        project_key="ct-unit-test",
        client_id="a96e59be-24da-4f41-a6cf-d61d7b6e1766",
        client_secret="98c32de8-1a6c-45a9-a718-d3cce5201799",
        scopes="manage_project:ct-unit-test",
        project_settings=types.CommercetoolsProjectSettings(
            languages=["nl-NL"],
            countries=["NL"],
            currencies=["EUR"],
        ),
        stores=[
            types.CommercetoolsStore(
                name={
                    "en-GB": "Default store",
                },
                key="main-store",
            ),
            types.CommercetoolsStore(
                name={
                    "en-GB": "Some other store",
                },
                key="other-store",
            ),
            types.CommercetoolsStore(
                name={
                    "en-GB": "Forgotten store",
                },
                key="forgotten-store",
            ),
        ],
    )
    config.components[0].integrations = ["aws", "commercetools"]
    data = tf.generate(parse.parse_config(config))

    assert len(data.resource.commercetools_store) == 3
    assert "main-store" in data.resource.commercetools_store
    assert "other-store" in data.resource.commercetools_store
    assert "forgotten-store" in data.resource.commercetools_store

    assert len(data.module["api-extensions"].ct_stores) == 3

    for store_key, store in data.module["api-extensions"].ct_stores.items():
        assert store["key"] == store_key
        assert not store["variables"]
        assert not store["secrets"]

    config.sites[0].components[0].store_variables = {
        "main-store": {
            "FOO": "BAR",
            "EXTRA": "VALUES",
        },
        "other-store": {
            "FOO": "SOMETHING ELSE",
        },
    }
    config.sites[0].components[0].store_secrets = {
        "main-store": {
            "PAYMENT_KEY": "TLrlDf6XhKkXFGGHeQGY",
        },
    }

    data = tf.generate(parse.parse_config(config))
    main_store = data.module["api-extensions"].ct_stores["main-store"]
    other_store = data.module["api-extensions"].ct_stores["other-store"]
    assert len(main_store.variables) == 2
    assert len(other_store.variables) == 1
    assert len(main_store.secrets) == 1
    assert not other_store.secrets


def test_commercetools_frontend_credentials(config: types.MachConfig, tf_mock):
    config.sites[0].commercetools = types.CommercetoolsSettings(
        project_key="ct-unit-test",
        client_id="a96e59be-24da-4f41-a6cf-d61d7b6e1766",
        client_secret="98c32de8-1a6c-45a9-a718-d3cce5201799",
        scopes="manage_project:ct-unit-test",
        project_settings=types.CommercetoolsProjectSettings(
            languages=["nl-NL"],
            countries=["NL"],
            currencies=["EUR"],
        ),
        stores=[
            types.CommercetoolsStore(
                name={
                    "en-GB": "Default store",
                },
                key="main-store",
            ),
            types.CommercetoolsStore(
                name={
                    "en-GB": "Some other store",
                },
                key="other-store",
            ),
            types.CommercetoolsStore(
                name={
                    "en-GB": "Forgotten store",
                },
                key="forgotten-store",
            ),
        ],
    )
    data = tf.generate(parse.parse_config(config))
    assert list(data.resource.commercetools_api_client.keys()) == [
        "frontend_credentials_main-store",
        "frontend_credentials_other-store",
        "frontend_credentials_forgotten-store",
    ]

    assert data.resource.commercetools_api_client[
        "frontend_credentials_main-store"
    ].scope == [
        "create_anonymous_token:ct-unit-test",
        "manage_my_profile:ct-unit-test",
        "manage_my_profile:ct-unit-test:main-store",
        "manage_my_orders:ct-unit-test",
        "manage_my_orders:ct-unit-test:main-store",
        "manage_my_shopping_lists:ct-unit-test",
        "manage_my_payments:ct-unit-test",
        "view_products:ct-unit-test",
        "view_project_settings:ct-unit-test",
    ]

    config.sites[0].commercetools.frontend.create_credentials = False
    data = tf.generate(parse.parse_config(config))
    assert "commercetools_api_client" not in data.resource

    config.sites[0].commercetools.frontend.create_credentials = True
    config.sites[0].commercetools.frontend.permission_scopes = [
        "manage_my_profile",
        "manage_my_orders",
        "view_products",
        "manage_my_payments",
        "create_anonymous_token",
        "view_stores",
    ]

    data = tf.generate(parse.parse_config(config))
    assert data.resource.commercetools_api_client[
        "frontend_credentials_main-store"
    ].scope == [
        "manage_my_profile:ct-unit-test",
        "manage_my_profile:ct-unit-test:main-store",
        "manage_my_orders:ct-unit-test",
        "manage_my_orders:ct-unit-test:main-store",
        "view_products:ct-unit-test",
        "manage_my_payments:ct-unit-test",
        "create_anonymous_token:ct-unit-test",
        "view_stores:ct-unit-test",
    ]

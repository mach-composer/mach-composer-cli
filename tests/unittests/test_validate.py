import pytest
from mach import parse, types, validate
from mach.exceptions import ValidationError


def test_validate_sentry():
    with pytest.raises(ValidationError):
        validate.validate_sentry_config(types.SentryConfig())

    validate.validate_sentry_config(types.SentryConfig(dsn="12345"))

    with pytest.raises(ValidationError):
        validate.validate_sentry_config(types.SentryConfig(auth_token="12345"))

    with pytest.raises(ValidationError):
        validate.validate_sentry_config(
            types.SentryConfig(auth_token="12345", dsn="12345")
        )

    validate.validate_sentry_config(
        types.SentryConfig(
            auth_token="12345", project="my-project", organization="my-organization"
        )
    )


def test_validate_aws_default_endpoint(config: types.MachConfig):
    """It must be possible for a component to use the default API Gateway endpoint."""
    config.components[0].endpoints = {
        "main": "public",
    }
    config = parse.parse_config(config)

    with pytest.raises(ValidationError) as e:
        validate.validate_config(config)

    assert str(e.value) == "Missing required endpoints public"

    config.components[0].endpoints = {
        "main": "default",
    }
    config = parse.parse_config(config)
    validate.validate_config(config)


def test_validate_azure_default_endpoint(azure_config: types.MachConfig):
    """It must be possible for a component to use the default Frontdoor endpoint."""
    config = azure_config
    config.components[0].endpoints = {
        "main": "public",
    }
    config = parse.parse_config(config)

    with pytest.raises(ValidationError) as e:
        validate.validate_config(config)

    assert str(e.value) == "Missing required endpoints public"

    config.components[0].endpoints = {
        "main": "default",
    }
    config = parse.parse_config(config)
    validate.validate_config(config)


def test_validate_component_endpoint(config: types.MachConfig):
    """An endpoint defined on a component must exist for all sites that use that component."""
    config.components[0].endpoints = {
        "main": "public",
    }
    config = parse.parse_config(config)

    with pytest.raises(ValidationError):
        validate.validate_config(config)

    config.sites[0].endpoints = [
        types.Endpoint(
            key="public",
            url="api.mach-example.com",
        ),
        types.Endpoint(
            key="services",
            url="services.mach-example.com",
        ),
    ]
    validate.validate_config(config)

    new_site = types.Site(
        identifier="unittest-nl",
        components=[],
        aws=types.SiteAWSSettings(
            account_id=1234567890,
            region="eu-central-1",
        ),
    )
    config.sites.append(new_site)
    validate.validate_config(config)

    new_site.components.append(
        types.Component(
            name="api-extensions",
        )
    )
    config = parse.parse_config(config)

    with pytest.raises(ValidationError):
        validate.validate_config(config)


@pytest.mark.parametrize(
    "name, valid",
    (
        ("m", False),
        ("main store", False),
        ("main-store", True),
        ("main_store", True),
    ),
)
def test_validate_store_keys(name, valid):
    ct = types.CommercetoolsSettings(
        project_key="ct-unit-test",
        client_id="a96e59be-24da-4f41-a6cf-d61d7b6e1766",
        client_secret="98c32de8-1a6c-45a9-a718-d3cce5201799",
        scopes="manage_project:ct-unit-test",
        languages=["nl-NL"],
        countries=["NL"],
        currencies=["EUR"],
        stores=[
            types.Store(
                key=name,
                name={
                    "en-GB": "Default store",
                },
            ),
        ],
    )

    if not valid:
        with pytest.raises(ValidationError):
            validate.validate_store_keys(ct)
    else:
        validate.validate_store_keys(ct)

    ct.stores[0].key = "main-store"
    validate.validate_store_keys(ct)

    ct.stores.append(
        types.Store(
            key="main-store",
            name={
                "en-GB": "Another store",
            },
        ),
    )

    with pytest.raises(ValidationError):
        # Duplicate key
        validate.validate_store_keys(ct)

    ct.stores[1].key = "other-store"
    validate.validate_store_keys(ct)


def test_validate_stores(parsed_config: types.MachConfig):
    """Tests if the stores used in the store variables match the defined commercetools stores."""
    config = parsed_config
    site = config.sites[0]

    site.components[0].store_variables = {
        "main-store": {
            "FOO": "BAR",
        }
    }

    # It should fail because we refer a store that hasnt been defined yet
    with pytest.raises(ValidationError):
        validate.validate_config(config)

    site.commercetools = types.CommercetoolsSettings(
        project_key="ct-unit-test",
        client_id="a96e59be-24da-4f41-a6cf-d61d7b6e1766",
        client_secret="98c32de8-1a6c-45a9-a718-d3cce5201799",
        scopes="manage_project:ct-unit-test",
        languages=["nl-NL"],
        countries=["NL"],
        currencies=["EUR"],
        stores=[
            types.Store(
                name={
                    "en-GB": "Default store",
                },
                key="main-store",
            ),
        ],
    )

    validate.validate_config(config)


def test_validate_azure_service_plans(parsed_azure_config: types.MachConfig):
    config = parsed_azure_config
    config.components[0].azure.service_plan = "premium"

    # It should fail because we refer a store that hasnt been defined yet
    with pytest.raises(ValidationError) as e:
        validate.validate_config(config)
    assert str(e.value) == (
        "Component api-extensions requires service plan premium which is not defined in the "
        "Azure configuration."
    )

    config.general_config.azure.service_plans["premium"] = types.ServicePlan(
        kind="Linux",
        tier="PremiumV2",
        size="P2v2",
    )

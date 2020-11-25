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


def test_validate_endpoints(parsed_config: types.MachConfig):
    config = parsed_config

    config.sites[0].endpoints = {
        "public": "api.mach-example.com",
        "services": "services.mach-example.com",
    }

    with pytest.raises(ValidationError):
        validate.validate_config(config)

    config.sites[0].aws.route53_zone_name = "mach-example.com"
    validate.validate_config(config)

    # Change one of the components that does not match the DNS zone anymore
    config.sites[0].endpoints["services"] = "api.mach-services.com"
    with pytest.raises(ValidationError):
        validate.validate_config(config)


def test_validate_component_endpoint(parsed_config: types.MachConfig):
    """An endpoint defined on a component must exist for all sites that use that component."""
    config = parsed_config

    config.components[0].endpoint = "public"

    with pytest.raises(ValidationError):
        validate.validate_config(config)

    config.sites[0].endpoints = {
        "public": "api.mach-example.com",
        "services": "services.mach-example.com",
    }
    config.sites[0].aws.route53_zone_name = "mach-example.com"
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

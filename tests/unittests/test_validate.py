import pytest
from mach import types, validate
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


def test_validate_endpoints(config: types.MachConfig):
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

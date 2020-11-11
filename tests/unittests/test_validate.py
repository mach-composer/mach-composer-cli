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

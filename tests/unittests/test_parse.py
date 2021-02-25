from mach import parse, types


def test_resolve_sentry_configs(config: types.MachConfig):
    """Check if sentry configurations get merged correctly."""
    # Sanity check
    assert config.general_config.sentry and config.general_config.sentry.dsn
    sentry_dsn = config.general_config.sentry.dsn

    # Initially, not site or component specific sentry config should be set
    assert not config.sites[0].sentry
    assert not config.sites[0].components[0].sentry

    parse.resolve_site_configs(config)

    # In order for a correct Terraform file, the parser ensures
    # the basic sentry settings are passed on to the lowest level (the component)
    assert config.sites[0].sentry.dsn == sentry_dsn
    assert config.sites[0].components[0].sentry.dsn == sentry_dsn

    # Reset component sentry and add some extra settings to the site
    config.sites[0].components[0].sentry = None
    config.sites[0].sentry.rate_limit_window = 500
    config.sites[0].sentry.rate_limit_count = 100

    parse.resolve_site_configs(config)
    comp_sentry = config.sites[0].components[0].sentry
    assert comp_sentry.dsn == sentry_dsn
    assert comp_sentry.rate_limit_window == 500
    assert comp_sentry.rate_limit_count == 100

    # No set only one attribute on the component sentry
    # The site specific settings should be merged
    config.sites[0].components[0].sentry = types.SentryDsn(rate_limit_count=50)
    parse.resolve_site_configs(config)
    comp_sentry = config.sites[0].components[0].sentry
    assert comp_sentry.dsn == sentry_dsn
    assert comp_sentry.rate_limit_window == 500
    assert comp_sentry.rate_limit_count == 50


def test_parse_endpoints(config: types.MachConfig):
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

    config = parse.parse_config(config)
    for endpoint in config.sites[0].endpoints:
        assert endpoint.zone == "mach-example.com"


def test_parse_azure_service_plans(azure_config: types.MachConfig):
    config = azure_config
    # Sanity check
    assert not config.general_config.azure.service_plans
    for c in config.components:
        assert not (c.azure and c.azure.service_plan)
    for c in config.sites[0].components:
        assert not (c.azure and c.azure.service_plan)

    config = parse.parse_config(config)
    assert "default" in config.general_config.azure.service_plans
    assert config.general_config.azure.service_plans["default"] == types.ServicePlan(
        kind="FunctionApp", tier="Dynamic", size="Y1"
    )
    for c in config.components:
        if "azure" in c.integrations:
            assert c.azure and c.azure.service_plan == "default"
        else:
            assert not c.azure or not c.azure.service_plan

    for c in config.sites[0].components:
        if "azure" in c.integrations:
            assert c.azure and c.azure.service_plan == "default"
        else:
            assert not c.azure or not c.azure.service_plan


def test_apollo_federation_integration_set(apollo_config: types.MachConfig):
    config = apollo_config
    for c in config.components:
        assert "apollo_federation" in c.integrations

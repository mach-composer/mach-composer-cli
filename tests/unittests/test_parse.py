from mach import parse, types


def test_resolve_site_configs(config: types.MachConfig):
    assert config.general_config.sentry and config.general_config.sentry.dsn
    sentry_dsn = config.general_config.sentry.dsn

    assert not config.sites[0].sentry_dsn
    assert not config.sites[0].components[0].sentry_dsn

    parse.resolve_site_configs(config)

    assert config.sites[0].sentry_dsn == sentry_dsn
    assert config.sites[0].components[0].sentry_dsn == sentry_dsn

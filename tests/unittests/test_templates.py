import pytest
from mach import exceptions, templates


@pytest.mark.parametrize(
    "url, zone",
    (
        ("https://api.labd.io", "labd.io"),
        ("https://api.test.mach-examples.net", "test.mach-examples.net"),
        ("api.test.mach-examples.net", "test.mach-examples.net"),
        ("http://api.test.mach-examples.net", "test.mach-examples.net"),
    ),
)
def test_zone_name_filter(url, zone):
    assert templates.zone_name(url) == zone


def test_render_tfvalue():
    assert templates.render_tfvalue(True) == "true"
    assert templates.render_tfvalue(False) == "false"

    assert templates.render_tfvalue(12) == 12
    assert (
        templates.render_tfvalue(["value1", False, "value2"])
        == """["value1",false,"value2"]"""
    )

    assert (
        templates.render_tfvalue(r"${component.infra.db_password}")
        == "module.infra.db_password"
    )


def test_render_config_variable():
    assert (
        templates.parse_config_variable(r"${component.infra.db_password}")
        == "module.infra.db_password"
    )
    assert templates.parse_config_variable(r"${component.infra.db_password") is None

    with pytest.raises(exceptions.MachError):
        assert templates.parse_config_variable(r"${component.infra}")

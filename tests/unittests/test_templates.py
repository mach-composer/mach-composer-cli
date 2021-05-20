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


def test_render_variable():
    assert templates.render_variable(True) == "true"
    assert templates.render_variable(False) == "false"

    assert templates.render_variable(12) == 12
    assert (
        templates.render_variable(["value1", False, "value2"])
        == """["value1",false,"value2"]"""
    )

    assert (
        templates.render_variable(r"${component.infra.db_password}")
        == "module.infra.db_password"
    )


def test_render_config_variable():
    assert (
        templates.parse_config_variable(r"${component.infra.db_password}")
        == "module.infra.db_password"
    )
    assert templates.parse_config_variable(r"${component.infra.db_password") == None

    with pytest.raises(exceptions.MachError):
        assert templates.parse_config_variable(r"${component.infra}")

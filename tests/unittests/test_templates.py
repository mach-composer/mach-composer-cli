import pytest
from mach import templates


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

import pytest
from mach import utils


@pytest.mark.parametrize(
    "url, zone",
    (
        ("https://mach.example.com", "example.com"),
        ("api.test.example.com", "test.example.com"),
    ),
)
def test_dns_zone_from_url(url, zone):
    assert utils.dns_zone_from_url(url) == zone

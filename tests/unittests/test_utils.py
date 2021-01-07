import pytest
from mach import utils


@pytest.mark.parametrize(
    "url, zone",
    (
        ("https://mach.example.com", "example.com"),
        ("https://mach.example.co.uk", "example.co.uk"),
        ("api.test.example.com", "test.example.com"),
        ("api.test.example.co.uk", "test.example.co.uk"),
    ),
)
def test_dns_zone_from_url(url, zone):
    assert utils.dns_zone_from_url(url) == zone


@pytest.mark.parametrize(
    "url",
    (
        "https://example.com",
        "https://example.co.uk",
    ),
)
def test_dns_zone_from_url_invalid(url):
    with pytest.raises(ValueError):
        utils.dns_zone_from_url(url)

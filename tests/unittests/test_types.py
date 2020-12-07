from copy import deepcopy

import pytest
from mach import parse, types


def test_site(config: types.MachConfig):
    config.components.append(
        types.ComponentConfig(
            name="private-api",
            source="some-private-source/terraform",
            version="1.0",
        )
    )
    config.sites[0].components.append(types.Component(name="private-api"))

    config.sites[0].endpoints = [
        types.Endpoint(
            key="public",
            url="api.example.com",
        ),
        types.Endpoint(
            key="private",
            url="private-api.example.com",
        ),
    ]
    config = parse.parse_config(config)
    site = config.sites[0]

    assert site.used_endpoints == []
    config.components[0].endpoint = "public"
    config = parse.parse_config(config)
    assert site.used_endpoints == [
        types.Endpoint(
            key="public", url="api.example.com", components=[site.components[0]]
        )
    ]

    config.components[1].endpoint = "public"
    config = parse.parse_config(config)
    assert site.used_endpoints == [
        types.Endpoint(
            key="public",
            url="api.example.com",
            components=[site.components[0], site.components[1]],
        )
    ]

    config.components[1].endpoint = "private"
    config = parse.parse_config(config)
    assert site.used_endpoints == [
        types.Endpoint(
            key="public",
            url="api.example.com",
            components=[
                site.components[0],
            ],
        ),
        types.Endpoint(
            key="private",
            url="private-api.example.com",
            components=[site.components[1]],
        ),
    ]


def test_hybrid_endpoints():
    endpoints_flat = {
        "main": "api.example.com",
        "private": "private.example.com",
    }
    endpoints_complex = {
        "main": {
            "url": "api.example.com",
        },
        "private": {
            "url": "private.example.com",
        },
    }

    site_schema = types.Site.schema(infer_missing=True)

    for input_ in [endpoints_flat, endpoints_complex]:
        site = site_schema.load({"identifier": "nl-unittest", "endpoints": input_})

        assert site.endpoints == [
            types.Endpoint(
                key="main",
                url="api.example.com",
            ),
            types.Endpoint(
                key="private",
                url="private.example.com",
            ),
        ]

        serialized = site_schema.dump(site)
        assert serialized["endpoints"] == endpoints_flat


def test_hybrid_endpoints_wrong_value():
    with pytest.raises(Exception):
        types.Site.schema(infer_missing=True).load(
            {"identifier": "nl-unittest", "endpoints": ["bla", "bla"]}
        )

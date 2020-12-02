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

    config.sites[0].endpoints = {
        "public": "api.example.com",
        "private": "private-api.example.com",
    }
    config = parse.parse_config(config)
    site = config.sites[0]

    assert site.used_endpoints == []

    config.components[0].endpoint = "public"
    assert site.used_endpoints == [
        types.Endpoint(
            key="public", url="api.example.com", components=[site.components[0]]
        )
    ]

    config.components[1].endpoint = "public"
    assert site.used_endpoints == [
        types.Endpoint(
            key="public",
            url="api.example.com",
            components=[site.components[0], site.components[1]],
        )
    ]

    config.components[1].endpoint = "private"
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

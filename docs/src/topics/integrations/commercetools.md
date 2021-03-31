# commercetools

## Configuration

From a MACH configuration file, you can configure the following items in commercetools:

- currencies
- languages
- countries
- shipping zones
- channels
- taxes
- stores

For more information about these configuration options, see the [syntax](../../reference/syntax/sites.md#commercetools).

!!! tip "More fine-grained control"
    MACH provides a couple of basic configuration options which in most cases are sufficient and is setup very quickly.<br>
    If you need more control over the configuration in your project, you can load custom configuration through a new component which would include the necessary Terraform configuration for that.<br>
    For example a `commercetools-setup-component`.

## API client management

MACH is provided by admin credentials which allows MACH to fully manage your commercetools project.

Each component that implements functionality that needs to communicate with commercetools needs their own set of credentials.
The component is responsible for creating the necessary API client credentials for that specific component, with only the most necessary scopes set.

One way a component could facilitate in this would be to define the following in your Terraform config:

```terraform
locals {
  ct_scopes = formatlist("%s:%s", [
    "manage_orders",
  ], var.ct_project_key)
}

resource "commercetools_api_client" "client" {
  name  = "api-extension"
  scope = local.ct_scopes
}
```

## Stores

MACH provides built-in [commercetools Stores](https://docs.commercetools.com/api/projects/stores) support.

It's possible to

- Configure your stores within a MACH configuration (see [example](#example-configuration-block) below)
- Define store-specific variables and secrets per component

For example, when two stores are defined `uk-store` and `nl-store`, these can be used to configure your component:

```yaml
component:
  - my-component:
    variables:
      FROM_EMAIL: mach@example.com
    store_variables:
      uk-store:
        FROM_EMAIL: mach@example.co.uk
      nl-store:
        FROM_EMAIL: mach@example.nl
    store_secrets:
      uk-store:
        MAIL_API_KEY: 5kg6Z4HxgLMfBjYj5BOg
      nl-store:
        MAIL_API_KEY: i1hajIJ92LPNYGB2p3W1
```

The component can use the correct variables based on the context that is decided based on incoming requests or other information available.

!!! warning "Using store-aware context in components"
    Although MACH composer provides a data structure for managing your configuration in different contexts, it is the responsibility of the component itself to parse this data structure correctly, and apply the right configuration in the right context.

    Read the [parse store variables and secrets](../../howto/commercetools/store-vars.md) how-to for more information.

### Outside MACH configuration

If you don't want to manage the stores from the MACH configuration itself, you still have the possibility to define them in the configuration so that components can use store variables and receive a full list of the stores available.

This is done by setting a store to `managed: false`:

```yaml
commercetools:
  stores:
    - key: uk-store
      managed: false
    - key: nl-store
      managed: false
```

## Example configuration block

An example [commercetools configuration block](../../reference/syntax/sites.md#commercetools):

```yaml
commercetools:
  project_key: nl-mach-tst
  client_id: ixwVv7T3rnmJ
  client_secret: dElymMZvDqW3X5RLASMS
  scopes: manage_project:nl-mach-tst manage_api_clients:nl-mach-tst view_api_clients:nl-mach-tst
  languages:
      - en-GB
      - nl-NL
  currencies:
      - GBP
      - EUR
  countries:
      - GB
      - IE
      - NL
  channels:
    - key: INV
      roles:
        - InventorySupply
      name:
        en-GB: Inventory
      description:
        en-GB: Main inventory channel
    - key: DIST-EUR
      roles:
        - ProductDistribution
      name:
        en-GB: Europe Distribution
      description:
        en-GB: Europe distribution channel
  stores:
      - key: uk-store
        name:
          en: UK store
        languages:
          - en-GB
        distribution_channels:
          - DIST-EUR
        inventory_channels:
          - INV
      - key: nl-store
        name:
          en: NL store
        languages:
          - nl-NL
        distribution_channels:
          - DIST-EUR
        inventory_channels:
          - INV
```

## Integrate with components

When `commercetools` is set as an [component integration](../../reference/components/structure.md#integrations), the component should have the following Terraform variables defined:

- `ct_project_key`
- `ct_api_url`
- `ct_auth_url`
- `ct_stores`

!!! info ""
    More information on the [commercetools integration on components](../../reference/components/structure.md#commercetools).

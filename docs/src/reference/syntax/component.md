# Component

All component definitions in this configuration file. Components in this 
context define deployables that are available to roll out per site. These 
components are used and configured separately [per site](./site.md#nested-schema-for-components).

## Example

```yaml
components:
  - name: api-extension-products
    source: git::ssh://git@git.labdigital.nl/mach-components/api-extensions-component.git//products/terraform
    version: 3b8ab91
    paths:
      - products/terraform
  - name: api-extension-orders
    source: git::ssh://git@git.labdigital.nl/mach-components/api-extensions-component.git//orders/terraform
    version: 3b8ab91
    paths:
      - orders/terraform
  - name: ct-products-types
    source: git::ssh://git@git.labdigital.nl/mach-components/ct-product-types.git//terraform
    version: 1.4.0
    integrations:
      - commercetools
```

## Schema

### Required

- `name` (String) Name of the component. To be used as reference in the site
  definitions.
- `version` (String) A Git commit hash or tag
- `source` (String) Source definition of the terraform module

### Optional

- `integrations` (List of String) Defines a list of integrations for the given
  component. It controls what Terraform variables are passed on to the
  components [Terraform module](../../concepts/components/structure.md).
  Defaults to the cloud provider of the site as specified in
  the [global config](global.md).
- `endpoints` (Map of String, _deprecated_) Defines the endpoint that needs to
  connect to this component. Will set up Frontdoor routing or pass API Gateway
  information when set.
- `health_check_path` (String, _deprecated_) Defines a custom healthcheck path.
  Defaults to `/<name>/healthchecks`
- `paths`(String) Defines a list of paths which contain source code for a
  component. This is most useful when working with a monorepo, as it allows for
  filtering for updates (for example with `mach-composer update --check`).
  Default is empty, which will assume all changes are relevant
- `branch` (String) Configure the git branch of the component. If left empty
  `main` will be used. Only used to facilitate the `mach-composer update`
  CLI command.

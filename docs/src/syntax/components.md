# components definitions

All component definitions in this configuration file.<br>
These components are used and configured separately [per site](./sites.md#component-configurations).

Example:

```yaml
components:
  - name: api-extensions
    short_name: apiexts
    source: git::ssh://git@git.labdigital.nl/mach-components/api-extensions-component.git//terraform
    version: 3b8ab91
    endpoints: 
      main: default
  - name: ct-products-types
    source: git::ssh://git@git.labdigital.nl/mach-components/ct-product-types.git//terraform
    version: 1.4.0
    integrations: ["commercetools"]
```

- `name` - Name of the component. To be used as reference in the site definitions.
- `version` - A Git commit hash or tag
- `source` - Source definition of the terraform module
- `short_name` - Short name to be used in cloud resources. Should be at most 10 characters to avoid running into Resource naming limits.<br>
  Defaults to the given components `name`
- `integrations` - Defines a list of integrations for the given component. It controls what Terraform variables are passed on to the components [Terraform module](../components/structure.md#terraform-module).<br>
  Defaults to `["azure"]` or `["aws"]`, depending on your cloud provider.<br>
  Could be any of:
    - `azure`
    - `aws`
    - `commercetools`
    - `contentful`<br>
- `endpoints` - Defines the endpoint that needs to connect to this component.<br>
  Will setup Frontdoor routing or pass API Gateway information when set.
- `health_check_path` - Defines a custom healthcheck path.<br>
  Defaults to `/<name>/healthchecks`
- `branch` - Configure the git branch of the component. Only used to facilitate the `mach update` CLI command.
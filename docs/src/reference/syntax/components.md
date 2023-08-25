# components

All component definitions in this configuration file.<br>
These components are used and configured separately [per site](./sites.md#components).

Example:

```yaml
components:
  - name: api-extension-products
    source: git::ssh://git@git.labdigital.nl/mach-components/api-extensions-component.git//products/terraform
    version: 3b8ab91
    paths:
      - products/terraform
    endpoints:
      main: default
  - name: api-extension-orders
    source: git::ssh://git@git.labdigital.nl/mach-components/api-extensions-component.git//orders/terraform
    version: 3b8ab91
    paths:
      - orders/terraform
    endpoints:
      main: default
  - name: ct-products-types
    source: git::ssh://git@git.labdigital.nl/mach-components/ct-product-types.git//terraform
    version: 1.4.0
    integrations: [ "commercetools" ]
```

- `name` - Name of the component. To be used as reference in the site definitions.
- `version` - A Git commit hash or tag
- `source` - Source definition of the terraform module
- `integrations` - Defines a list of integrations for the given component. It controls what Terraform variables are passed on
  to the components [Terraform module](../components/structure.md#terraform-module).<br>

  Defaults to `["azure"]` or `["aws"]`, depending on your cloud provider.<br>
  Could be any of:
    - `azure`
    - `aws`
    - `commercetools`
    - `contentful`<br>

- `endpoints` - (_deprecated_) Defines the endpoint that needs to connect to this component.<br>
  Will setup Frontdoor routing or pass API Gateway information when set.
- `health_check_path` - Defines a custom healthcheck path. Defaults to `/<name>/healthchecks`
- `paths`- Defines a list of paths which contain source code for a component. This is most useful when working with a
  monorepo, as it allows for filtering for updates (for example with `mach-composer update --check`). Default is empty, 
  which will assume all changes are relevant

!!! tip "Development settings"
In addition to the default set of component settings, a couple of settings
can be defined during development.<br> These are not intended to be used
for a production deployment, but can facilitate local development:

    - `branch` - Configure the git branch of the component. Only used to
      facilitate the `mach-composer update` CLI command.
    - `artifacts` - Mapping of additional artifacts **AWS only**
      - `script` - Script file to build and package the component, relative to the workdir.
      - `filename` - Filename to be used for deployment, relative to the workdir
      - `workdir` - Work directory for the script/filename, relative to current work dir (default)

    Example:
    ```yaml
    components:
    - name: my-component
      source: ../mach-component-my-component/terraform
      integrations: ["aws", "commercetools", "sentry"]
      version: latest
      artifacts:
        service:
          filename: .serverless/my-component.zip
          script: yarn package
          workdir: ../mach-component-my-component/
    ```

    More info on [using MACH composer during development](../../topics/development/workflow.md)

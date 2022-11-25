# components

All component definitions in this configuration file.<br>
These components are used and configured separately [per site](./sites.md#components).

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
- `integrations` - Defines a list of integrations for the given component. It
  controls what Terraform variables are passed on to the components [Terraform
  module](../components/structure.md#terraform-module).<br>

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
- `azure` - Configuration block for [Azure-specific settings](#azure)

!!! tip "Development settings"
    In addition to the default set of component settings, a couple of settings
    can be defined during development.<br> These are not intendend to be used
    for a production deployment, but can facilitate local development:

    - `branch` - Configure the git branch of the component. Only used to
      facilitate the `mach update` CLI command.
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

    More info on [using MACH during development](../../topics/development/workflow.md)

## azure
Example:

```yaml
azure:
  service_plan: default
  short_name: apiexts
```

- `service_plan` - The service plan (defined in [`service_plans`](./global.md#service_plans))
  to use for this component. Set this to `default` if you want to use the
  MACH-managed Consumption plan.
- `short_name` - Short name to be used in cloud resources. Should be at most 10
  characters to avoid running into Resource naming limits.<br>
  Defaults to the given components `name`

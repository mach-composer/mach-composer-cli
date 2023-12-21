# Mach Composer

All MACH composer specific configurations options. These govern how the 
application as a whole will function, where variables can be sourced from 
and which plugins are loaded, among others.

## Example

```yaml
mach_composer:
  version: 1
  variables_file: variables.yml
```

## Schema

### Required

- `version` (Number) Schema version to be used for this configuration. Currently
  only version 1 is supported.

### Optional

- `variables_file` (String) Define a variables file. Can be used instead of
  using the `--var-file` option. See [variables](index.md#variables) for
  more information.
- `plugins` (List of Block) List of plugins to be used. See
  [plugins](../../plugins/index.md) for more information.
  By default, the amplience, aws, azure, commercetools, contentful and
  sentry plugins are loaded (
  see [below for nested schema](#nested-schema-for-plugins))
- `cloud` (Block) Cloud specific configuration. See
  [cloud](../../cloud/index.md) for more information. See [below for nested
  schema](#nested-schema-for-cloud)). If not set no cloud specific features
  will be used.
- `deployment` (Block) Deployment specific configuration. See
  [deployment](../../concepts/deployment/index.md) for more information. If not
  mach-composer will default to site-scoped deployments. See [below for nested
  schema](#nested-schema-for-deployment)).

## Nested schema for `plugins`

### Required

- `source` (String) The source of the plugin. This will be used to download the
  plugin from GitHub releases. The format is `organization/repository`.

### Optional

- `version` (String) The version of the plugin to be used. If not specified, the
  latest version will be used.

## Nested schema for `cloud`

### Required

- `organization` (String) The organization in mach-composer cloud this project
  belongs to.
- `project` (String) The project name in mach-composer cloud.

## Nested schema for `deployment`

{% include-markdown "./deployment.md" %}

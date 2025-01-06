# Overview

A configuration file can contain several sites with all different configurations
and all using a different mix of re-usable serverless microservice components.

It is common to have a single configuration file per environment since they
usually share the same general configurations.

## Schema

The root schema is show below.

```yaml
{% include "./schema.yaml" %}
```

This can however be extended through the use of plugins, which declare their own bits of configuration. For development
purposes it is possible to generate a JSON schema for the required configuration. This can be used to
configure autocompletion and syntax checking support in your favorite IDE.

To generate the schema, run the following command:

```bash
mach-composer schema
```

See the [CLI documentation](../cli/mach-composer_schema.md) for more information.

### Required

- `mach_composer` (Block) will determine the overall behaviour of
  the application. See [the mach_composer documentation](./mach_composer.md)
  for more information
- `global` (Block) will determine the global configuration for all sites. See [the
  global documentation](./global.md) for more information
- `sites` (List of Block) will determine the configuration for
  each site. See [the site documentation](./site.md) for more information
- `components` (List of Block) will determine the configuration for each
  component. See [the component documentation](./component.md) for more
  information

## General syntax

The configuration syntax follows the same basic rules as YAML, with some extra
features.

### Plugins extending functionality

Plugins can be used to extend the functionality of mach-composer. They can
be added to the configuration file using the `plugins` key in the
[`mach_composer` configuration](./mach_composer.md).

In essence a plugin will add additional terraform code to your module where
appropriate. These additional elements might add additional configuration
options on a global, site and site-component level which can change behaviour.
See [the documentation](../../plugins/index.md) of the
plugin for more information on what plugins are available, what they do and what
to add where.

### Including YAML files

Using the `$ref` syntax it is possible to load in other yaml files as part of
your configuration.

This can be used for example to manage your component definitions elsewhere
within the same directory like so;

```yaml
---
mach_composer: ...
global: ...
sites: ...
components:
  $ref: _components.yaml
```

## Variables

MACH composer support the usage of variables in a configuration file. This
is a good way to keep your configuration DRY and to keep sensitive
information separate.

The following types are supported;

- [`${component.}`](#component) component output references
- [`${var.}`](#var) variables file values
- [`${env.}`](#env) environment variables value

### Example

```yaml
mach_composer:
  version: 1
global:
  environment: ${env.MACH_ENVIRONMENT}
  cloud: aws
sites:
  - identifier: my-site
    aws:
      account_id: 1234567890
      region: eu-central-1
    endpoints:
      public: api.tst.mach-example.net
    components:
      - name: infra
      - name: payment
        variables:
          sns_topic: ${components.infra.sns_topic_arn}
        secrets:
          stripe_secret_key: ${var.stripe_secret}
```

- `${components.infra.sns_topic_arn}` uses the `sns_topic_arn` Terraform output
  as a value for the payment component
- `${var.stripe_secret}` reads the `stripe_secret` from a variables file
- `${env.MACH_ENVIRONMENT}` reads the `MACH_ENVIRONMENT` environment variable

### `component`

**Usage** `${component.<component-name>.<output-value>}`

You can use this to refer to any [Terraform output](https://www.terraform.io/docs/language/values/outputs.html) that
another component has defined.

So for example if a component called "email" has the following outputs:

```terraform
# outputs.tf

output "sqs_queue" {
  value = {
    id  = aws_sqs_queue.email_queue.id
    arn = aws_sqs_queue.email_queue.arn
  }
}
```

These can then be used in the configuration:

```yaml
components:
  - name: order-notifier
    variables:
      email_queue_id: ${component.email.sqs_queue.id}
```

### `var`

**Usage** `${var.<variable-key>}`

This can be used for using values from a *variables file*. This variable file must be set by using the [
`--var-file` CLI option](../cli/mach-composer_apply.md#options):

```bash
mach-composer apply -f main.yml --var-file variables.yml
```

From the [example](#example) above, the following configuration line:

```yaml
stripe_secret_key: ${var.stripe_secret}
```

will use the `stripe_secret` value from the given variables file.

!!! info ""
These values can be nested, so it's possible to define a
`${var.site1.stripe.secret_key}` with your `variables.yml` looking like:

        ```yaml
        ---
        site1:
          stripe:
            secret_key: vRBNcBH2XuNvHwAoPdDnhs2XyeVMOT
        site2:
          stripe:
            secret_key: 2hzctJCLjyMjUL07BNSh3Nyjt6r7aC
    ```

!!! tip "Note on encryption"
    You can [encrypt your `variables.yml` using SOPS](../../howto/security/encrypt.md#encrypted-variables).

    When doing so, MACH composer won't render the variable files directly into
    your Terraform configuration but uses
    [terraform-sops](https://github.com/carlpett/terraform-provider-sops) to
    refer you the SOPS encrypted variables within the Terraform file.

### `env`

**Usage** `${env.<variable-name>}`

Use environment variables in your MACH configuration:

```bash
export MACH_ENVIRONMENT=test
mach-composer apply
```

Will replace `${env.MACH_ENVIRONMENT}` in our [example](#example) with `test`.

### Examples

For examples see the [examples](../../tutorial/examples/index.md) directory in
the tutorial section.

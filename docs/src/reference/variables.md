# Variables

MACH composer support the usage of variables in a configuration file.

The following types are supported;

- [`${component.}`](#component) component output references
- [`${var.}`](#var) variables file values
- [`${env.}`](#env) environment variables value
- [`${include.}`](#include) file includes

## Example

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
components: ${include(components.yml)}
```

- `${components.infra.sns_topic_arn}` uses the `sns_topic_arn` Terraform output as a value for the payment component
- `${var.stripe_secret}` reads the `stripe_secret` from a variables file
- `${include(components.yml)}` includes `components.yml` and injects it in the configuration

## `component`
**Usage** `${component.<component-name>.<output-value>}`

You can use this to refer to any [Terraform output](https://www.terraform.io/docs/language/values/outputs.html) that another component has defined.

So for example if a component called "email" has the following outputs:

```terraform
# outputs.tf

output "sqs_queue" {
    value = {
      id = aws_sqs_queue.email_queue.id
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

## `var`
**Usage** `${var.<variable-key>}`

This can be used for using values from a *variables file*. This variable file must be set by using the [`--var-file` CLI option](./cli.md#apply):

```bash
mach-composer apply -f main.yml --var-file variables.yml
```

From the [example](#example) above, the following configuration line:
```yaml
stripe_secret_key: ${var.stripe_secret}
```

will use the `stripe_secret` value from the given variables file.

!!! info ""
    These values can be nested, so it's possible to define a `${var.site1.stripe.secret_key}` with your `variables.yml` looking like:

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
    You can [encrypt your `variables.yml` using SOPS](../howto/security/encrypt.md#encrypted-variables).

    When doing so, MACH won't render the variable files directly into your Terraform configuration but uses [terraform-sops](https://github.com/carlpett/terraform-provider-sops) to refer you the SOPS encrypted variables within the Terraform file.

## `env`
**Usage** `${env.<variable-name>}`

Use environment variables in your MACH configuration:

```bash
export MACH_ENVIRONMENT=test
mach-composer apply
```

Will replace `${env.MACH_ENVIRONMENT}` in our [example](#example) with `test`.

## `include`
**Usage** `${include(<filename>)}`

Any valid YAML file can be included here.

!!! info "Using `!include`"
    The `${include(...)}` syntax has the same effect as using `!include ...` in your YAML file.

    However, when using [SOPS](../topics/../howto/security/encrypt.md) to encrypt your configuration file, this tag will get stripped.
    Therefor, MACH also supports the MACH-specific syntax.


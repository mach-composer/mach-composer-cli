# Component structure

## Required variables

MACH expects each component to have a certain set of variables defined.

What variables it needs to have defined is controlled by the [plugin's](../../plugins/index.md) setting.

If `integrations` is set to an empty list `[]`, no variables will be needed.

!!! tip ""
    An example of a component that takes no variables could be a component that
    creates custom product types in commercetools. This component operates with
    the same Terraform commercetools provider which is configured for the
    correct project already, so no additional information will be needed in the
    Terraform module itself.


### Optional variables

Regardless of what [plugins](../../plugins/index.md) the components has, two
variables must be set on the component if those are actually set in your MACH
configuration.

So only if you have the following configuration for your component:

```yaml
components:
  - name: my-component
    variables:
      FOO: bar
    secrets:
      SECRET: 12345
```

The following variables must be defined in your component:

```terraform
variable "variables" {
  type        = any
  description = "Generic way to pass variables to components."
}

variable "secrets" {
  type        = any
  description = "Map of secret values. Can be placed in a key vault."
}
```

## Integrations

By defining a set of `integrations` in the
[component definitions](../../reference/syntax/site.md#nested-schema-for-components), MACH knows what variables need
to be passed on to the components.

This way the components don't need to define **all possible variables** a
component might have.

Available plugins can be found [here](../../plugins/index.md).

By default, integrations are set on the given cloud provider. So when no
`integrations` definition is given, it defaults to `['aws']` in case of an AWS
deployment.

!!! tip "Non-cloud components"
    As an example; you might have a component defining some custom commercetools
    product types. No further cloud infrastructure is needed.<br>
    In this case, that component will have `integrations: ['commeretools']` and
    MACH won't pass any of the cloud-specific variables.

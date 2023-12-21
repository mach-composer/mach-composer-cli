# Plugins

Each integration brings its own set of managed resources that gets created by
MACH composer when configured for a specific MACH composer stack.

### commercetools

The commercetools integration will manage all resources that can be
[configured](../../../plugins/commercetools.md#configuration) in the
[commercetools settings](../../../plugins/commercetools.md).

It uses the [Terraform commercetools provider](https://registry.terraform.io/providers/labd/commercetools/latest/docs) for this.

!!! note ""
    [More info](../../../plugins/commercetools.md) about the commercetools
    integration.

### Sentry

It will create a [Sentry key](https://registry.terraform.io/providers/jianyuan/sentry/latest/docs/resources/key)
per component that has the [Sentry integration](../../../plugins/sentry.md) defined.

!!! note ""
    [more info](../../../plugins/sentry.md) about the Sentry integration

### Contentful

It uses the [Terraform contentful provider](https://registry.terraform.io/providers/labd/contentful/latest)
to manage the following resources per site:

- Contentful space
- An API key

!!! note ""
    [more info](../../../plugins/contentful.md) about the Contentful integration

### Amplience

It will use the information from the [amplience settings](../../../plugins/amplience.md)
to pass to the components that have `amplience` included in their integrations.

!!! note ""
    [more info](../../../plugins/amplience.md) about the Amplience integration

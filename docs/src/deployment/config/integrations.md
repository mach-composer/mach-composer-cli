# Integrations

Each integration brings it's own set of managed resources that gets created by MACH when configured for a specific MACH stack.

### commercetools

The commercetools integration will manage all resources that can be [configured](../../integrations/commercetools.md#configuration) in the [commercetools settings](../../syntax.md#commercetools).

It uses the [Terraform commercetools provider](https://registry.terraform.io/providers/labd/commercetools/latest/docs) for this.

!!! note ""
    [More info](../../integrations/commercetools.md) about the commercetools integration.

### Sentry

It will create a [Sentry key](https://registry.terraform.io/providers/jianyuan/sentry/latest/docs/resources/key) per component that has the [Sentry integration](../../integrations/sentry.md) defined.

!!! note ""
    [more info](../../integrations/sentry.md) about the Sentry integration

### Contentful

It uses the [Terraform contentful provider](https://registry.terraform.io/providers/labd/contentful/latest) to manage the following resources per site:

- Contentful space
- An API key

!!! note ""
    [more info](../../integrations/contentful.md) about the Contentful integration
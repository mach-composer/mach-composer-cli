# Component structure

## Required variables

MACH expects each component to have a certain set of variables defined.

What variables it needs to have defined is controlled by the [integrations](../syntax/components.md) setting.

If `integrations` is set to an empty list `[]`, no variables will be needed.

!!! tip ""
    An example of a component that takes no variables could be a component that creates custom product types in commercetools. This component operates with the same Terraform commercetools provider which is configured for the correct project already, so no additional information will be needed in the Terraform module itself.


### Optional variables

Regardless of what [integration](#integrations) the components has, two variables must be set on the component if those are actually set in your MACH configuration.

So only if you have following configuration for your component:

```yaml
components:
  - my-component
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

By defining a set of `integrations` in the [component definitions](../syntax/components.md), MACH knows what variables need to be passed on to the components.

This way the components don't need to define **all possible variables** a component might have.

Available integrations are:

- `aws`
- `azure`
- `commercetools`
- `sentry`
- `contentful`
- `amplience`
- `apollo_federation`

By default, integrations are set on the given cloud provider. So when no `integrations` defintion is given, it defaults to `['aws']` in case of an AWS deployment.

!!! tip "Non-cloud components"
    As an example; you might have a component defining some custom commercetools product types. No further cloud infrastructure is needed.<br>
    In this case, that component will have `integrations: ['commeretools']` and MACH won't pass any of the cloud-specific variables.

### cloud integration

The following variables are expected in case of either `aws` or `azure`.

```terraform
variable "component_version" {
  type        = string
  description = "Version to deploy"
}

variable "environment" {
  type        = string
  description = "Specify what environment it's in (e.g. `test` or `production`)"
}

variable "site" {
  type        = string
  description = "Identifier of the site"
}

variable "variables" {
  type        = any
  description = "Generic way to pass variables to components."
}

variable "secrets" {
  type        = any
  description = "Map of secret values. Can be placed in a key vault."
}

variable "tags" {
  type        = map(string)
  description = "Tags to be used on resources."
}
```

!!! info "Cloud specific variables"
      See [AWS variables](./aws.md#terraform-variables) and [Azure variables](./azure.md#terraform-variables) for cloud-specific variables that a component needs in addition to the base set.


### commercetools

The following variable is given when `commercetools` integration is defined.

```terraform
variable "ct_project_key" {
  type        = string
  description = "commercetools project key"
}

variable "ct_api_url" {
  type        = string
  description = "commercetools API URL"
}

variable "ct_auth_url" {
  type        = string
  description = "commercetools Auth URL"
}

variable "ct_stores" {
  type = map(object({
    key       = string
    variables = map(string)
    secrets   = map(string)
  }))
  default = {}
}
```

### sentry

The following variable is given when `sentry` integration is defined.

```terraform
variable "sentry_dsn" {
  type        = string
  default     = ""
  description = "Sentry DSN - only when Sentry is configured"
}
```

### contentful

The following variable is given when `contentful` integration is defined.

```terraform
variable "contentful_space_id" {
  type        = string
  description = "Contentful Space ID"
}
```

### amplience

The following variable is given when `amplience` integration is defined.

```terraform
variable "amplience_client_id" {
  type        = string
  description = "Amplience client id"
}

variable "amplience_client_secret" {
  type        = string
  description = "Amplience client secret"
}

variable "amplience_hub_id" {
  type        = string
  description = "Amplience hub id"
}
```

### apollo federation

The following variable is given when `apollo_federation` integration is defined.

```terraform
variable "apollo_federation" {
  type = object({
    api_key       = string
    graph         = string
    graph_variant = string
  })
}
```

## Serverless function

The component might contain code for a serverless function to run on the cloud provider.

What kind of language/runtime is used for that is irrelevant to MACH. Two things the component needs to contain:

- **Build/deploy script** to build, package and upload the serverless function to a repository
- A **Terraform configuration** for the serverless function

## Cloud provider specifics

A component is always tailored for a specific cloud provider.

Continue for details about the specifics:

- [Azure components](./azure.md)
- [AWS components](./aws.md)

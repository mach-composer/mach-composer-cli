# MACH configuration deployment

A MACH configuration deployment (or simply put: **MACH deployment**) will generate and apply a Terraform configuration **per site** so that it can deploy

#### 1. MACH-managed resources
The resources that are managed by MACH depend on the cloud integration:

- AWS (See [AWS deployments](./aws.md))
- Azure (See [Azure deployments](./azure.md))
#### 2. Integration resources
Resources needed for the integrations such as

- [commercetools](./integrations.md#commercetools)
- [sentry](./integrations.md#sentry)
- [contentful](./integrations.md#contentful)

#### 3. Components
Since components are loaded into the configuration as [Terraform modules](../../components/structure.md#terraform-module), during a MACH deployment the resources defined in the component will get created.

1. The [**first stage**](../components.md) of a component deployment (uploading the assets to a component repository) is done before a component is deployed as part of a MACH stack.

2. The [**second stage**](./components.md) is getting the previously deployed component assets actually up and running in your MACH stack and to create other necessary resources.

More info about the [second stage deployment](./components.md).

!!! info "Component deployment - first and second stage"
    Not all components have a '*first stage*' which means: some components might **just** have a Terraform configuration to be applied and no serverless function assets.<br>
    In that case, there is no need of a '*first stage*' component deployment.


## Providing credentials

MACH needs to be able to access:

- The components repositories
- The AWS account / Azure subscription it needs to manage resources in
  
TODO: Describe steps
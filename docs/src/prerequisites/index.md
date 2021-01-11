# Prerequisites

Before MACH can be used to deploy your environment, a couple of things need to be prepared.

Wether you're using AWS or Azure (or a mix of both), a choice has to be made where MACH stores its generated Terraform state. This can be done on an Azure Blob Container or AWS S3 bucket.

## Development and build environment

The machine MACH runs on needs to have [Docker](https://www.docker.com) installed.

!!! info  "From source"
    It's also possible to run MACH from source.<br>
    In order to do so, the following needs to be installed on the system:

    - Python 3.8
    - Terraform 0.13

## commercetools

!!! note "Optional"
    This step is only necessary if you are going to use the commercetools integration in your MACH stack

Create a API client *'to rule them all'*.

Required scopes:

- `manage_api_clients`
- `manage_project`
- `view_api_clients`

!!! note ""
    This client is used the MACH composer to create other necessary commercetools clients for each individual component.

Use the credentials for this client to configure each site's [commercetools settings](../syntax/sites.md#commercetools).

## Cloud environment

Setup your [Azure](./azure.md) or [AWS](./aws.md) environment.

## Other built-in integrations

MACH also provides other integrations that you can include in your MACH environment.

- [Commercetools integration](../integrations/commercetools.md)
- [Contentful integration](../integrations/contentful.md)
- [Amplience integration](../integrations/amplience.md)
- [Sentry integration](../integrations/sentry.md)
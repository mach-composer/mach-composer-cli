# Getting started

## Setup configuration file

To setup a MACH configuration, create a YAML file with the following structure

``` yaml
---
general_config:
    environment: test
    cloud: azure
    terraform_config:
        azure_remote_state:
            resource_group: my-shared-rg
            storage_account: mysharedsaterra
            container_name: tfstate
            state_folder: test
    azure:
        tenant_id: e180345a-b3e1-421f-b448-672ab50d8502
        subscription_id: 086bd7e7-0755-44ab-a730-7a0b8ad4883f
        region: westeurope
sites:
    - identifier: my-site
      commercetools:
          project_key: my-site-tst
          client_id: ...
          client_secret: ...
          scopes: manage_project:my-site-tst manage_api_clients:my-site-tst view_api_clients:my-site-tst
          languages:
              - en-GB
              - nl-NL
          currencies:
              - GBP
              - EUR
          countries:
              - GB
              - NL
      components:
          - name: api-extensions
            variables:
                ORDER_PREFIX: mysitetst
components:
    - name: api-extensions
      short_name: apiexts
      source: git::ssh://git@github.com/your-project/components/api-extensions-component.git//terraform
      version: e638e57
```

See [Syntax](./syntax.md) for all configuration options.

## Deploy using MACH configuration

You can deploy your current configuration by running

```bash
$ mach generate
```

If you wish to review the changes before applying them, run

```bash
$ mach plan
```

!!! tip "Using Docker image"
    You can invoke MACH by running the Docker image:  
    `$ docker run --rm --volume $(pwd):/code docker.pkg.github.com/labd/mach-composer/mach plan`

    You do need to provide the docker container with the necessary environment variables to be able to authenticate with the cloud provider. More info on that in the [deployment section](./deployment/config.md)


## Additional options
See the [Deployment section](./deployment/config.md) for more deployment options.

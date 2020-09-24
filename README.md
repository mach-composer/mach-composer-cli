# MACH Framework

The idea is to define a configuration file in yml with one or more Commercetools
sites with all the required configuration for that site, what components to use
and which version of components. 

I suggest keeping component versions global to avoid version hell between
different sites, and a git commit for versions (preferably just the latest
master git commit). Eventually you're gonna end up with 15 sites, each with 6
components. By always releasing the latest version maintenance and releasing
becomes a lot easier compared to the alternative.


### Site
A site is one Commercetools instance (with one or more Commercetools stores) and
a cloud environment (Azure resource group or AWS account)

### Component
A component is an isolated Terraform module. F.e. a component can define
Commercetools product types (so you can use different product types per site).
Typically the other use case for a component is a serverless function that
either adds an API to the site or listens to Commercetools subscriptions (f.e.
order created or order status changes).

## Installation

- Make sure you have Terraform installed  
tip: Use [tfenv](https://github.com/tfutils/tfenv) to support multiple versions of Terraform)
- Install https://github.com/labd/terraform-provider-commercetools  
note: you have to specify the version in the file name (f.e. terraform-provider-commercetools_v0.21.1 in the Terraform plugin directory)
- Create a new virtualenv with Python 3.8
- Run the following commands:

```
$ make install
$ az login
``` 
 

## Running MACH

To generate the files:

`mach generate # generates all available configs.` 

`mach generate -f main.yml`

To plan Terraform:

`mach plan`

To apply Terraform config:

`mach apply`

## Checking for updates

MACH can check your components for available updates.

To do this, run:

`mach update -f main.yml`


## Code style
The Python source code should be formatted using
[black](https://github.com/python/black) and the JavaScript code should be
formatted using [prettier](https://prettier.io/). You can use `make format`
to run both black and prettier.

This project uses [pre-commit](https://pre-commit.com) to validate the changed
files before they are committed. You can install it on MacOS using brew:

    $ brew install pre-commit

In the repository you need to register the hooks in git the first time using:

    $ pre-commit install

The pre-commit config (`.pre-commit-config.yaml`) currently runs black and
flake8.

## Example yaml file

```yaml
---
general_config:
    environment: test
    terraform_config:
        azure_remote_state:
            resource_group_name: my-shared-rg
            storage_account_name: mysharedsaterra
            container_name: tfstate
            state_folder: test
sites:
    - identifier: my-site
      azure:
          tenant_id: e180345a-b3e1-421f-b448-672ab50d8502
          subscription_id: 086bd7e7-0755-44ab-a730-7a0b8ad4883f
          region: westeurope
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
                CT_CLIENT_SCOPES: manage_products:my-site-tst manage_orders:my-site-tst
                ORDER_PREFIX: mysitetst
components:
    - name: api-extensions
      short_name: apiexts
      source: git::ssh://git@github.com/your-project/components/api-extensions-component.git//terraform
      version: e638e57
```

## Terminology and settings

### Azure

#### service_object_id
The Azure Active Directory principal ID which can access the necessary resources.  
[More info](https://docs.microsoft.com/en-us/azure/active-directory/develop/app-objects-and-service-principals#service-principal-object)

#### tenant_id
The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.  
[More info](https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-scenarios#tenants)

## TODO

- Revisit discussion about secrets filled in on the main.yml or fill in on Azure Portal.
- Add optional encryption (f.e. using sops https://github.com/mozilla/sops/ like Danone and integrate it with Azure Key Vault in shared rg)
- Switch to Ubuntu, alpine is always difficult to install extra dependencies.

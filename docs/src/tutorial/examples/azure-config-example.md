# Azure MACH config

The below is a simplified example built on Azure and including Commercetools.

```yaml
mach_composer:
  version: 1
  plugins:
    aws:
      source: mach-composer/azure-minimal #unopinioned version of Azure integration
      version: 0.1.0
    commercetools:
      source: mach-composer/commercetools
      version: 0.1.8
global:
  environment: test
  cloud: azure
  terraform_config:
    remote_state:
      plugin: azure
      resource_group: mach-shared-we-rg
      storage_account: machsharedwesaterra
      container_name: tfstate
      state_folder: test
  azure:
    tenant_id: <your-tenant-id>
    subscription_id: <your-subscription-id>
    region: westeurope
    resources_prefix: ""
sites:
  - identifier: my-site
    commercetools:
      project_key: my-site
      client_id: <client-id>
      client_secret: <client-secret>
      scopes: manage_api_clients:my-site manage_project:my-site view_api_clients:my-site
      project_settings:
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
      - name: your-component
        variables:
          FOO_VAR: my-value
        secrets:
          MY_SECRET: secretvalue
components:
  - name: your-component
    source: git::https://github.com/<username>/<your-component>.git//terraform
    version: 0.1.0
    azure:
      short_name: yourcomp
    integrations:
      - azure
      - commercetools
```

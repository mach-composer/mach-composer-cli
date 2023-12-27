# GCP Mach Config

The below is a simplified example built on GCP and including Commercetools.

```yaml
mach_composer:
  version: 1
  plugins:
    gcp:
      source: mach-composer/gcp
      version: 0.1.4
    commercetools:
      source: mach-composer/commercetools
      version: 0.1.8

global:
  environment: test
  cloud: gcp
  terraform_config:
    remote_state:
      plugin: gcp
      bucket: <your bucket>
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
    integrations:
      - gcp
      - commercetools
```

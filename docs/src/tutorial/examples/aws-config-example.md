# AWS MACH config

The below is a simplified example of a full configuration that should work.

```yaml
mach_composer:
  version: 1
  plugins:
    aws:
      source: mach-composer/aws
      version: 0.1.0
    commercetools:
      source: mach-composer/commercetools
      version: 0.1.8

global:
  environment: test
  terraform_config:
    aws_remote_state:
      bucket: <your bucket>
      key_prefix: mach
      region: eu-central-1
  cloud: aws
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
    aws:
      account_id: 123456789
      region: eu-central-1
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
      - aws
      - commercetools
```
# Configuration syntax

A configuration file can contain several sites with all different configurations and all using a different mix of re-usable serverless microservice components.

It is common to have a single configuration file per environment since they usually share the same general configurations.

The configuration file has the following structure:

- **[mach_composer](./mach.md)**
- **[global](./global.md)**
    - **[environment](./global.md)**
    - **[terraform_config](./global.md#terraform_config)**
    - **[cloud](./global.md)**
    - [azure](./global.md#azure)
    - [sentry](./global.md#sentry)
    - [contentful](./global.md#contentful)
    - [amplience](./global.md#amplience)
- **[sites](./sites.md)**
    - **[identifier](./sites.md)**
    - [commercetools](./sites.md#commercetools)
    - [contentful](./sites.md#contentful)
    - [amplience](./sites.md#amplience)
    - [azure](./sites.md#azure)
    - [aws](./sites.md#aws)
    - [stores](./sites.md#stores)
    - [components](./sites.md#components)
- [components](./components.md)


!!! tip "JSON schema"
    A JSON schema for the syntax is [available on GitHub](https://github.com/labd/mach-composer/blob/master/schema.json). This can be used to configure IntelliSense autocompletion support in VSCode.

## Including YAML files

Using the `!include` tag it's possible to load in another yaml file as part of your configuration.

This can be used for example to manage your component definitions elsewhere like so;

```yaml
---
mach_composer: ...
global: ...
sites: ...
components: !include components.yml
```

Or load them from an external location;

=== "Git"
    ```yaml
    ---
    mach_composer: ...
    global: ...
    sites: ...
    components: !include git::https://github.com/your-org/mach-config.git@9f42fe2//components.yml
    ```
=== "HTTPS"
    ```yaml
    ---
    mach_composer: ...
    global: ...
    sites: ...
    components: !include https://www.your-org.com/mach/components.yml
    ```

!!! info "SOPS compatability"
    Using the `!include` tag doesn't play well with [SOPS](../../howto/security/encrypt.md) yet.

    SOPS will remove the `!` from the tag when encrypting.

## Full example

=== "AWS"
    ```yaml
    ---
    mach_composer:
      version: 1.0.0
    global:
      environment: test
      cloud: aws
      terraform_config:
        aws_remote_state:
          bucket: my-tfstate-tst
          key_prefix: mach-composer-tst
          lock_table: my-tfstate-tst-lock
          region: eu-central-1
        providers:
          aws: 3.28.0
      amplience:
        client_id: 2d02e7da-053a-449c-a954-bc0bb22f70e6
        client_secret: ybPPCtPlgzBCuAbaWoD8eDLAifVgJhT2HPepn4mCCJZibmGKno
      sentry:
        auth_token: 56jIifdp7viYv1DMgTBZq7DZlUbV4mF7nL036KN5
        organization: my-org
        project: my-project
        rate_limit_window: 21600
        rate_limit_count: 100
      contentful:
        cma_token: sbhdaxeux8l54q9no2nng1pgopsu9u
        organization_id: OcF6DjbBh9AtQXcy8CZ0
    sites:
      - identifier: site-tst
        endpoints:
          main: https://site-tst.my-commerce-services.net
          internal: https://site-tst.my-commerce-services.internal
        aws:
          account_id: 00000000001
          region: eu-central-1
        amplience:
          hub_id: UWebNaB5MbYIVyWPUzZE
        contentful:
          space: my-site
        commercetools:
          project_key: site-tst
          client_id: rA53Hcm8ZtBuO8QwgmRZ
          client_secret: A94ElFbJYEt2sTyBzWcFToWp5lipS1
          scopes: manage_project:site-tst view_api_clients:site-tst manage_api_clients:site-tst
          token_url: https://auth.us-central1.gcp.commercetools.com
          api_url: https://api.us-central1.gcp.commercetools.com
          languages:
            - en-US
          countries:
            - US
          currencies:
            - USD
          channels:
            - key: US
              roles:
                - ProductDistribution
                - InventorySupply
              name:
                en-US: USA
          stores:
            - key: US
              name:
                en-US: USA
              distribution_channels:
                - US
              supply_channels:
                - US
        apollo_federation:
          api_key: service:mach-poc-123:Abc00kHbB89h
          graph: mach-poc-123
          graph_variant: current
        components:
          - name: ecommerce-content
          - name: api-extensions
            variables:
              ORDER_PREFIX: mysitetst
          - name: payment
            variables:
              PAYMENT_PUBLIC_KEY: 6H6HHLrQH03OZtUVfWwY
            secrets:
              PAYMENT_PRIVATE_KEY: eGye6jX9FNp3peytjeyQEmnJsBB6gG
    components:
      - name: api-extensions
        source: git::https://github.com/your-company/api-extensions-component//terraform
        integrations: [aws, commercetools, sentry]
        version: 0dd4814
      - name: payment
        source: git::https://github.com/your-company/payment-component//terraform
        integrations: [aws, commercetools, sentry]
        endpoints:
          main: main
          internal: internal
        version: 06b3cf8
      - name: ecommerce-content
        integrations: [amplience, contentful]
        source: git::https://github.com/your-company/ecommerce-content//terraform
        version: a410ce6

    ```

=== "Azure"
    ```yaml
    ---
    mach_composer:
      version: 1.0.0
    global:
      environment: test
      cloud: azure
      terraform_config:
        azure_remote_state:
          resource_group: my-shared-we-rg
          storage_account: mysharedwesaterra
          container_name: tfstate
          state_folder: test
        providers: 
          azure: 2.51.0
      amplience:
        client_id: 2d02e7da-053a-449c-a954-bc0bb22f70e6
        client_secret: ybPPCtPlgzBCuAbaWoD8eDLAifVgJhT2HPepn4mCCJZibmGKno
      sentry:
        auth_token: 56jIifdp7viYv1DMgTBZq7DZlUbV4mF7nL036KN5
        organization: my-org
        project: my-project
        rate_limit_window: 21600
        rate_limit_count: 100
      contentful:
        cma_token: sbhdaxeux8l54q9no2nng1pgopsu9u
        organization_id: OcF6DjbBh9AtQXcy8CZ0
      azure:
        service_object_ids:
          developers: 034f9027-4934-4565-99c6-0fb58b2be1fe # developer TST group
        tenant_id: 486df2a5-2faa-47eb-b1b5-5ace2032e59b
        subscription_id: e12224d5-997d-4c72-a576-60f4919296db
        region: westeurope
        resources_prefix: my-
        frontdoor:
          resource_group: my-shared-we-rg
          suppress_changes: true
    sites:
      - identifier: cas-site-tst
        endpoints:
          main: https://cas-site-tst.my-commerce-services.net
        azure:
          region: centralus
          alert_group:
            name: critical
            alert_emails:
              - alerting@example.com
            webhook_url: https://example.com/api/alert-me/
          service_plans:
            default:
              kind: elastic
              tier: ElasticPremium
              size: EP1
        amplience:
          hub_id: UWebNaB5MbYIVyWPUzZE
        contentful:
          space: my-site
        commercetools:
          project_key: cas-site-tst
          client_id: rA53Hcm8ZtBuO8QwgmRZ
          client_secret: A94ElFbJYEt2sTyBzWcFToWp5lipS1
          scopes: manage_project:cas-site-tst view_api_clients:cas-site-tst manage_api_clients:cas-site-tst
          token_url: https://auth.us-central1.gcp.commercetools.com
          api_url: https://api.us-central1.gcp.commercetools.com
          languages:
            - en-US
          countries:
            - US
          currencies:
            - USD
          channels:
            - key: US
              roles:
                - ProductDistribution
                - InventorySupply
              name:
                en-US: USA
          stores:
            - key: US
              name:
                en-US: USA
              distribution_channels:
                - US
              supply_channels:
                - US
        apollo_federation:
          api_key: service:mach-poc-123:Abc00kHbB89h
          graph: mach-poc-123
          graph_variant: current
        components:
          - name: ecommerce-content
          - name: api-extensions
            variables:
              ORDER_PREFIX: mysitetst
          - name: payment
            variables:
              PAYMENT_PUBLIC_KEY: 6H6HHLrQH03OZtUVfWwY
            secrets:
              PAYMENT_PRIVATE_KEY: eGye6jX9FNp3peytjeyQEmnJsBB6gG
    components:
      - name: api-extensions
        source: git::https://github.com/your-company/api-extensions-component//terraform
        integrations: [azure, commercetools, sentry]
        version: 0dd4814
        azure:
          short_name: apiexts
          service_plan: default
      - name: payment
        source: git::https://github.com/your-company/payment-component//terraform
        integrations: [azure, commercetools, sentry]
        endpoints:
          main: main
        azure:
          service_plan: default
        version: 06b3cf8
      - name: ecommerce-content
        integrations: [amplience, contentful]
        azure:
          short_name: pt_ecom
        source: git::https://github.com/your-company/ecommerce-content//terraform
        version: a410ce6
    ```

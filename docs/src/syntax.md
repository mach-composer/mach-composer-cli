# Configuration syntax

A configuration file can contain several sites with all different configurations and all using a different mix of re-usable serverless microservice components.

It is common to have a single configuration file per environment since they usually share the same general configurations.

The configuration file has the following structure:

- **[general_config](#general_config)**
    - **[environment](#general_config)**
    - **[terraform_config](#terraform_config)**
    - **[cloud](#general_config)**
    - [azure](#azure)
    - [sentry](#sentry)
    - [contentful](#contentful)
    - [amplience](#amplience)
- **[sites](#sites)**
    - **[identifier](#sites)**
    - [commercetools](#commercetools)
    - [contentful](#contentful_1)
    - [amplience](#amplience_1)
    - [azure](#azure_1)
    - [aws](#aws)
    - [stores](#stores)
    - [components](#component-configurations)
- [components](#components)


!!! tip "JSON schema"
    A JSON schema for the syntax is [available on GitHub](https://github.com/labd/mach-composer/blob/master/schema.json). This can be used to configure IntelliSense autocompletion support in VSCode.

## general_config
All 'shared' configuration that applies to all sites.

- **`environment`** - (Required) [environment](#environment) Identifier for the environment. For example `development`, `test` or `production`.<br>
Is used to set the `environment` variable of any [Terraform component](./components/structure.md#terraform-component)
- **`terraform_config`** - (Required) [terraform_config](#terraform_config) block
- `cloud` - Either `azure` or `aws`
- `azure` - [Azure](#azure) block
- `sentry` - [Sentry](#sentry) block
- `contentful` - [Contentful](#contentful) block


### terraform_config
Terraform configuration block.

Can be used to configure the state backend and Terraform provider versions.

- **`azure_remote_state`** - [Azure](#azure_remote_state) remote state configuration
- **`aws_remote_state`** - [AWS](#aws_remote_state) remote state configuration
- `providers` - [Providers](#providers) configuration block

#### azure_remote_state
An Azure state backend can be defined as:

```yaml
terraform_config:
  azure_remote_state:
    resource_group: <your resource group>
    storage_account: <storage account name>
    container_name: <container name>
    state_folder: <state folder>
```

!!! tip ""
    A good convention is to give the state_folder the same name as [`environment`](#environment)

#### aws_remote_state
An AWS S3 state backend can be defined as:

```yaml
terraform_config:
  aws_remote_state:
    bucket: mach-statefiles
    key_prefix: test-statefiles
    role_arn: arn:aws:iam::1234567890:role/deploy
```

- **`bucket`** - (Required) S3 bucket name
- **`key_prefix`** - (Required) Key prefix for each individual Terraform state
- `role_arn` - Role ARN to access S3 bucket with
- `lock_table` - DynamoDB lock table
- `encrypt` - Enable server side encryption of the state file. Defaults to `True`

#### providers

Can be used to overwrite the MACH defaults for the Terraform provider versions.

Example:

```yaml
terraform_config:
  providers:
    aws: 3.21.0
```

- `aws` [aws provider](https://registry.terraform.io/providers/hashicorp/aws) version overwrite
- `azure` [azurerm provider](https://registry.terraform.io/providers/hashicorp/azurerm) version overwrite
- `commercetools` [commercetools provider](https://registry.terraform.io/providers/labd/commercetools) version overwrite
- `sentry` [sentry provider](https://registry.terraform.io/providers/jianyuan/sentry) version overwrite
- `contentful` [contentful provider](https://registry.terraform.io/providers/labd/contentful) version overwrite


!!! tip "Cache your providers"
    If you're overwriting the provider versions, make sure you [mount the plugins cache](./deployment/config/index.md#cache-terraform-providers)

### sentry
Defines a Sentry configuration.

This could be a predefined DSN to be used in the components, or MACH can manage the keys for you and pass the correct DSN to the components to be used.

=== "Example (managed)"
    ```yaml
    sentry:
      auth_token: <your-sentry-auth-token>
      organization: <organization-name>
      project: <project-name>
    ```
=== "predefined DSN"
    ```yaml
    sentry:
      dsn: https://LhNrqROZRIl2c5ciidkN82DObJfgtiLd@sentry.io/123456
    ```

- **`dsn`** - DSN to use in the components

or

- **`auth_token`** - Auth token to manage keys with
- **`organization`** - Organization name
- **`project`** - Project to create the key for
- `rate_limit_window` - The rate limit window that applies to a generated key
- `rate_limit_count` - The rate limit count that applies to a generated key

When defined, a `sentry` integration can be used in the components to expose a Sentry DSN value.

!!! tip ""
    The Sentry settings can be overwritten on [site](#sentry_1) and [component](#sentry_2) level

### azure

General Azure settings. Values can be overwritten [per site](#azure_1).

Example:
```yaml
azure:
  tenant_id: f2e03b8b-fe10-4fbc-9f5c-76dad9ac52e2
  subscription_id: a5b51c09-a2da-45b8-918a-67cf42456ab3
  region: westeurope
  resources_prefix: my-
  service_object_ids:
    gitlab-sp: d1114ea6-88f9-45b2-9de4-031291090380 # gitlab-sp
    developers: 3d280212-934f-4d32-876d-1b541a7697ba # developers tst group
```

- **`tenant_id`** - (Required)
- **`subscription_id`** - (Required)
- **`region`** - (Required)
- `resources_prefix` - Prefix to be used for all Azure resources, for example `my-`
- `frontdoor` - [Front-door](#frontdoor) settings
- `service_object_ids` - Map of objects IDs that should have access to things like KeyVaults created for components


!!! tip "Don't lock yourself out"
    Make sure that, as a minimum, you set the `service_object_ids` to the objects IDs of the users or groups that perform the `mach apply`


#### frontdoor

Example:
```yaml
frontdoor:
  resource_group: my-shared-rg
  dns_zone: my-services-domain.net
  ssl_key_vault_name: mysharedwekvcdn
  ssl_key_vault_secret_name: wildcard-my-services-domain-net
  ssl_key_vault_secret_version: IOlB8XmYLH1keYcpkcji23sp
```

- **`resource_group`** - (Required)
- **`dns_zone`** - (Required)
- **`ssl_key_vault_name`** - (Required)
- **`ssl_key_vault_secret_name`** - (Required)
- **`ssl_key_vault_secret_version`** - (Required)


### contentful
Defines global Contentful credentials to manage the spaces

- **cma_token** - (Required)
- **organization_id** - (Required)


### amplience
Defines global Amplience credentials to manage hubs

- **client_id** - (Required)
- **client_secret** - (Required)

## sites
All site definitions.


- **`identifier`** - (Required)<br>
  Unique identifier for this site.<br>
  Will be used for the Terraform state and naming all cloud resources.
- `endpoints` - [Endpoint definitions](#endpoints) to be used in the API Gateway or Frontdoor routing
- `commercetools` - [commercetools configuration](#commercetools) block
- `sentry` - [Sentry configuration](#sentry_1) block
- `contentful` - [Contentful configuration](#contentful_1) block
- `azure` - [Azure](#azure_1) settings
- `aws` - [AWS](#aws_1) settings
- `components` - [Component configurations](#component-configurations)

### endpoints

Endpoint definitions to be used in the API Gateway or Frontdoor routing.

Each component might require a different endpoint. In the [component definition](#components) it can be defined which endpoint it expects. The actual endpoint can be defined here using the unique key.

Example:

```yaml
endpoints:
  main: api.tst.mach-example.net
  services: services.tst.mach-example.net
```

!!! info "Azure support"
    At the moment, this option is not supported when using Azure and simply ignored.

    For Azure, the endpoints that are created for the APIs are constructed by using the commercetools project key as DNS record. More on that in the [Azure routing](./deployment/config/azure.md#http-routing) section

### commercetools

commercetools configuration.

Example:

```yaml
commercetools:
  project_key: my-site-tst
  client_id: T9J5g5bJe-VV8aVvN5Q
  client_secret: FIo3PGHJDThCM17wok_irLakRzCA
  scopes: manage_api_clients:my-site-tst manage_project:my-site-tst view_api_clients:my-site-tst
  languages:
    - en-GB
    - nl-NL
  currencies:
    - GBP
    - EUR
  countries:
    - GB
    - NL
```

- **`project_key`** - (Required) commercetools project key
- **`client_id`** - (Required) API client ID
- **`client_secret`** - (Required) API client secret
- **`scopes`** - (Required) Required scopes for given API client ID.
- **`currencies`** - (Required) List of three-digit currency codes as per ISO 4217
- **`languages`** - (Required) List of IETF language tag
- **`countries`** - (Required) List of two-digit country codes as per [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- `token_url` - Defaults to `https://auth.europe-west1.gcp.commercetools.com`
- `api_url` - Defaults to `https://api.europe-west1.gcp.commercetools.com`
- `messages_enabled` - When false the creation of messages is disabled.<br>
  Defaults to True
- `channels` - List of [channel definitions](#channels)
- `taxes` - List of [tax definitions](#taxes)
- `stores` - List of [store definitions](#stores) if multiple (store) contexts are going to be used.
- `create_frontend_credentials` - Defines if frontend API credentials must be created
  Defaults to `true`

#### channels

Example
```yaml
channels:
  - key: INV
    roles:
      - InventorySupply
    name:
      en-GB: Inventory
    description:
      en-GB: Our main inventory channel
  - key: DIST
    roles:
      - ProductDistribution
    name:
      en-GB: Distribution
    description:
      en-GB: Our main distribution channel
```

- **`key`** - (Required) 
- **`roles`** - (Required) List of [channel roles](https://docs.commercetools.com/http-api-projects-channels#channelroleenum).<br>
    Can be one of `InventorySupply`, `ProductDistribution`, `OrderExport`, `OrderImport` or `Primary`
- `name` - Name of the channel. Localized string [^1]
- `description` - Description of the channel. Localized string [^1]

#### taxes

Defines tax rates for various countries.

- **`country`** - (Required) A two-digit country code as per [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- **`amount`** - (Required) Number Percentage in the range of [0..1]
- **`name`** - (Required) Tax rate name

#### stores
Defines [commercetools stores](https://docs.commercetools.com/http-api-projects-stores).

Example:
```yaml
stores:
  - name:
      en-GB: my store
    key: mystore
    distribution_channels:
      - EU-DIST
    supply_channels:
      - EU-SUPPL
```

- **`name`** - (Required) Name of the store. Localized string [^1]
- **`key`** - (Required) Store key
- `languages` - List of languages
- `distribution_channels` - List of supply channel keys used for [product projection store filtering](https://docs.commercetools.com/http-api-projects-productProjections#prices-beta)
- `supply_channels` - List of supply channel keys used for [product projection store filtering](https://docs.commercetools.com/http-api-projects-productProjections#prices-beta)


### sentry

Overwrites any value specified in the general configs [Sentry block](#sentry)

- `dsn` - DSN to use in the components
- `rate_limit_window` - The rate limit window that applies to a generated key
- `rate_limit_count` - The rate limit count that applies to a generated key

### contentful

Contentful configuration.

Example:

```yaml
contentful:
  space: "MySpace"
```

- **`space`** - (Required) Name of the space to be created
- `default_locale` - Set default locale. Defaults to "en-US"
- `cma_token` - Override default `cma_token` setting (site-specific)
- `organization_id` - Override default `organization_id` setting (site-specific)


### amplience

Amplience configuration.

Example:

- `client_id` - Overrides default `client_id` settings (site-specific)
- `client_secret` - Override default `client_secret` setting (site-specific)
- `hub_id` - Override default `hub_id` setting (site-specific)


### azure
Site-specific Azure settings.<br>
Can overwrite any value from the generic [Azure settings](#azure):

- `tenant_id`
- `service_object_ids`
- `frontdoor`
- `subscription_id`
- `region`

And adds the following exta attributes:

- `alert_group` - List of [Alert groups](#alert_group)
- `resource_group` - Name of an already existing resource group.

!!! warning
    Use `resource_group` with care.<br>
    By default, MACH will manage the site resource groups for you. If you add this option later, the managed resource group will get **deleted**.<br>
    So only use for new site definitions

#### alert_group
Example:

```yaml
alert_group:
  name: critical
  alert_emails:
    - alerting@example.com
  logic_app: my-shared-we-rg.my-shared-we-alerts-slack
  webhook_url: https://example.com/api/alert-me/
```

- `name` - (Required) The name of the alert group
- `alert_emails` - Hook alert group to these email addresses
- `webhook_url` - Hooks alert group to a webhook
- `logic_app` - Reference to a Logic App the alert group needs to be connected to.<br>
  Format is `<resource_group>.<logic_app_name>`

### aws
Site-specific AWS settings.

Example:

```yaml
aws:
  account_id: 1234567890
  region: eu-west-1
  deploy_role_arn: deploy
  extra_providers:
    - name: email
      region: eu-west-1
```

- **`account_id`** - (Required) AWS account ID for this site
- **`region`** - AWS region to deploy site in
- `deploy_role_arn` - The [IAM role](./prerequisites/aws#iam-deploy-role) ARN needed for deployment
- `extra_providers`


### Component configurations

Configures the components for the site. The must reference a defined component (defined in the [component definitions](#components))

Example:

```yaml
components:
  - name: api-extensions
    variables:
      ORDER_PREFIX: mysitetst
  - name: order-mailer
    variables:
      FROM_EMAIL: mach@example.com
    secrets:
      SENDGRID_API_KEY: my-api-token
    store_variables:
      brand-a:
        FROM_EMAIL: mach@brand-a.com
      other-brand:
        FROM_EMAIL: mach@other-brand.com
```

- **`name`** - (Required) Reference to a [component](#component) definition
- `variables` - Variables for this component
- `secrets` - Variables for this component that should be stored in a encrypted key-value store
- `store_variables` - Store-specific variables for this component
- `store_secrets` - Store-specific variables for this component that should be stored in a encrypted key-value store
- `health_check_path` - Defines a custom healthcheck path.<br>
  Overwrites the default `health_check_path` defined in the component definition
- `sentry` - [Sentry configuration](#sentry_2) block

#### sentry

Overwrites any value specified in the site configs [Sentry block](#sentry_1)

- `dsn` - DSN to use in the components
- `rate_limit_window` - The rate limit window that applies to a generated key
- `rate_limit_count` - The rate limit count that applies to a generated key

## components definitions

All component definitions in this configuration file.<br>
These components are used and configured separately [per site](#component-configurations).

Example:

```yaml
components:
  - name: api-extensions
    short_name: apiexts
    source: git::ssh://git@git.labdigital.nl/mach-components/api-extensions-component.git//terraform
    version: 3b8ab91
    endpoint: main
  - name: ct-products-types
    source: git::ssh://git@git.labdigital.nl/mach-components/ct-product-types.git//terraform
    version: 1.4.0
    integrations: ["commercetools"]
```

- `name` - Name of the component. To be used as reference in the site definitions.
- `version` - A Git commit hash or tag
- `source` - Source definition of the terraform module
- `short_name` - Short name to be used in cloud resources. Should be at most 10 characters to avoid running into Resource naming limits.<br>
  Defaults to the given components `name`
- `integrations` - Defines a list of integrations for the given component. It controls what Terraform variables are passed on to the components [Terraform module](./components/structure.md#terraform-module).<br>
  Defaults to `["azure"]` or `["aws"]`, depending on your cloud provider.<br>
  Could be any of:
    - `azure`
    - `aws`
    - `commercetools`
    - `contentful`<br>
- `endpoint` - Defines the endpoint that needs to connect to this component.<br>
  Will setup Frontdoor routing or pass API Gateway information when set.
- `health_check_path` - Defines a custom healthcheck path.<br>
  Defaults to `/<name>/healthchecks`
- `branch` - Configure the git branch of the component. Only used to facilitate the `mach update` CLI command.


[^1]: commercetools uses [Localized strings](https://docs.commercetools.com/http-api-types#localizedstring) to be able to define strings in mulitple languages.<br>
Whenever a localized string needs to be defined, this can be done in the following format:
```yaml
some-string:
  - en-GB: My value
  - nl-NL:  Mijn waarde
```

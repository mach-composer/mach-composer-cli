# sites
All site definitions.


- **`identifier`** - (Required)<br>
  Unique identifier for this site.<br>
  Will be used for the Terraform state and naming all cloud resources.
- `endpoints` - [Endpoint definitions](#endpoints) to be used in the API Gateway or Frontdoor routing
- `commercetools` - [commercetools configuration](#commercetools) block
- `sentry` - [Sentry configuration](#sentry) block
- `contentful` - [Contentful configuration](#contentful) block
- `azure` - [Azure](#azure) settings
- `aws` - [AWS](#aws) settings
- `components` - [Component configurations](#component-configurations)

## endpoints

Endpoint definitions to be used in the API Gateway or Frontdoor routing.

Each component might require a different endpoint. In the [component definition](./components.md) it can be defined which endpoint it expects. The actual endpoint can be defined here using the unique key.

Example:

```yaml
endpoints:
  main: api.tst.mach-example.net
  services: services.tst.mach-example.net
```

!!! info "Azure support"
    At the moment, this option is not supported when using Azure and simply ignored.

    For Azure, the endpoints that are created for the APIs are constructed by using the commercetools project key as DNS record. More on that in the [Azure routing](../../topics/deployment/config/azure.md#http-routing) section

## commercetools

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

### channels

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

### taxes

Defines tax rates for various countries.

- **`country`** - (Required) A two-digit country code as per [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- **`amount`** - (Required) Number Percentage in the range of [0..1]
- **`name`** - (Required) Tax rate name

### stores
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


## sentry

Overwrites any value specified in the general configs [Sentry block](#sentry)

- `dsn` - DSN to use in the components
- `rate_limit_window` - The rate limit window that applies to a generated key
- `rate_limit_count` - The rate limit count that applies to a generated key

## contentful

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


## amplience

Amplience configuration.

Example:

- `client_id` - Overrides default `client_id` settings (site-specific)
- `client_secret` - Override default `client_secret` setting (site-specific)
- `hub_id` - Override default `hub_id` setting (site-specific)


## azure
Site-specific Azure settings.<br>
Can overwrite any value from the generic [Azure settings](#azure):

- `tenant_id`
- `service_object_ids`
- `frontdoor`
- `subscription_id`
- `region`
- `service_plans`

And adds the following exta attributes:

- `alert_group` - List of [Alert groups](#alert_group)
- `resource_group` - Name of an already existing resource group.

!!! warning
    Use `resource_group` with care.<br>
    By default, MACH will manage the site resource groups for you. If you add this option later, the managed resource group will get **deleted**.<br>
    So only use for new site definitions

### alert_group
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

## aws
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


## Component configurations

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

- **`name`** - (Required) Reference to a [component](./components.md) definition
- `variables` - Variables for this component
- `secrets` - Variables for this component that should be stored in a encrypted key-value store
- `store_variables` - [commercetools store](#stores)-specific variables for this component
- `store_secrets` - [commercetools store](#stores)-specific variables for this component that should be stored in a encrypted key-value store
- `health_check_path` - Defines a custom healthcheck path.<br>
  Overwrites the default `health_check_path` defined in the component definition
- `sentry` - [Sentry configuration](#sentry_1) block
- `azure` - [Azure configuration](#azure_1) block

### sentry

Overwrites any value specified in the site configs [Sentry block](#sentry)

- `dsn` - DSN to use in the components
- `rate_limit_window` - The rate limit window that applies to a generated key
- `rate_limit_count` - The rate limit count that applies to a generated key


[^1]: commercetools uses [Localized strings](https://docs.commercetools.com/http-api-types#localizedstring) to be able to define strings in mulitple languages.<br>
Whenever a localized string needs to be defined, this can be done in the following format:
```yaml
some-string:
  - en-GB: My value
  - nl-NL:  Mijn waarde
```

### azure
Overwrites any value specified in the [component definition](./components.md#azure).

Example:

```yaml
azure:
  service_plan: premium
```

- `service_plan` - The service plan (defined in [`service_plans`](./general_config.md#service_plans)) to use for this component. Defaults to `default`

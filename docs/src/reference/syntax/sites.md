# sites
All site definitions.


- **`identifier`** - (Required)<br>
  Unique identifier for this site.<br>
  Will be used for the Terraform state and naming all cloud resources.
- `endpoints` - [Endpoint definitions](#endpoints) to be used in the API Gateway or Frontdoor routing
- `commercetools` - [commercetools configuration](#commercetools) block
- `sentry` - [Sentry configuration](#sentry) block
- `contentful` - [Contentful configuration](#contentful) block
- `apollo_federation` - [Apollo Federation configuration](#apollo-federation) block
- `azure` - [Azure](#azure) settings
- `aws` - [AWS](#aws) settings
- `components` - [Component configurations](#components)

## endpoints

Endpoint definitions to be used in the API Gateway or Frontdoor routing.

Each component might require a different endpoint. In the [component definition](./components.md) it can be defined which endpoint it expects. The actual endpoint can be defined here using the unique key.

Basic example:

```yaml
endpoints:
  main: api.tst.mach-example.net
  services: services.tst.mach-example.net
```

Complex example:

```yaml
endpoints:
  internal:
    url: internal-api.tst.mach-example.net
    zone: tst.mach-example.net
    aws:
      throttling_burst_limit: 5000
      throttling_rate_limit: 10000
      enable_cdn: true
```

- **`url`** - (Required) url of the endpoint
- `zone` - DNS zone to use, if missing it's determined based on the given `url`
- `throttling_burst_limit` - Set burst limit for API Gateway endpoints
- `throttling_rate_limit` - Set burst limit for API Gateway endpoints
- `enable_cdn` - Defaults to false. Sets a CDN in front of this endpoint for better global availability. For AWS creates a CloudFront distribution

## commercetools

commercetools configuration.

Example:

```yaml
commercetools:
  project_key: my-site-tst
  client_id: T9J5g5bJe-VV8aVvN5Q
  client_secret: FIo3PGHJDThCM17wok_irLakRzCA
  scopes: manage_api_clients:my-site-tst manage_project:my-site-tst view_api_clients:my-site-tst
```

- **`project_key`** - (Required) commercetools project key
- **`client_id`** - (Required) API client ID
- **`client_secret`** - (Required) API client secret
- **`scopes`** - (Required) Required scopes for given API client ID.
- `project_settings` - [Project settings](#project_settings) configuration block
- `token_url` - Defaults to `https://auth.europe-west1.gcp.commercetools.com`
- `api_url` - Defaults to `https://api.europe-west1.gcp.commercetools.com`
- `channels` - List of [channel definitions](#channels)
- `taxes` - List of [tax definitions](#taxes)
- `tax_categories` - List of [tax_category definitions](#tax_categories)
- `zones` - List of [zone definitions](#zones)
- `stores` - List of [store definitions](#stores) if multiple (store) contexts are going to be used.
- `frontend` - [Frontend configuration block](#frontend)

### project_settings

Configuration block to define [project settings](https://docs.commercetools.com/api/projects/project).

Example:

```yaml
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
```

- **`currencies`** - (Required) List of three-digit currency codes as per ISO 4217
- **`languages`** - (Required) List of IETF language tag
- **`countries`** - (Required) List of two-digit country codes as per [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- `messages_enabled` - When false the creation of messages is disabled.<br>
  Defaults to True

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

### zones

Defines shipping zones.

Example:
```yaml
zones:
  - name: US zone 1
    locations:
      - country: US
        state: Nevada
  - name: US zone 2
    locations:
      - country: US
  - name: CA
    locations:
      - country: CA
  - name: Benelux
    locations:
      - country: NL
      - country: BE
      - country: LU

```

- **`name`** - (Required) Name for the zone
- **`locations`** - List of [locations](#locations) that are part of this zone
- `description` - Zone description

#### locations

- **`country`** - (Required) Two-digit country codes as per [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- `state` - State

### taxes

Defines tax rates for various countries but will only create a single tax category.
Please see [tax_categories](#tax_categories) for setups demanding multiple tax categories.

**Cannot be used in conjunction with `tax_categories`**

- **`country`** - (Required) A two-digit country code as per [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- **`amount`** - (Required) Number Percentage in the range of [0..1]
- **`name`** - (Required) Tax rate name
- **`included_in_price`** - Tax included in price

### tax_categories

Defines [commercetools tax categories](https://docs.commercetools.com/tutorials/explore-merchant-center#creating-tax-categories) with commercetools tax rates.

**Cannot be used in conjunction with `taxes`**

Example:
```yaml
tax_categories:
- name: Non-standard taxes
  key: non-standard-taxes
  rates:
  - name: rate 2
    amount: 0.02
    country: GB
    included_in_price: false
  - name: rate 10
    amount: 0.1
    country: NL
    included_in_price: true
- name: Standard taxes
  key: some-standard
  rates:
  - name: rate 21
    amount: 0.21
    country: GB
  - name: rate 6
    amount: 0.06
    country: NL
```

- **`name`** - (Required) Name of the category
- **`key`** - (Required) Key of the category
- **`rates`** - List of [tax rates](#taxes)

### stores

Defines [commercetools stores](https://docs.commercetools.com/http-api-projects-stores).

Example:
=== "Managed"
    ```yaml
    stores:
      - key: mystore
        name:
          en-GB: my store
        distribution_channels:
          - EU-DIST
        supply_channels:
          - EU-SUPPL
    ```
=== "Not managed"
    ```yaml
    stores:
      - key: mystore
        managed: false
      - key: our-other-store
        managed: false
      - key: global-store
        managed: false
    ```
- **`name`** - (Required) Name of the store. Localized string [^1]
- **`key`** - (Required) Store key
- `languages` - List of languages
- `distribution_channels` - List of supply channel keys used for [product projection store filtering](https://docs.commercetools.com/http-api-projects-productProjections#prices-beta)
- `supply_channels` - List of supply channel keys used for [product projection store filtering](https://docs.commercetools.com/http-api-projects-productProjections#prices-beta)
- `managed` - Indicate whether this store should be managed by MACH or not. Default to `true`

### frontend

Example:
```yaml
frontend:
  create_credentials: true
  permission_scopes: [manage_my_profile, manage_my_orders, view_states, manage_my_shopping_lists, view_products, manage_my_payments, create_anonymous_token, view_project_settings]
```

- `create_credentials` - Defines if frontend API credentials must be created
  Defaults to `true`
- `permission_scopes` - List of [scopes](https://docs.commercetools.com/api/scopes) excluding the project key

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

- **`client_id`** - (Required) Overrides default `client_id` settings (site-specific)
- **`client_secret`** - (Required) Override default `client_secret` setting (site-specific)
- **`hub_id`** - (Required) Override default `hub_id` setting (site-specific)


## apollo_federation

Apollo Federation configuration.

Example:

```yaml
apollo_federation:
  api_key: service:mach-poc-123:Abc00kHbB89h
  graph: mach-poc-123
  graph_variant: current
```

- `api_key` - Apollo Studio API key
- `graph` - Graph name to use for this site
- `graph_variant` - Graph variant


## azure
Site-specific Azure settings.<br>
Can overwrite any value from the generic [Azure settings](#azure):

- `tenant_id` - Tenant to deploy in
- `service_object_ids`
- `frontdoor` - [Front-door](#frontdoor) settings
- `subscription_id` - Subscription to deploy in
- `region` - Azure region to deploy in
- `service_plans` - Map of custom [service plan configurations](#service_plans)

And adds the following exta attributes:

- `alert_group` - List of [Alert groups](#alert_group)
- `resource_group` - Name of an already existing resource group.

!!! warning
    Use `resource_group` with care.<br>
    By default, MACH will manage the site resource groups for you. If you add this option later, the managed resource group will get **deleted**.<br>
    So only use for new site definitions

### frontdoor

Example:
```yaml
frontdoor:
  resource_group: my-shared-rg
  ssl_key_vault:
    name: mysharedwekvcerts
    resource_group: my-shared-we-rg
    secret_name: my-domain-net
```

- **`dns_resource_group`** - (Required) Resource group name where the DNS zone can be found
- `suppress_changes` - Suppress changes to the Frontdoor instance. This is a temporary work-around for some issues in the Azure Terraform provider.
- `ssl_key_vault` - SSL certificate configuration when Frontdoor should use your own certificate

#### ssl_key_vault
- **`name`** - (Required) KeyVault name
- **`resource_group`** - (Required) KeyVault resource group
- **`secret_name`** - (Required) Certificate name


### service_plans

Map of service plan definitions if you want to define additional service plans your components should run on, or if you want to overwrite the default.

Example:
=== "Additional plan"
    ```yaml
    # Here we add an additional service plan 'premium'
    service_plans:
      premium:
        kind: "Linux"
        tier: "PremiumV2"
        size: "P2v2"
        capacity: 2
    ```
=== "Default overwrite"
    ```yaml
    # Here we configure the default service plan to run Premium
    # and also offer a service plan running Windows
    service_plans:
      default:
        kind: "Linux"
        tier: "PremiumV2"
        size: "P2v2"
      windows:
        kind: "Windows"
        tier: "PremiumV2"
        size: "P2v2"
    ```

- **`kind`** - (Required) The kind of the App Service Plan to create. `Windows`, `Linux`, `elastic` or `FunctionApp`.
- **`tier`** - (Required) Specifies the plan's pricing tier.
- **`size`** - (Required) Specifies the plan's instance size.
- `capacity` - Specifies the number of workers associated with this App Service Plan.
- `dedicated_resource_group` - Indicates of the service plan should run on a dedicated resource group. This might be useful when, due to Azure hosting restrictions, a service plan cannot run on the same resource group as an existing one. Defaults to `false`.
- `per_site_scaling` - Can Apps assigned to this App Service Plan be scaled independently? If set to `false` apps assigned to this plan will scale to all instances of the plan. Defaults to `false`.

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
  deploy_role_name: deploy
  extra_providers:
    - name: email
      region: eu-west-1
```

- **`account_id`** - (Required) AWS account ID for this site
- **`region`** - AWS region to deploy site in
- `deploy_role_name` - The [IAM role](./prerequisites/aws#iam-deploy-role) name needed for deployment
- `extra_providers`

## components

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
  service_plan: default
```

- `service_plan` - The service plan (defined in [`service_plans`](./global.md#service_plans)) to use for this component. Set this to `default` if you want to use the MACH-managed Consumption plan.

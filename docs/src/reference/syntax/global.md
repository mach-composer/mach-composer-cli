# global
All 'shared' configuration that applies to all sites.

- **`environment`** - (Required) [environment](#environment) Identifier for the environment. For example `development`, `test` or `production`.<br>
Is used to set the `environment` variable of any [Terraform component](../components/structure.md#terraform-component)
- **`terraform_config`** - (Required) [terraform_config](#terraform_config) block
- `cloud` - Either `azure` or `aws`
- `azure` - [Azure](#azure) block
- `sentry` - [Sentry](#sentry) block
- `contentful` - [Contentful](#contentful) block


## terraform_config
Terraform configuration block.

Can be used to configure the state backend and Terraform provider versions.

- **`azure_remote_state`** - [Azure](#azure_remote_state) remote state configuration
- **`aws_remote_state`** - [AWS](#aws_remote_state) remote state configuration
- `providers` - [Providers](#providers) configuration block

### azure_remote_state
An Azure state backend can be defined as:

```yaml
azure_remote_state:
  resource_group: <your resource group>
  storage_account: <storage account name>
  container_name: <container name>
  state_folder: <state folder>
```

!!! tip ""
    A good convention is to give the state_folder the same name as [`environment`](#environment)

### aws_remote_state
An AWS S3 state backend can be defined as:

```yaml
aws_remote_state:
  bucket: mach-statefiles
  key_prefix: test-statefiles
  role_arn: arn:aws:iam::1234567890:roldeploy
```

- **`bucket`** - (Required) S3 bucket name
- **`key_prefix`** - (Required) Key prefix for each individual Terraform state
- `role_arn` - Role ARN to access S3 bucket with
- `lock_table` - DynamoDB lock table
- `encrypt` - Enable server side encryption of the state file. Defaults to `True`

### providers

Can be used to overwrite the MACH defaults for the Terraform provider versions.

Example:

```yaml
providers:
  aws: 3.21.0
```

- `aws` [aws provider](https://registry.terraform.io/providers/hashicorp/aws) version overwrite
- `azure` [azurerm provider](https://registry.terraform.io/providers/hashicorp/azurerm) version overwrite
- `commercetools` [commercetools provider](https://registry.terraform.io/providers/labd/commercetools) version overwrite
- `sentry` [sentry provider](https://registry.terraform.io/providers/jianyuan/sentry) version overwrite
- `contentful` [contentful provider](https://registry.terraform.io/providers/labd/contentful) version overwrite


!!! tip "Cache your providers"
    If you're overwriting the provider versions, make sure you [mount the plugins cache](../../topics/deployment/config/index.md#cache-terraform-providers)

## sentry
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
    The Sentry settings can be overwritten on [site](./sites.md#sentry) and [component](./sites.md#sentry_1) level

## azure

General Azure settings. Values can be overwritten [per site](./sites.md#azure).

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
  frontdoor:
    dns_resource_group: my-shared-rg
  service_plans:
    premium:
      kind: "Linux"
      tier: "PremiumV2"
      size: "P2v2"
```

- **`tenant_id`** - (Required) Tenant to deploy in
- **`subscription_id`** - (Required) Subscription to deploy in
- **`region`** - (Required) Azure region to deploy in
- `resources_prefix` - Prefix to be used for all Azure resources, for example `my-`
- `frontdoor` - [Front-door](#frontdoor) settings
- `service_object_ids` - Map of objects IDs that should have access to things like KeyVaults created for components
- `service_plans` - Map of custom [service plan configurations](#service_plans)


!!! tip "Don't lock yourself out"
    Make sure that, as a minimum, you set the `service_object_ids` to the objects IDs of the users or groups that perform the `mach apply`

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

### frontdoor

Example:
```yaml
frontdoor:
  dns_resource_group: my-shared-rg
```

- **`dns_resource_group`** - (Required) Resource group name where the DNS zone can be found
- `suppress_changes` - Suppress changes to the Frontdoor instance. This is a temporary work-around for some issues in the Azure Terraform provider.

## contentful
Defines global Contentful credentials to manage the spaces

- **cma_token** - (Required)
- **organization_id** - (Required)


## amplience
Defines global Amplience credentials to manage hubs

- **client_id** - (Required)
- **client_secret** - (Required)

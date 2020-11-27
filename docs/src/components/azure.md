# Azure components

Unless a component is flagged as `is_software_component: false`, a component within a Azure-based MACH configuration is considered an *Azure component*.

To be able to create the resources needed, a couple of extra [Terraform variables](#terraform-variables) are set by MACH.

In addition to this, the component itself is responsible for [packaging and deploying](#packaging-and-deploying) the correct assets in case of a Function App.

## Terraform variables

In addition to the [base variables](./index.md#required-variables), an Azure component expects the following:

- `short_name` - Short name; will not be more than 10 characters to prevent Azure naming limits
- `name_prefix` - Name prefix to be used for all Azure resources. See [naming conventions](#nmaing-conventions)
- `subscription_id` - The current subscription ID
- `tenant_id` - The current tenant ID
- `service_object_ids` - User/Group/Principle IDs that should be able to access for example a KeyVault.<br>
Type: `map(string)`
- `region` - Azure region
- `resource_group_name` - Name of the resource group the component should be created in
- `resource_group_location` - The resource group location
- `app_service_plan_id` - The [App service plan](../deployment/config/azure.md#app-service-plan) managed by MACH
- `tags` - Azure tags to be used on resources<br>
  Type: `map(string)`
- `monitor_action_group_id` - [action group](../deployment/config/azure.md#action-groups) ID when [alert_group](../syntax.md#azure) is configured.

```terraform
variable "short_name" {}
variable "name_prefix" {}
variable "subscription_id" {}
variable "tenant_id" {}
variable "service_object_ids" {
  type        = map(string)
  default     = {}
}
variable "region" {}
variable "resource_group_name" {}
variable "resource_group_location" {}
variable "app_service_plan_id" {}
variable "tags" {
  type        = map(string)
}
variable "monitor_action_group_id" {
  type        = string
  default     = ""
}
```


### With `endpoint`

In order to support the [`endpoint`](../deployment/config/azure.md#http-routing) attribute on the component, the component needs to have the following variables defined:

- `frontdoor_id` - Frontdoor ID on which routing is created for that component


## Packaging and deploying

For Azure functions, the deployment process constist of two steps:

- Packaging the function
- Deploying it to the [function app storage](../prerequisites/azure.md#create-function-app-storage)

[Read more](../deployment/components.md) about Azure component deployments.

### Configure runtime
When defining your Azure function app resource, you can reference back to the asset that is deployed:

```terraform
data "azurerm_storage_account" "shared" {
  name                = "mysharedwesacomponents"
  resource_group_name = "my-shared-we-rg"
}

locals {
  package_name = format("yourcomponent-%s.zip", var.component_version)
}

resource "azurerm_function_app" "your_component" {
    app_settings = {
        WEBSITE_RUN_FROM_ZIP = "https://${data.azurerm_storage_account.shared.name}.blob.core.windows.net/code/${local.package_name}${data.azurerm_storage_account_blob_container_sas.code_access.sas}"
        ...
    }
}
```
## HTTP routing

MACH will provide the correct HTTP routing for you.<br>
To do so, the following has to be configured:

- [Frontdoor](../syntax.md#front_door) settings in the Azure configuration
- The component needs to have an [`endpoint`](../syntax.md#components) defined

More information in the [deployment section](../deployment/config/azure.md#http-routing).

## Naming conventions

When setting up Terraform components for Azure we need to follow the following naming conventions and mind the current naming restrictions.

- [Azure best practices on naming things](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
- [Azure naming restrictions](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules)

An important thing to highlight is the following:

> Resource names have length limits. Balancing the context embedded in a name with its scope and length is important when you develop your naming conventions. For more information about naming rules for allowed characters, scopes, and name lengths for resource types, see [Naming conventions for Azure resources](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules).


### Resource prefixing

MACH creates a name prefix which can be used to name all other resources.

This prefix is built up using

- The configured [resources prefix](../syntax.md#azure)
- The [site identifier](../#syntax.md#sites)
- Region

So for example, when creating a function app, we can define this as:

```
resource "azurerm_function_app" "example_component" {
  name = lower(format("%s-func-%s", var.name_prefix, var.functionapp_name))
}
```

With `resources_prefix` configured as `mach-` and a site id of `uk-example-prd` (UK site of Example site on production) the result will be `mach-uk-example-p-we-func-example`.

!!! note ""
    To save character space [^1], the following adjustments are made:

    - The values `tst`, `dev` and `prd` will be shortened to `t`, `d`, and `p`
    - The region is shortened to two characters. For example: `westeurope` becomes `we`


### Azure naming restrictions

Some resources might have a maximum length to be used as name.<br>
To work with these limitations, the `short_name` variable can be used.

For example, a Storage Account can be created using

```terraform
resource "azurerm_storage_account" "main" {
    name = replace(lower(format("%s-sa-%s", var.name_prefix, var.short_name)), "-", "")
    ...
}
```

Where `sa` stands for *'storage account'*

[^1]: Some resources have a name restriction of max 24 characters. Obviously we want to avoid hitting that limit. See [Azure naming restrictions](#azure-naming-restrictions) on how to avoid that.

## Function App

We recommend using the [`azurerm_function_app`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/function_app)

```terraform
resource "azurerm_function_app" "example_component" {
  name                       = lower(format("%s-func-%s", var.name_prefix, var.short_name))
  location                   = var.resource_group_location
  resource_group_name        = var.resource_group_name
  app_service_plan_id        = var.app_service_plan_id
  storage_account_name       = azurerm_storage_account.main.name
  storage_account_access_key = azurerm_storage_account.main.primary_access_key
  app_settings               = local.environment_variables
  os_type                    = "linux"
  version                    = "~3"
  https_only                 = true

  site_config {
    linux_fx_version = "PYTHON|3.8"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = var.tags
}
```

See also [notes on using the serverless framework](../deployment/config/components.md#serverless-framework)
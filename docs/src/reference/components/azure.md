# Azure components

All components within a Azure-based MACH composer configuration are
automatically considered to have a 'azure' integration by default. Only if
'azure' is explicitely omitted from the `integrations` definition, it won't
require any Azure-specific variables.

To be able to create the resources needed, a couple of extra
[Terraform variables](#terraform-variables) are set by MACH composer.

In addition to this, the component itself is responsible for
[packaging and deploying](#packaging-and-deploying) the correct assets in case
of a Function App.

## Terraform variables

In addition to the [base variables](./structure.md#required-variables), an Azure
component expects the following:

```terraform
variable "azure_short_name" {
  type        = string
  description = "Short name; will not be more than 10 characters to prevent Azure naming limits"
}

variable "azure_name_prefix" {
  type        = string
  description = "Name prefix to be used for all Azure resources"
}

variable "azure_subscription_id" {
  type        = string
  description = "The current subscription ID"
}

variable "azure_tenant_id" {
  type        = string
  description = "The current tenant ID"
}

variable "azure_service_object_ids" {
  type        = map(string)
  description = "User/Group/Principle IDs that should be able to access for example a KeyVault."
}

variable "azure_region" {
  type        = string
  description = "Azure region"
}

variable "azure_resource_group" {
  type = object({
    name     = string
    location = string
  })
  description = "Information of the resource group the component should be created in"
}

variable "azure_monitor_action_group_id" {
  type        = string
  default     = "Action group ID when alert group is configured.
}
```

!!! info "Monitor action group"
    `monitor_action_group_id` is set to the
    [action group](../../topics/deployment/config/azure.md#action-groups)
    ID when a [alert_group](../syntax/global.md#azure) is configured.

### With `endpoints`

In order to support the [`endpoints`](../../topics/deployment/config/azure.md#http-routing)
attribute on the component, the component needs to define what endpoints it
expects **and** needs to provide output about the
[routing options](../../topics/deployment/config/azure.md#http-routing).

For example, if the component requires two endpoints (`main` and `webhooks`) to
be set, the following variables and outputs needs to be defined:

```terraform
variable "azure_endpoint_main" {
  type = object({
    url          = string
    frontdoor_id = string
  })
}

variable "azure_endpoint_webhooks" {
  type = object({
    url          = string
    frontdoor_id = string
  })
}

output "azure_endpoint_main" {
  value = {
    address = azurerm_app_service.main.default_site_hostname
  }
}

output "azure_endpoint_webhooks" {
  value = {
    address = azurerm_app_service.main.default_site_hostname
    routes = [
      {
        name = "hooks"
        patterns = [
          "/hooks/*",
        ]
      },
    ]
  }
}
```

#### Defining `outputs`

As shown in the example above, the component needs to have an output *per
endpoint* defined in order to instruct MACH how to setup Frontdoor routing.

This output needs to have a name in the form of `azure_endpoint_<endpoint-name>` and contain the following attributes:

- **`address`** - (Required) The host address to route traffic to
- `host_header` - The value to use as the host header sent to the backend. By default will take the value of `address`.
- `http_port` -  The HTTP TCP port number. Possible values are between `1` - `65535`. Defaults to `80`.
- `https_port` - The HTTPS TCP port number. Possible values are between `1` - `65535`. Defaults to `443`.
- `health_probe_path` - The path to use for the Health Probe. If left empty, health probe won't be enabled.
- `health_probe_protocol` - Protocol scheme to use for the Health Probe. Defaults to `Http`
- `health_probe_method` - Specifies HTTP method the health probe uses when querying the service. Possible values include: `GET` and `HEAD`. Defaults to `GET`.
- `routes` - A list of custom Frontdoor routing rules. By default MACH will generate one default routing rule for each component.

**`routes` options**

- `patterns` - List of patterns to match
- `cache_enabled` - Specifies whether to Enable caching or not. Valid options are `true` or `false`. Defaults to `false`.
- `custom_forwarding_path` - Path to use when constructing the request to
  forward to the backend. This functions as a URL Rewrite. Default behaviour
  preserves the URL path.


**Full example**
```terraform
output "azure_endpoint_main" {
  value = {
    address = azurerm_app_service.main.default_site_hostname
    health_probe_path = "/"
    health_probe_method = "Get"
    host_header = "www.example.com/something-else/"
    https_port = 9000
    routes = [
      {
        patterns = [
          "/*",
        ]
      },
      {
        patterns = [
          "/media/*",
        ]
        cache_enabled = true
      }
    ]
  }
}
```

### With `service_plan`

When a component has been configured with a [`service_plan`](../syntax/sites.md#azure_1),
MACH composer manages the service plan for you and passes the information to the
component with a `app_service_plan` variable:

```terraform
variable "azure_app_service_plan" {
  type = object({
    id                  = string
    name                = string
    resource_group_name = string
  })
}
```

## Packaging and deploying

For Azure functions, the deployment process constist of two steps:

- Packaging the function
- Deploying it to the [function app storage](../../tutorial/azure/step-3-setup-azure.md)

[Read more](../../topics/deployment/components.md) about Azure component deployments.

### Configure runtime
When defining your Azure function app resource, you can reference back to the
asset that is deployed:

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

MACH composer will provide the correct HTTP routing for you.<br>
To do so, the following has to be configured:

- [Frontdoor](../syntax/global.md#frontdoor) settings in the Azure configuration
- The component needs to have [`endpoints`](../syntax/components.md) defined

!!! tip "Default endpoint"
    If you assign `default` to one of your components endpoints, no additional
    Frontdoor settings are needed.

    MACH composer will create a Frontdoor instance for you without any custom domain.

More information in the [deployment section](../../topics/deployment/config/azure.md#http-routing).

## Naming conventions

When setting up Terraform components for Azure we need to follow the following
naming conventions and mind the current naming restrictions.

- [Azure best practices on naming things](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
- [Azure naming restrictions](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules)

An important thing to highlight is the following:

> Resource names have length limits. Balancing the context embedded in a name
with its scope and length is important when you develop your naming conventions.
For more information about naming rules for allowed characters, scopes, and name
lengths for resource types, see [Naming conventions for Azure
resources](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules).


### Resource prefixing

MACH composer creates a name prefix which can be used to name all other resources.

This prefix is built up using

- The configured [resources prefix](../syntax/global.md#azure)
- The [site identifier](../#syntax/sites.md)
- Region

So for example, when creating a function app, we can define this as:

```
resource "azurerm_function_app" "example_component" {
  name = lower(format("%s-func-%s", var.azure_name_prefix, var.azure_short_name))
}
```

With `resources_prefix` configured as `mach-` and a site id of `uk-example-prd`
(UK site of Example site on production) the result will be
`mach-uk-example-p-we-func-example`.

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
    name = replace(lower(format("%s-sa-%s", var.azure_name_prefix, var.azure_short_name)), "-", "")
    ...
}
```

Where `sa` stands for *'storage account'*

[^1]: Some resources have a name restriction of max 24 characters. Obviously we
want to avoid hitting that limit. See
[Azure naming restrictions](#azure-naming-restrictions) on how to avoid that.

## Function App

We recommend using the [`azurerm_function_app`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/function_app)

```terraform
resource "azurerm_function_app" "example_component" {
  name                       = lower(format("%s-func-%s", var.azure_name_prefix, var.azure_short_name))
  location                   = var.azure_resource_group.location
  resource_group_name        = var.azure_resource_group.name
  app_service_plan_id        = var.azure_app_service_plan.id
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

See also [notes on using the serverless framework](../../topics/deployment/config/components.md#serverless-framework)

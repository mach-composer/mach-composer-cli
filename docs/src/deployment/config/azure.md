# Azure deployments

## Resource groups

MACH will create a **[resource group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group) per site**.

!!! info ""
    Only when a [`resource_group`](../../syntax.md#azure_1) is explicitly set, it won't be managed by MACH.

## HTTP routing

Only when a MACH stack contains components that are configures as [`has_public_api`](../../syntax.md#components), MACH will setup the necessary resources to be able to route traffic to that component:

- Frontdoor instance
- DNS record

It will use the information from the [`front_door` configuration](../../syntax.md#front_door) to setup the Frontdoor instance.

The information needed for components to add custom routes to that API Gateway are provided through [Terraform variables](../../components/azure.md#terraform-variables).

## App service plan

MACH will create an [App service plan](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_plan) that can be used for any MACH component that implements an Azure function.

## Action groups

When an [Alert group](../../syntax.md#alert_group) is configured, an [Action group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_action_group) will be created.

Components can use that action group to attach alert rules to.
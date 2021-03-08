# Azure deployments

## Resource groups

MACH will create a **[resource group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group) per site**.

!!! info ""
    Only when a [`resource_group`](../../../reference/syntax/sites.md#azure) is explicitly set, it won't be managed by MACH.

## HTTP routing

Only when a MACH stack contains components that have an [`endpoint`](../../../reference/syntax/components.md) defined, MACH will setup a **Frontdoor instance** to be able to route traffic to that component.

### Default endpoint

If you have defined your component with a `default` endpoint, MACH will create a Frontdoor instance for you which includes the default Azure domain.

```
components:
  - name: payment
    source: git::ssh://git@github.com/your-project/components/payment-component.git//terraform
    endpoints:
      public: default
    version: ....
```

!!! note ""
    This `default` endpoint doesn't need to be defined in your [endpoints definition](../../../reference/syntax/sites.md#endpoints).

### Custom endpoint

Whenever a custom endpoint from your [endpoints definition](../../../reference/syntax/sites.md#endpoints) is used, MACH will require that you have configured [`frontdoor`](../../../reference/syntax/global.md#frontdoor) for additional DNS information that it needs to setup your Frontdoor instance.

In addition to that it will also setup the necessary DNS record.

### Routes to the component

For each component with an `endpoint` MACH composer will add a route to the Frontdoor instance using the name of the component.

So when having the following components defined:

```yaml
components:
  - name: payment
    source: git::ssh://git@github.com/your-project/components/payment-component.git//terraform
    endpoints: 
      public: main
    version: ....
  - name: api-extensions
    source: git::ssh://git@github.com/your-project/components/api-extensions-component.git//terraform
    version: ....
  - name: graphql
    source: git::ssh://git@github.com/your-project/components/graphql-component.git//terraform
    endpoints: 
      public: main
    version: ....
```

The routing in Frontdoor that will be created:

![Frontdoor routes](../../../_img/azure/frontdoor_routes.png)

## App service plans

MACH can create an [App service plan](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_plan) that can be used for any MACH component that uses the `service_plan` configuration option (either in the [component definition](../../../reference/syntax/components.md#azure) or the [site-specific configuration](../../../reference/syntax/sites.md#azure_1)).

What kind and how many service plans it will create depends on what service plans are needed by the components, and how the [`service_plans`](../../../reference/syntax/global.md#service_plans) configuration looks like.

By default, MACH will create a `default` service plan which is a Consumption plan, if any of the components have `service_plan: default` set in their configuration.

In this case, you don't need to define the service plan yourself, MACH automatically creates this default for you:

```yaml
service_plans:
  default:
    kind: "FunctionApp"
    tier: "Dynamic"
    size: "Y1"
```

## Action groups

When an [Alert group](../../../reference/syntax/sites.md#alert_group) is configured, an [Action group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_action_group) will be created.

Components can use that action group to attach alert rules to.

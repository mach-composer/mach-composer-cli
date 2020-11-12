# commercetools

## Configuration

From a MACH configuration file, you can configure the following items in commercetools:

- currencies
- languages
- countries
- channels
- taxes
- stores

For more information about these configuration options, see the [syntax](../syntax.md#commercetools).

!!! tip "More fine-grained control"
    MACH provides a couple of basic configuration options which in most cases are sufficient and is setup very quickly.  
    If you need more control over the configuration in your project, you can load custom configuration through a new component which would include the necessary Terraform configuration for that.  
    For example a `commercetools-setup-component`.

## API client management

MACH is provided by admin credentials which allows MACH to fully manage your commercetools project.

Each component that implements functionality that needs to communicate with commercetools needs their own set of credentials.
The component is responsible for creating the necessary API client credentials for that specific component, with only the most necessary scopes set.

One way a component could facilitate in this would be to define the following in your Terraform config:

```terraform
locals {
  ct_scopes = formatlist("%s:%s", [
    "manage_orders",
  ], var.ct_project_key)
}

resource "commercetools_api_client" "client" {
  name  = format("%s_api_extension", var.name_prefix)
  scope = local.ct_scopes
}
```
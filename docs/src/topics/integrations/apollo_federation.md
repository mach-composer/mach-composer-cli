# Apollo Federation

You can define a site specific configuration for [managed Apollo federation](https://go.apollo.dev/s/managed-federation).

## Retrieve credentials

Go to Apollo Studio and you will be prompted with a popup for the credentials.

The *graph* setting can be determined from the **APOLLO_KEY** string, use the name between the two `:` characters.

## Example site configuration block

An example site block:

```yaml
apollo_federation:
  api_key: service:mach-poc-123:Abc00kHbB89h
  graph: mach-poc-123
  graph_variant: current
```
      
## Expose Apollo Federation to components

MACH needs to know what components want to use the Apollo Federation configuration.<br>
For this you need to include `apollo_federation` to the list of integrations.

Components with the `apollo_federation` integration need to have an extra required variable [^1]:

```terraform
variable "apollo_federation" {
  type = object({
    api_key       = string
    graph         = string
    graph_variant = string
  })
}
```

## Further reading

For more instructions on how to integrate Apollo Federation into your MACH stack, [read our how-to](../../howto/apollo-federation.md).

[^1]: Read more about the [Apollo Federation structure for components](../../reference/components/structure.md#apollo-federation)

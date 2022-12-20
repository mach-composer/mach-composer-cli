# Apollo Federation

You can define a site specific configuration for [managed Apollo federation](https://go.apollo.dev/s/managed-federation).

## Retrieve credentials

Go to Apollo Studio and you will be prompted with a popup for the credentials.

The *graph* setting can be determined from the **APOLLO_KEY** string, use the
name between the two `:` characters.

## Example site configuration block

An example site block:

```yaml
apollo_federation:
  api_key: service:mach-poc-123:Abc00kHbB89h
  graph: mach-poc-123
  graph_variant: current
```

## Integrate with components

When `apollo_federation` is set as an [component integration](../../reference/components/structure.md#integrations),
the component should have the following Terraform variables defined:

- `apollo_federation`

!!! info ""
    More information on the [apollo_federation integration on components](../../reference/components/structure.md#apollo-federation).

## Further reading

For more instructions on how to integrate Apollo Federation into your MACH composer
stack, [read our how-to](../../howto/apollo-federation.md).

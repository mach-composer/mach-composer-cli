# Apollo Federation

You can define a site specific configuration for [managed Apollo federation](https://go.apollo.dev/s/managed-federation).

## Retrieve credentials

Go to Apollo Studio and you will be prompted with a popup for the credentials.

The *graph* setting can be determined from the **APOLLO_KEY** string, use the name between the two `:` characters.

## Example configuration block

An example site block:

```yaml
apollo_federation:
  api_key: service:mach-poc-123:Abc00kHbB89h
  graph: mach-poc-123
  graph_variant: current
```
      
## Expose Apollo Federation to components

MACH needs to know what components want to use the Apollo Federation configuration.<br>
For this you need to include `apollo_federation` to the list of integrations.<br>
When doing so, MACH expects the component to have one variable `apollo_federation` defined ([more info](../../reference/components/structure.md#apollo federation))

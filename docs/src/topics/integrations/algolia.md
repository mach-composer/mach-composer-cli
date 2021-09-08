# Algolia

You can define a site specific configuration for [Algolia Search](https://www.algolia.com/doc/) which uses the unofficial [Algolia TF provider](https://registry.terraform.io/providers/k-yomo/algolia/latest).

## Retrieve credentials

Login to Algolia and go to API keys to retrieve the `Admin API Key`

## Example site configuration block

An example site block:

```yaml
algolia:
  application_id: ab12345
  api_key: abcd12345
```
      
## Integrate with components

When `algolia` is set as an [component integration](../../reference/components/structure.md#integrations), the component should have the following Terraform variables defined:

- `algolia`

!!! info ""
    More information on the [algolia integration on components](../../reference/components/structure.md#algolia).

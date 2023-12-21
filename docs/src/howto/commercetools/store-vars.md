# Parse store variables and secrets

When [integrating with commercetools](../../plugins/commercetools.md)
MACH composer makes it possible to define
[store specific variables and secrets](../../plugins/commercetools.md).

You could have for example the following configuration

```yaml
component:
  - my-component:
    variables:
      FROM_EMAIL: mach@example.com
    store_variables:
      uk-store:
        FROM_EMAIL: mach@example.co.uk
      nl-store:
        FROM_EMAIL: mach@example.nl
    store_secrets:
      uk-store:
        MAIL_API_KEY: 5kg6Z4HxgLMfBjYj5BOg
      nl-store:
        MAIL_API_KEY: i1hajIJ92LPNYGB2p3W1
```

## Incoming format

The above example will result in the following variable definition:

```terraform
ct_stores = {
  uk-store = {
    key = "UK"
    variables = {
      FROM_EMAIL = "mach@example.co.uk"
    }
    secrets = {
      MAIL_API_KEY = "5kg6Z4HxgLMfBjYj5BOg"
    }
  }
  nl-store = {
    key = "NL"
    variables = {
      FROM_EMAIL = "mach@example.nl"
    }
    secrets = {
      MAIL_API_KEY = "i1hajIJ92LPNYGB2p3W1"
    }
  }
}
```

## Using store variables

One way of approaching the store variables is to provide them all in your
runtime environment prefixed with the store key.

In your function, depending on the Store context, you can choose what
environment setting to use.

In this example, we'd like to set the following environment variables on our
function runtime:

```bash
UK_FROM_EMAIL = "mach@example.co.uk"
NL_FROM_EMAIL = "mach@example.nl"
```

This can be done with the following Terraform definition:

```terraform
locals {
  store_variables = flatten([
    for store in values(var.ct_stores) : [
      for variable_key, variable_value in store.variables : {
        "${store.key}_${variable_key}" : variable_value
      }
    ]
  ])
  env_store_variables = zipmap(
    flatten([for item in local.store_variables : keys(item)]),
    flatten([for item in local.store_variables : values(item)])
  )
}
```

## Using store secrets

For the store secrets you can use the same technique as for the [store variables](#using-store-variables).
This especially will be sufficient for most cases when implementing for Azure.

On AWS, it might be a better option to combine secrets per store into one secrets
value to avoid too much latency when fetching those secrets.

### AWS

#### Combined secrets
```terraform
resource "aws_secretsmanager_secret" "store_secret" {
  for_each = var.ct_stores
  name     = "my-component/${each.value.key}-secrets"

  tags = {
    lambda = "my-component"
  }
}

resource "aws_secretsmanager_secret_version" "store_secret" {
  for_each      = var.ct_stores
  secret_id     = aws_secretsmanager_secret.component_secret[each.key].id
  secret_string = jsonencode(each.value.secrets)
}
```

!!! tip "Combining secrets"
    Depending on your use-case you could choose to store **all** secrets into
    one AWS secret or to have a logical split.

    More considerations about this in the [
      'Handle secrets in components'](../security/handle-secrets.md#combine-or-split-up) how-to.

#### Separate secrets
```terraform
locals {
  store_secrets = flatten([
    for store in values(var.ct_stores) : [
      for variable_key, variable_value in store.secrets : {
        "${store.key}_${variable_key}" : variable_value
      }
    ]
  ])
  env_store_secrets = zipmap(
    flatten([for item in local.store_secrets : keys(item)]),
    flatten([for item in local.store_secrets : values(item)]),
  )
  secrets = merge(local.env_store_secrets, var.secrets)
}

resource "aws_secretsmanager_secret" "component_secret" {
  for_each = local.secrets
  name     = "my-component/${replace(each.key, "_", "-")}-secret-${random_id.main.hex}"

  tags = {
    lambda = "my-component"
  }
}

resource "aws_secretsmanager_secret_version" "component_secret" {
  for_each      = local.secrets
  secret_id     = aws_secretsmanager_secret.component_secret[each.key].id
  secret_string = each.value
}
```

### Azure

```terraform
locals {
  store_secrets = flatten([
    for store in values(var.ct_stores) : [
      for variable_key, variable_value in store.secrets : {
        "${store.key}_${variable_key}" : variable_value
      }
    ]
  ])
  env_store_secrets = zipmap(
    flatten([for item in local.store_secrets : keys(item)]),
    flatten([for item in local.store_secrets : values(item)]),
  )
  secrets = merge(local.env_store_secrets, var.secrets)
}

# Key Vault definitions
# ...

resource "azurerm_key_vault_secret" "secrets" {
  for_each = local.secrets

  name         = replace(each.key, "_", "-")
  value        = each.value
  key_vault_id = azurerm_key_vault.main.id
  tags         = var.tags
}
```

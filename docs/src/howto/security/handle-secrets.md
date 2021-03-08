# Handle secrets in components

Secrets values are passed on to the component by MACH using the [`secrets` attribute](../../reference/syntax/sites.md#components).

It is up to the component to use those secret values in a secure manner.

## Handle secrets in AWS

In AWS we recommend storing the secret values in the [Secrets Manager](https://aws.amazon.com/secrets-manager/) and provide references to those secrets in the Lambda environment variables

### Store secrets

```terraform
# Having a random_id in the secrets name avoids issues when a component gets removed
# and added again. No need to import secrets into the state and recover them
resource "random_id" "main" {
  byte_length = 5
  keepers = {
    # Generate a new id each time set of secrets change
    secrets = join("", tolist(keys(var.secrets)))
  }
}

resource "aws_secretsmanager_secret" "component_secret" {
  for_each = var.secrets
  name     = "my-component/${replace(each.key, "_", "-")}-secret-${random_id.main.hex}"

  tags = {
    lambda = "my-component"
  }
}

resource "aws_secretsmanager_secret_version" "component_secret" {
  for_each      = var.secrets
  secret_id     = aws_secretsmanager_secret.component_secret[each.key].id
  secret_string = each.value
}
```

### Reference secret values

Make sure your Lambda function knows where to find the secrets.

By providing the references to the Secret Manager secrets the Lambda function can use the AWS SDK to retreive the values.

```terraform
locals {
  secret_references = {
    for key in keys(var.secrets) : "${key}_SECRET_NAME" => aws_secretsmanager_secret.component_secret[key].name
  }
}

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  ...

  environment_variables = merge(
    local.secret_references,
    {
      COMPONENT_NAME = "my-component",
      ENVIRONMENT    = var.environment,
      ...
    }
  ) 
```

### Configure your Lambda IAM policy

Make sure your Lambda has the correct policies to access the secrets **but only the secrets of that component**.

One way of achieving this is to use the `tags` that we've set on the secret itself.

Snippet of our IAM policy that we assign to the Lambda:
```terraform
# Secrets manager
statement {
  actions = [
    "secretsmanager:GetSecretValue",
  ]

  resources = [
    "*",
  ]

  condition {
    test     = "StringEquals"
    variable = "secretsmanager:ResourceTag/lambda"
    values   = ["my-component"]
  }
}
```

## Handle secrets in Azure

In Azure we can store the secrets in a KeyVault and pass the KeyVaul references to the Function App that needs those values.

### Store the secrets

```terraform
resource "azurerm_key_vault" "main" {
  name                        = replace("${var.azure_name_prefix}-kv-${var.azure_short_name}", "-", "")
  location                    = var.azure_resource_group.location
  resource_group_name         = var.azure_resource_group.name
  tenant_id                   = var.azure_tenant_id
  enabled_for_disk_encryption = true
  sku_name                    = "standard"
  soft_delete_enabled         = true
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  tags = var.tags
}

resource "azurerm_key_vault_secret" "secrets" {
    for_each     = var.secrets

    name         = replace(each.key, "_", "-")
    value        = each.value
    key_vault_id = azurerm_key_vault.main.id
}
```

### Reference secret values

```terraform
locals {
  secret_variables = {
    for k, v in azurerm_key_vault_secret.secrets : replace(k, "-", "_") => "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault.main.vault_uri}secrets/${v.name}/)" }
}


resource "azurerm_function_app" "crm_365_component" {
  ...

  app_settings = merge(
    local.secret_variables,
    {
      COMPONENT_NAME = "my-component",
      ENVIRONMENT    = var.environment,
      ...
    }
  )
  depends_on   = [azurerm_key_vault_secret.secrets]
}
```

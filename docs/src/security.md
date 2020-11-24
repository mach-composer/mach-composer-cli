# Security

## Component secrets
The MACH configuration provides a [`secrets` attribute](./syntax.md#component-configurations) in which you can pass secret values to the components.

It is up to the component to use those secret values in a secure manner.

=== "Azure"
    ```terraform
    resource "azurerm_key_vault_secret" "secrets" {
        for_each     = var.secrets

        name         = replace(each.key, "_", "-")
        value        = each.value
        key_vault_id = azurerm_key_vault.main.id
    }
    ```
=== "AWS"
    ```terraform
    resource "aws_secretsmanager_secret" "mail_client_secret" {
        name = "my_component/mail-client-secret"
    }

    resource "aws_secretsmanager_secret_version" "mail_client_secret" {
        secret_id     = aws_secretsmanager_secret.mail_client_secret.id
        secret_string = var.secrets["MAIL_CLIENT_SECRET"]
    }
    ```

## Encrypt your MACH configuration

A MACH configuration typically contains secrets that are configured on the components as well as secrets used to configure the integrations.

We recommend using [SOPS](https://github.com/mozilla/sops) to encrypt your MACH configuration files or a part of it.

#### Encrypting

Encrypting your file can be done with the `sops --encrypt` command:

=== "AWS"
    ```bash
    $ export SOPS_KMS_ARN="arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500"
    $ export SOPS_PGP_FP="C9CAB0AF1165060DB58D6D6B2653B624D620786D"
    $ sops -e main.yml > main.enc.yml
    $ mv main.enc.yml main.yml
    ```
=== "Azure"
    ```bash
    $ sops --encrypt --azure-kv https://yoursharedsops.vault.azure.net/keys/sops-key/<your-key> main.yml > main.enc.yml
    $ mv main.enc.yml main.yml
    ```

#### Decrypt using deployment

In order to make this work with a MACH deployment you'll need to add an extra step to your CI/CD process:

```bash
$ sops -d main.yml --output-type=yaml
$ mach apply -f main.yml.dec
```
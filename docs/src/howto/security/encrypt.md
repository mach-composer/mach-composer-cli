# Encrypt MACH configuration

A MACH configuration typically contains secrets that are configured on the components as well as secrets used to configure the integrations.

We recommend using [SOPS](https://github.com/mozilla/sops) to encrypt your MACH configuration files or a part of it.

#### Encrypting

Encrypting your file can be done with the `sops --encrypt` command:

=== "AWS"
    ```bash
    $ export SOPS_KMS_ARN="arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500"
    $ sops -e --encrypted-regex '^(.*(secret|token).*)$' main.yml > main.enc.yml
    $ mv main.enc.yml main.yml
    ```
=== "Azure"
    ```bash
    $ export SOPS_AZURE_KEYVAULT_URLS="https://yoursharedsops.vault.azure.net/keys/sops-key/<your-key>"
    $ sops -e --encrypted-regex '^(.*(secret|token).*)$' main.yml > main.enc.yml
    $ mv main.enc.yml main.yml
    ```

#### Decrypt during deployment

In order to make this work with a MACH deployment you'll need to add an extra step to your CI/CD process:

```bash
$ sops -d main.yml --output-type=yaml > main.yml.dec
$ mach apply -f main.yml.dec
```

Make sure that your CI/CD environment has access to the appropriate encryption keys in AWS KMS or Azure KeyVault.

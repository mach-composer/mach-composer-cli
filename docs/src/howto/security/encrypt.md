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

MACH composer offers built-in support for decrypting sops-encrypted files automatically.

When MACH composer encounters an encrypted YAML file, it will attempt to decrypt the file prior to the execution of `generate`, `plan` or `apply`. 
Make sure that your CI/CD environment has access to the appropriate encryption keys in AWS KMS or Azure KeyVault.

##### Decrypting manually
Manual decrypting of the configuration can be done as follows:

```bash
$ sops -d main.yml --output-type=yaml > main.yml.dec
```

And you can then execute MACH composer with this decrypted file:
```bash
$ mach apply -f main.yml.dec
```

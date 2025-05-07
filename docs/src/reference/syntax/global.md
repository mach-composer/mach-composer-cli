# Global

All 'shared' configuration that applies to all sites.

## Schema

### Required

- `environment` (String) Identifier for the environment. For example `development`, `test` or `production`. Is used to
  set the `environment` variable for terraform components
- `cloud` (String) Either `azure`, `aws` or `gcp`. Defines the cloud provider to use. This will be used to load the
  correct cloud specific configuration.
- `terraform_config` (Block) This configuration will determine what type of backend terraform will use to store its
  state. See [below for nested schema](#nested-schema-for-terraform_config)

### Optional
- `variables` (Map of String) Variables for this configuration. Note that variables with the same name set in the site
  configuration or site component configuration will override these values
- `secrets` (Map of String) Variables for this configuration that should be stored in an encrypted key-value store . Note that
  variables with the same name set in the site configuration or site component configuration will override these values

### Dynamic

{% include-markdown "./dynamic.md" %}

## Nested schema for `terraform_config`

Terraform configuration block. Can be used to configure the state backend
and Terraform provider versions.

### Optional

- `providers` (Map of String) Can be used to overwrite the MACH Plugin
  defaults for the Terraform provider versions. The format is
  `provider_name: version`. For example `aws: 3.0.0`. If left empty the
  plugin defaults will be used.
- `remote_state` (Block) [Remote state](#nested-schema-for-remote_state)
  configuration. If left empty local state will be used. It is recommended
  to use one of the supported cloud providers for remote state instead.

## Nested schema for `remote_state`

### Required

- `plugin` (String) The plugin to use. One of `aws`, `azure`, `gcp`, `http`, `local` or `terraform_cloud`.
  This will determine what remote state backend configs will be available

### Dynamic

Depending on the `plugin` value, the following blocks will be merged into
the `remote_state` block:

- `aws` (Block) [AWS](#nested-schema-for-aws) state configuration for AWS
  backend
- `azure` (Block) [Azure](#nested-schema-for-azure) state configuration for
  Azure backend
- `gcp` (Block) [GCP](#nested-schema-for-gcp) state configuration for GCP
  backend
- `http` (Block) [HTTP](#nested-schema-for-http) state configuration for
  HTTP backend
- `local` (Block) [Local](#nested-schema-for-local) state configuration for
  local backend
- `terraform_cloud` (Block) [Terraform Cloud](#nested-schema-for-terraform_cloud)
  state configuration for Terraform Cloud backend

## Nested schema for `aws`

An AWS S3 state backend can be configured with the following options

### Example

```yaml
remote_state:
  plugin: aws
  bucket: <your bucket>
  region: <your region>
  key_prefix: <your key prefix>
  role_arn: <your role arn>
```

### Required

- `bucket` (String) S3 bucket name
- `region` (String) AWS region
- `key_prefix` (String) Key prefix for each individual Terraform state

### Optional

- `role_arn` - Role ARN to access S3 bucket with
- `lock_table` - DynamoDB lock table
- `encrypt` - Enable server side encryption of the state file. Defaults
  to `True`

## Nested schema for `azure`

An Azure state backend can be configured with the following options

### Example

```yaml
remote_state:
  plugin: azure
  resource_group: <your resource group>
  storage_account: <storage account name>
  container_name: <container name>
  state_folder: <state folder>
```

!!! tip ""
A good convention is to give the state_folder the same name
as the environment

## Required

- `resource_group` (String) Resource group name
- `storage_account` (String) Storage account name
- `container_name` (String) Container name

### Optional

- `state_folder` (String) Folder name for each individual Terraform state.
  If left empty the site identifier will be used

## Nested schema for `gcp`

A GCP state backend can be configured with the following options

### Example

```yaml
remote_state:
  plugin: gcp
  bucket: <your bucket>
  prefix: <your prefix>
```

### Required

- `bucket` (String) GCS bucket name
- `prefix` (String) Prefix for each individual Terraform state

## Nested schema for `http`

An HTTP state backend can be configured with the following options.

### Example

```yaml
remote_state:
  plugin: http
  address: https://example.com/state
```

### Required

- `address` (String) The address of the REST endpoint.

### Optional

- `update_method` (String) HTTP method to use when updating state. Defaults to `POST`.
- `lock_address` (String) The address of the lock REST endpoint. Defaults to disabled.
- `lock_method` (String) The HTTP method to use when locking. Defaults to `LOCK`.
- `unlock_address` (String) The address of the unlock REST endpoint. Defaults to disabled.
- `unlock_method` (String) The HTTP method to use when unlocking. Defaults to `UNLOCK`.
- `username` (String) The username for HTTP basic authentication.
- `password` (String) The password for HTTP basic authentication.
- `skip_cert_verification` (Boolean) Whether to skip TLS verification. Defaults to `false`.
- `retry_max` (Integer) The number of HTTP request retries. Defaults to `2`.
- `retry_wait_min` (Integer) The minimum time in seconds to wait between HTTP request attempts. Defaults to `1`.
- `retry_wait_max` (Integer) The maximum time in seconds to wait between HTTP request attempts. Defaults to `30`.
- `client_certificate_pem` (String) A PEM-encoded certificate used by the server to verify the client during mutual TLS (mTLS) authentication.
- `client_private_key_pem` (String) A PEM-encoded private key, required if `client_certificate_pem` is specified.
- `client_ca_certificate_pem` (String) A PEM-encoded CA certificate chain used by the client to verify server certificates during TLS authentication.

## Nested schema for `local`

A local state backend can be configured with the following options

### Example

```yaml
remote_state:
  plugin: local
  path: <your path>
```

### Optional

- `path` (String) Local path to store state files. Defaults to
  `./terraform.tfstate`

## Nested schema for `terraform_cloud`

A Terraform Cloud state backend can be configured with the following options.

### Example

```yaml
remote_state:
  plugin: terraform_cloud
  organization: <your organization>
```

### Required

- `organization` (String) The name of the Terraform Cloud organization.

### Optional

- `hostname` (String) The hostname of the Terraform Cloud instance. Defaults to `app.terraform.io`.
- `token` (String) The token used to authenticate with the Terraform Cloud backend. It is recommended to omit this field and use `terraform login` or configure credentials in the CLI config file instead.
- `workspaces` (Block) Configuration for workspaces:
  - `name` (String) The name of the workspace.
  - `prefix` (String) A prefix for dynamically created workspaces.
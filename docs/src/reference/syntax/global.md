# Global

All 'shared' configuration that applies to all sites.

## Schema

### Required

- `environment` (String) Identifier for the
  environment. For example `development`, `test` or `production`. Is used to set
  the `environment` variable for terraform components
- `cloud` (String) Either `azure`, `aws` or `gcp`. Defines the cloud provider to
  use. This will be used to load the correct cloud specific configuration.
- `terraform_config` (Block) This configuration will determine what type of
  backend terraform will use to store its state.
  See [below for nested schema](#nested-schema-for-terraform_config)

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

- `plugin` (String) The plugin to use. One of `aws`, `gcp`, `azure` or `local`.
  This will determine what remote state backend configs will be available

### Dynamic

Depending on the `plugin` value, the following blocks will be merged into
the `remote_state` block:

- `azure` (Block) [Azure](#nested-schema-for-azure) state configuration for
  Azure backend
- `aws` (Block) [AWS](#nested-schema-for-aws) state configuration for AWS
  backend
- `gcp` (Block) [GCP](#nested-schema-for-gcp) state configuration for GCP
  backend
- `local` (Block) [Local](#nested-schema-for-local) state configuration for
  local backend

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

## Nested schema for `local`

A GCP state backend can be configured with the following options

### Example

```yaml
remote_state:
  plugin: local
  path: <your path>
```

### Optional

- `path` (String) Local path to store state files. Defaults to
  `./terraform.tfstate`

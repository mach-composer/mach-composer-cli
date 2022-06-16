# Changelog

## 2.2.1 (2022-06-10)
- Fixed inconsistencies between 1.2 and 2.x:
  - Add `branch` option to component definitions to be able to perform a `mach
    update` and stay within a certain branch (during development)
  - Skip non-MACH configuration files when processing all yaml files in a directory.<br>
    This allows you to run things like `mach apply` or `mach update` without having to specify the `-f main.yml` option if you only have one valid MACH configuration file in your directory. Fixes #150

## 2.2 (2022-06-10)
- Fixed inconsistencies between 1.2 and 2.x:
  - Upgrade Terraform providers in golang version of the MACH composer to match the 1.2 release:
    - Upgraded commercetools provider to 0.30.0
    - Upgraded Amplience provider to 0.3.7
    - Upgraded Azure provider to 2.99.0
  - Add `variables_file` option to the `mach_composer` configuration block to define a variable file
  - Fix auto add cloud integration (aws or azure) when `integration` list is left empty
  - Add ability to define a custom provider version including the version operator
- Deprecate `commercetools.frontend` block, will be removed in a later release.


## 2.1.1 (2022-04-22)
- Don't crash when running `mach-composer apply` without `--auto-approve`


## 2.1.0 (2022-04-22)
- Add back support to update sops encrypted config files
- Properly implement the `--check` flag on `update` command


## 2.0.2 (2022-04-05)
- Pass environment variables to terraform command


## 2.0.1 (2022-04-05)
- Add aws-cli to the Docker container


## 2.0.0 (2022-04-05)
Rewrite of the Python codebase to Go. Goal is to make it easier to distribute
mach-composer in a cross-platform way.

A number of features which were minimal used are removed.
  - The `mach bootstrap` command is no longer present. It was a simple wrapper
    around Python cookiecutter. This can still be used separately
  - The `mach sites` and `mach components` commands since they were unused.
  - The `--with-sp-login` is removed. This flags used to run `az login`. If this
    is needed it needs to be run before mach-composer is run.
  - The `--ignore-version` flag is removed. The version in the config file now
    indicates a schema version. Only version 1 is supported and updates within
    this schema version should always be backwards compatible.


## 1.2 (2022-04-11)

**general**
- Add `mach init` command
- Skip non-MACH configuration files when processing all yaml files in a directory.<br>
  This allows you to run things like `mach apply` or `mach update` without having to specify the `-f main.yml` option if you only have one valid MACH configuration file in your directory. Fixes #150
- Ignore missing variables when running `mach sites` and `mach components`
- Add `--destroy` flag to the `plan` and `apply` commands
- Add `variables_file` option to the `mach_composer` configuration block to define a variable file
- Show commit author in `mach update` output
- Upgraded commercetools provider to [`0.30.0`](https://github.com/labd/terraform-provider-commercetools/blob/main/CHANGELOG.md#v0300-2021-08-04)
- Upgraded Amplience provider to [`0.3.7`](https://github.com/labd/terraform-provider-amplience/blob/main/CHANGELOG.md#v037-2022-03-14)
- Upgraded Azure provider to [`2.99.0`](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md#2990-march-11-2022)

**AWS**
- Upgraded Terraform AWS provider to [`3.74.1`](https://github.com/hashicorp/terraform-provider-aws/blob/main/CHANGELOG.md#3741-february-7-2022)
- Add support for default tags on provider level
  ```yaml
  aws:
    account_id: 123456789
    region: eu-central-1
    default_tags:
      environment: test
      owner: john
  ```


## 1.1 (2021-11-25)

**General**

- Variable support:
  - `${var.}` to be used with the `--var-file` command line option
  - `${component.}` to use component output values
  - `${env.}` to include environment variables in the configuration file

**AWS**
- Upgraded Terraform AWS provider to [`3.66.0`](https://github.com/hashicorp/terraform-provider-aws/blob/main/CHANGELOG.md#3660-november-18-2021)
- Add AWS specific endpoint options;
  - `enable_cdn` creates a CDN in front of an endpoint
  - `throttling_burst_limit` and `throttling_rate_limit` controls throttling on the API gateway

**Azure**

- Upgraded Terraform Azure provider to [`2.86.0`](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG.md#2860-november-19-2021)
- Add extra Frontdoor frontend_endpoint options:
  - session affinity
  - waf policy support
- Add extra Frontdoor routing options to components such as:
  - Custom routing paths
  - Health probe settings
  - Custom host address and ports
  - Caching options
- Include `frontend_endpoint` in ignore list when `suppress_changes` is used
- Add Frontdoor `ssl_key_vault` option to supply your own SSL certificate for your endpoints
- Add Azure specific endpoint options:
  - `internal_name` Overwrites the frontend endpoint name
  - `waf_policy_id` Defines the Web Application Firewall policy ID for the endpoint
  - `session_affinity_enabled`  Whether to allow session affinity
  - `session_affinity_ttl_seconds` The TTL to use in seconds for session affinity
- Add new `service_plans` option `per_site_scaling`
- Fix: set correct root-level DNS record (`@`) when endpoint URL is the same as the zone

**Commercetools**

- Upgraded Terraform commercetools provider to [`0.25.3`](https://github.com/labd/terraform-provider-commercetools/blob/master/CHANGELOG.md#v0293-2021-06-16)
- Add `tax_categories` to allow more complex tax setups. Does not work in conjunction with `taxes`

### Upgrade notes

**For Azure**

- Each component that has an `endpoint` defined needs to have an Terraform output defined for that endpoint. For example:<br>
  ```terraform
  output "azure_endpoint_main" {
    value = {
      address = azurerm_function_app.main.default_hostname
    }
  }
  ```
  Read more about the [configuration options](https://docs.machcomposer.io/reference/components/azure.html#defining-outputs).
- Remove endpoints restrictions: Azure components can now use multiple endpoints.
- Changes have been made in the Frontdoor configuration in the underlying Terraform Azure provider.<br>
  If you are using endpoints with a custom domain, you'll need to import the new [`azurerm_frontdoor_custom_https_configuration`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/frontdoor_custom_https_configuration#import) into your Terraform state.<br>
  More on how to work with the Terraform state in [our troubleshooting guide](https://docs.machcomposer.io/topics/development/troubleshooting.html).



## 1.0 (2021-05-10)

**New platforms**

- Add Amplience support
- Add Apollo Federation support
- Add Sentry DSN management options


**General**

- Add `mach_composer` configuration block to configure required MACH composer
  version
- [SOPS](https://github.com/mozilla/sops) support: SOPS-encrypted configuration files will get decrypted before being parsed further
- Add `--ignore-version` to disable the MACH composer version check
- Improved development workflow:
  - Improved git log parsing
  - Add `mach bootstrap` commands:
      - `mach bootstrap config` for creating a new MACH configuration
      - `mach bootstrap component` for creating a new MACH component
  - Add `--site` option to the `generate`, `plan` and `apply` commands
  - Add `--component` option to the `plan` and `apply` commands
  - Add `--reuse` flag to the `plan` and `apply` commands to suppress a
    `terraform init` call
  - Add support for relative paths to components
  - Add extra component definition settings `artifacts` to facilitate local
    deployments
- Improved dependencies between components and MACH-managed commercetools
  configurations
- Add option to override Terraform provider versions
- Add support for multiple API endpoints:
    - `base_url` replaced with `endpoints`
    - `has_public_api` replaced with `endpoints`
    - Supports a `default` endpoint that doesn't require custom domain
      settings
- Add support for including yaml files using the `${include(...)}`
  ```
  components: ${include(components.yml)}
  components: ${include(git::https://github.com/labd/mach-configs.git@9f42fe2//components.yml)}
  ```


**commercetools**

- Move `currencies`, `languages`, `countries`, `messages_enabled` to `project_settings` configuration block
- Add support for commercetools Store-specific variables and secrets on
  components included in new variable: `ct_stores`
- Add `managed` setting to commercetools store. Set to false it will indicate the store should not be managed by MACH composer
- Add support for commercetools shipping zones
- Make commercetools frontend API client scopes configurable with new
  `frontend` configuration block


**AWS**

- AWS: Set `auto-deploy` on API gateway stage
- AWS: Add new component variable `tags`


**Azure**

- Add configuration options for Azure service plans
- Upgraded Terraform to `0.14.5`
- Upgraded Terraform commercetools provider to `0.25.3`
- Upgraded Terraform AWS provider to `3.28.0`
- Upgraded Terraform Azure provider to `2.47.0`
- Azure: Remove `project_key` from `var.tags` and add `Environment` and `Site`
- Azure: Add `--with-sp-login` option to `mach plan` command
- Azure: Remove function app sync bash command: this is now the
  responsibility of the component


### Breaking changes

**Generic**

- **config**: Rename `general_config` to `global`
- **config**: `base_url` has been replaced by the `endpoints` settings:<br>
  ```yaml
  sites:
  - identifier: mach-site-eu
    base_url: https://api.eu-tst.mach-example.net
  ```
  becomes
  ```yaml
  sites:
  - identifier: mach-site-eu
    endpoints:
      main: https://api.eu-tst.mach-example.net
  ```
  When you name the endpoint that replaces `base_url` "main", it will have
  the least effect on your existing Terraform state.<br><br>
  When endpoints are defined on a component, the component needs to define
  endpoint Terraform variables
  ([AWS](https://docs.machcomposer.io/reference/components/aws.html#with-endpoints)
  and
  [Azure](https://docs.machcomposer.io/reference/components/azure.html#with-endpoints))
- **config**: commercetools `create_frontend_credentials` is replaced with new `frontend` block:
  ```terraform
  commercetools:
    frontend:
      create_credentials: false
  ```
  default is still `true`
- **config** Default scopes for commercetools frontend API client changed:
    - If you want to maintain previous scope set, define the following in the `frontend` block:
    ```terraform
    commercetools:
      frontend:
        permission_scopes: [manage_my_profile, manage_my_orders, view_states, manage_my_shopping_lists, view_products, manage_my_payments, create_anonymous_token, view_project_settings]
    ```
    - Old scope set didn't include store-specific [`manage_my_profile:project:store`](https://docs.commercetools.com/api/scopes#manage_my_profileprojectkeystorekey) scope. If you're using the old set as described above, MACH will need to re-create the store-specific API clients in order to add the extra scope. For migration options, see next point
    - In case the scope needs to be updated but (production) frontend
      implementations are already using the current API client credentials, a
      way to migrate is to;
      1. Remove the old API client resource with `terraform state rm commercetools_api_client.frontend_credentials`
      2. Repeat step for the store-specific API clients in your Terraform
         state
      3. Perform `mach apply` to create the new API clients with updated
         scope
      4. Your commercetools project will now contain API clients with the
         same name. Once the frontend implementation is migrated, the older one
         can safely be removed.
- **component**: Components with a `commercetools` integration require a new variable `ct_stores`:
  ```terraform
  variable "ct_stores" {
    type = map(object({
      key       = string
      variables = any
      secrets   = any
    }))
    default = {}
  }
  ```
- **component**: The folowing deprecated values in the `var.variables` are removed:
  ```terraform
  var.variables["CT_PROJECT_KEY"]
  var.variables["CT_API_URL"]
  var.variables["CT_AUTH_URL"]
  ```
  See [0.5.0 release notes](#050-2020-11-09)
- **component**: The `var.environment_variables` won't be set by MACH anymore. Use `var.variables` for this


**AWS**

- **config**: The AWS `route53_zone_name` setting has been removed in favour of multiple endpoint support
- **config**: The `deploy_role` setting has been renamed to `deploy_role_name`
- **component**: Introduced new variable `tags`:
  ```terraform
  variable "tags" {
    type        = map(string)
    description = "Tags to be used on resources."
  }
  ```
- **component**: Add `aws_endpoint_*` variable when the `endpoints` configuration option is used. [More information](https://docs.machcomposer.io/reference/components/aws.html#with-endpoints) on defining and using endpoints in AWS.

**Azure**

- **config**: The `front_door` configuration block has been renamed to `frontdoor`
- **config**: The Azure frontdoor settings `dns_zone` and `ssl_key_*` settings have been removed;<br>
  Certificates are now managed by Frontdoor and dns_zone is auto-detected.
- **config**: The Azure frontdoor settings `resource_group` has been renamed to `dns_resource_group`
- **config**: Moved component `short_name` to new `azure` configuration block
- **state**: The Terraform `azurerm_dns_cname_record` resources have been renamed; they now take the name of the associated endpoint key. For the smoothest transition, rename them in your Terraform state:<br>
  ```bash
  terraform state mv azurerm_dns_cname_record.<project-key> azurerm_dns_cname_record.<endpoint-key>
  ```
- **component**: Prefixed **all** Azure-specific variables with `azure_`
- **component**: The `FRONTDOOR_ID` value is removed from the `var.variables`
  of a component. Replaced with `var.azure_endpoint_*`. [More
  information](https://docs.machcomposer.io/reference/components/azure.html#with-endpoints)
  on defining and using endpoints in Azure.
- **component**: `app_service_plan_id` has been replaced with
  `azure_app_service_plan` containing both an `id` and `name` so the
  [azurerm_app_service_plan](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/app_service_plan)
  data source can be used in a component.<br>

  It will *only be set* when `service_plan` is configured in the component
  [definition](https://docs.machcomposer.io/reference/syntax/components.html#azure)
  or [site configuration](https://docs.machcomposer.io/reference/syntax/sites.html#azure_1)
  ```terraform
  variable "azure_app_service_plan" {
    type = object({
      id                  = string
      name                = string
      resource_group_name = string
    })
  }
  ```
- **component**: Replaced `resource_group_name` and `resource_group_location` with `azure_resource_group`:
  ```terraform
  variable "azure_resource_group" {
    type = object({
      name     = string
      location = string
    })
  }
  ```

## 0.5.1 (2020-11-10)
- Removed `aws` block in general_config
- Add `branch` option to component definitions to be able to perform a `mach
  update` and stay within a certain branch (during development)


## 0.5.0 (2020-11-09)
- Add new CLI options:
    - `mach components` to list all components
    - `mach sites` to list all sites
- Improved `update` command:
    - Supports updating (or checking for updates) on all components based on their git history
    - This can now also be used to manually update a single component; `mach update my-component v1.0.4`
    - Add `--commit` argument to automatically create a git commit message
- Add new AWS configuration option `route53_zone_name`
- Remove unused `api_gateway` attribute on AWS config
- Remove restriction from `environment` value; can now be any. Fixes #9

### Breaking changes

- Require `ct_api_url` and `ct_auth_url` for components with `commercetools` integration

### Deprecations

In a component, the use of the following variables have been deprecated;

```
var.variables["CT_PROJECT_KEY"]
var.variables["CT_API_URL"]
var.variables["CT_AUTH_URL"]
```

Instead you should use:

```
var.ct_project_key
var.ct_api_url
var.ct_auth_url
```

## 0.4.3 (2020-11-04)
- Make AWS role definitions optional so MACH can run without an 'assume role' context


## 0.4.2 (2020-11-02)
- Add 'encrypt' option to AWS state backend
- Correctly depend component modules to the commercetools project settings resource
- Extend Azure regions mapping


## 0.4.1 (2020-10-27)
- Fixed TypeError when using `resource_group` on site Azure configuration


## 0.4.0 (2020-10-27)
- Add Contentful support

### Breaking changes

- `is_software_component` has been replaced by the `integrations` settings

```
components:
  - name: my-product-types
    source: git::ssh://git@github.com/example/product-types/my-product-types.git
    version: 81cd828
    is_software_component: false
```

becomes

```
components:
  - name: my-product-types
    source: git::ssh://git@github.com/example/product-types/my-product-types.git
    version: 81cd828
    integrations: ["commercetools"]
```

or `integrations: []` if no integrations are needed at all.


## 0.3.0 (2020-10-21)
- Add option to specify custom resource group per site

### Breaking changes

- All `resource_group_name` attributes is renamed to `resource_group`
- The `storage_account_name` attribute is renamed to `storage_account`


## 0.2.2 (2020-10-15)
- Fixed Azure config merge: not all generic settings where merged with
  site-specific ones
- Only validate short_name length check for Azure implementations
- Setup Frontdoor per 'public api' component regardless of global Frontdoor
  settings


## 0.2.1 (2020-10-06)
- Fixed rendering of STORE environment variables in components
- Updated Terraform version to 0.13.4
- Fix `--auto-approve` option on `mach apply` command


## 0.2.0 (2020-10-06)
- Add AWS support
- Add new required attribute `cloud` in general config


## 0.1.0 (2020-10-01)
- Initial release

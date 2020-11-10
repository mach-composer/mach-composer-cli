# Changelog

## 0.5.1 (2020-11-10)
- Removed `aws` block in general_config
- Add `branch` option to component definitions to be able to perform a `mach update` and stay within a certain branch (during development)
  

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
- Fixed Azure config merge: not all generic settings where merged with site-specific ones
- Only validate short_name length check for Azure implementations
- Setup Frontdoor per 'public api' component regardless of global Frontdoor settings


## 0.2.1 (2020-10-06)
- Fixed rendering of STORE environment variables in components
- Updated Terraform version to 0.13.4
- Fix `--auto-approve` option on `mach apply` command


0.2.0 (2020-10-06)
=================
- Add AWS support
- Add new required attribute `cloud` in general config
  

## 0.1.0 (2020-10-01)
- Initial release
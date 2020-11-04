
0.4.3 (2020-11-04)
==================
- Make AWS role definitions optional so MACH can run without an 'assume role' context


0.4.2 (2020-11-02)
==================
- Add 'encrypt' option to AWS state backend
- Correctly depend component modules to the commercetools project settings resource
- Extend Azure regions mapping
  

0.4.1 (2020-10-27)
==================
- Fixed TypeError when using `resource_group` on site Azure configuration


0.4.0 (2020-10-27)
==================
- Add Contentful support

Breaking changes
----------------
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


0.3.0 (2020-10-21)
==================
- Add option to specify custom resource group per site
  
Breaking changes
----------------
- All `resource_group_name` attributes is renamed to `resource_group`
- The `storage_account_name` attribute is renamed to `storage_account`


0.2.2 (2020-10-15)
==================
- Fixed Azure config merge: not all generic settings where merged with site-specific ones
- Only validate short_name length check for Azure implementations
- Setup Frontdoor per 'public api' component regardless of global Frontdoor settings


0.2.1 (2020-10-06)
==================
- Fixed rendering of STORE environment variables in components
- Updated Terraform version to 0.13.4
- Fix `--auto-approve` option on `mach apply` command


0.2.0 (2020-10-06)
=================
- Add AWS support
- Add new required attribute `cloud` in general config
  

0.1.0 (2020-10-01)
==================
- Initial release
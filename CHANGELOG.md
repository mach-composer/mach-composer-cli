0.2.2 (unreleased)
==================
- Fixed Azure config merge: not all generic settings where merged with site-specific ones
- Only validate short_name length check for Azure implementations
- Setup Frontdoor per 'public api' component regardless of global Frontdoor settings


0.2.1 (06-10-2020)
==================
- Fixed rendering of STORE environment variables in components
- Updated Terraform version to 0.13.4
- Fix `--auto-approve` option on `mach apply` command


0.2.0 (06-10-2020)
=================
- Add AWS support
- Add new required attribute `cloud` in general config
  

0.1.0 (01-10-2020)
==================
- Initial release
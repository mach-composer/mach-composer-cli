# Step 2. Setup commercetools

!!! note "Optional"
    This step is only necessary if you are going to use the commercetools integration in your MACH stack

Create a API client *'to rule them all'*.

Required scopes:

- `manage_api_clients`
- `manage_project`
- `view_api_clients`

!!! note ""
    This client is used the MACH composer to create other necessary commercetools clients for each individual component.

Use the credentials for this client to configure each site's [commercetools settings](../syntax/sites.md#commercetools).

## Next steps

Setup your [Azure](./azure.md) or [AWS](./aws.md) environment.

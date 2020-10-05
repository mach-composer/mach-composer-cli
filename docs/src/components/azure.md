# Azure components

# Azure Functions

For Azure functions, this means uploading a packaged ZIP file to a storage account from where fuction apps can download the function app code.

!!! tip ""
    An example build and deploy script is provided in the [component cookiecutter](https://git.labdigital.nl/mach/component-cookiecutter)

## HTTP routing

MACH will provide the correct HTTP routing for you.  
To do so, the following has to be configured:

- [Frontdoor](../syntax.md#front_door) settings in the Azure configuration
- The component has to be marked as [`has_public_api`](../syntax.md#components)

More information in the [deployment section](../deployment/azure.md#http-routing).

# Azure dashboard configuration

!!! Todo
    Future implementation
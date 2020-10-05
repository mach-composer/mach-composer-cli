# Components

A MACH component is in it's bare minimum a [*Terraform module*](https://www.terraform.io/docs/configuration/modules.html).

Other then the [Terraform configuration](#terraform-component), a component might include a:

- [serverless function](#serverless-function)
- [Azure dashboard configuration](#azure-dashboard-configuration)

#### Deployment process

The deployment of a full-fledged component typically flows through the following steps:

1. Serverless function is built, packaged up and uploaded to a shared resources all environments and sites can access.  
   **Note** at this point, no actual deployment is made; the function doesn't run yet.
2. At the moment the MACH composer deploys a site's Terraform configuration, it uses the component's Terraform configuration to make the necessary modifications to the resources.  
   For example: create the function app instance, necessary routing, etc.
3. MACH composer will use the packages function (from step 1) to deploy the function itself


## Terraform module

The component must be able to instruct MACH wat resources to create in the cloud infrastructure.

This is done by providing the necessary Terraform module files in the `terraform/` directory.


## Serverless function

The component might contain code for a serverless function to run on the cloud provider.

What kind of language/runtime is used for that is irrelevant to MACH. Two things the component needs to contain:

- Build/deploy script to build, package and upload the serverless function to a repository
- A Terraform configuration for the serverless function

### Azure Functions

For Azure functions, this means uploading a packaged ZIP file to a storage account from where fuction apps can download the function app code.

!!! tip ""
    An example build and deploy script is provided in the [component cookiecutter](https://git.labdigital.nl/mach/component-cookiecutter)

#### HTTP routing

MACH will provide the correct HTTP routing for you.  
To do so, the following has to be configured:

- [Frontdoor](./syntax.md#frontdoor) settings in the Azure configuration
- The component has to be marked as [`has_public_api`](./syntax.md#components)

More information in the [deployment section](./deployment/azure.md#http-routing).

### AWS Lambda

AWS Lambda functions need to be uploaded to a S3 bucket. From there AWS Lambda will run the functions for you once instructed by the Terraform deployment.

#### HTTP routing

MACH will provide the correct HTTP routing for you.  
To do so, the following has to be configured:

- [Frontdoor](./syntax.md#frontdoor) settings in the Azure configuration
- The component has to be marked as [`has_public_api`](./syntax.md#components)

More information in the [deployment section](./deployment/azure.md#http-routing).


## Azure dashboard configuration

!!! Todo
    Future implementation
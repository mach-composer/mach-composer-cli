# Step 3. Prepare your Azure environment

In this step we will set up the following in Azure:

1. Storage account for Terraform states
2. Storage account for function app code packages
3. Registered providers needed for various services

This section will demonstrate how to setup Azure using Terraform.

## 1. Prepare

Login to your Azure subscription through the CLI.

```bash
$ az login
```
Follow the prompts. On success, the CLI will respond with a JSON object of the subscriptions available to you.<br>
Make sure the subscription you want to work in is set to default. If it is not, you can run

```bash
$ az account set --subscription <name or id>
```

## 2. Register providers

Make sure the following providers are registered on the subscription.

```bash
$ az provider register --namespace Microsoft.Web
$ az provider register --namespace Microsoft.KeyVault
$ az provider register --namespace Microsoft.Storage
$ az provider register --namespace Microsoft.Insights
```

!!! info "More info"
    [https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/error-register-resource-provider](https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/error-register-resource-provider#solution-2---azure-cli)

## 3. Terraform configuration
Setup your Terraform configuration.

In main.tf:

```terraform
terraform {
  required_version = ">= 0.14.0"
}
```

In modules.tf:

```terraform
module "shared_infra" {
  source         = "git@github.com:labd/terraform-azure-mach-shared.git"
  name_prefix    = "mach-shared-we"  # replace 'mach' with your desired prefix
  region         = "westeurope"
  dns_zone_name  = "example.com"
  certificate_access_object_ids = [
     ... # User, group or principle IDs that should have access to the resources
  ]
}
```

[More info](https://github.com/labd/terraform-azure-mach-shared) about the shared infra module.

## 4. Apply Terraform
1. Run the following commands:
```bash
$ terraform init
$ terraform apply
```
2. For a new Terraform setup, initially it will store the Terraform state locally and should be named `terraform.tfstate`.<br>
   We'll move this state to the Storage Account that has been created by the shared infra module.<br>
   To do this, add a backend setting to project like below
```terraform
terraform {
 required_version = ">= 0.14.0"
 backend "azurerm" {
 }
}
```
3. Now run:
```bash
$ terraform init -reconfigure 
```
Terraform will detect that you're trying to move your state into Azure and ask; "*Do you want to copy existing state to the new backend?*".<br>
Enter **"yes"**.<br>
Now the state is stored in the Storage Account and the DynamoDB table will be used to lock the state to prevent concurrent modifications.
4. Check if `terraform.tfstate` is empty and remove it.<br>
   Repeat the above three steps for all other environments

## Example

See the [examples directory](https://github.com/labd/mach-composer/tree/master/examples/azure/infra/) for an example of a Terraform setup


## Manual setup

See instructions on how to [setup Azure manually](./manual.md).
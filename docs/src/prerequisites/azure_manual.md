# Manual Azure setup

## Create Terraform storage

### Storage account
Create a storage account which will be used as Terraform state backend.

![Create storage account](../_img/azure/terraform_storage_account.png)


!!! tip
    A good convention is to place the Terraform state backend storage account in a 'shared' resource group which can be used for various shared resources accross all your environments and sites.<br>
    For example:<br>
    **Resource group**: `my-shared-we-rg`<br>
    **Storage account** `mysharedwesaterra`<br>
    Where 'my' is replaced by a prefix of your choosing.

### Create container
Create a container in the storage account. Name it for example `tfstate`.

## Create function app storage
All packaged function app code should be stored on the shared environment from where all other envirnoment can access those assets.

### Storage account

Create a new `BlockBlobStorage` with a Premium account tier for improved performace.

!!! tip
    Again, like the Terraform state, place this in a 'shared' resource group
    For example:<br>
    **Resource group**: `my-shared-we-rg`<br>
    **Storage account** `mysharedwesacomponents`<br>
    Where 'my' is replaced by a prefix of your choosing.

### Create container

Create a blob container called `code`. Make this private.

## Register providers

Make sure the following providers are registered on the subscription:

- `Microsoft.Web`
- `Microsoft.KeyVault`
- `Microsoft.Storage`
- `Microsoft.Insights`

More info:<br>
[https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/error-register-resource-provider](https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/error-register-resource-provider)

# Setting up Azure

## Create Terraform storage

### Storage account
Create a storage account which will be used as Terraform state backend.

![Create storage account](../_img/azure/terraform_storage_account.png)


!!! tip
    A good convention is to place the Terraform state backend storage account in a 'shared' resource group which can be used for various shared resources accross all your environments and sites.  
    For example:  
    **Resource group**: `my-shared-we-rg`  
    **Storage account** `mysharedwesaterra`  
    Where 'my' is replaced by a prefix of your choosing.

### Create container
Create a container in the storage account. Name it for example `tfstate`.

## Create function app storage

### Storage account

### Create container
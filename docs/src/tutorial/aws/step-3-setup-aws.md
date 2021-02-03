# Step 3. Prepare your AWS environment

As described in the [tenenacy model](../guidance/tenancy.md#aws-tenancy), we advice to setup your MACH environment by creating **one service AWS account** containing shared resources and create an **AWS account per stack**.

This way, all resources are strictly seperated from eachother.

These account contain at least the following resources:

**Service account**

1. Terraform state backend
2. Component registry
3. Route53 zone to route all other accounts from

**Site-specific account**

1. Terraform state backend
2. The Route53 hosted zones needed for the endpoints
3. `deploy` IAM role for MACH to manage your resources


## Service account setup

!!! todo
    Describe setup using Terraform


!!! tip
    Keep your created component repository information at hand for later: you'll need it when creating a new component.

    This means; the resource group, storage account and container name in Azure or the AWS S3 bucket name.

## Site-specific account setup

!!! todo
    Describe setup using Terraform

## Example

See the [examples directory](https://github.com/labd/mach-composer/tree/master/examples/aws/infra/) for an example of a Terraform setup

## Manual setup

See instructions on how to [setup AWS manually](./aws_manual.md).
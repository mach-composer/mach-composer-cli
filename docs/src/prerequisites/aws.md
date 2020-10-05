# Setting up AWS

## Create S3 state backend
Create a S3 bucket which will be used as Terraform state backend.

!!! info "Setting up AWS"
    For more information on how to setup a S3 state backend including the correct IAM roles, see the [Terraform documentation](https://www.terraform.io/docs/backends/types/s3.html#s3-bucket-permissions)


## Create lambda repository

!!! Todo
    Describe how to setup buckets and permissions to deploy lambda functions and be able to access them from different accounts.


## Setup AWS account per site

It is recommended to use a dedicated AWS account **per site**.  
This way, all resources are strictly seperated from eachother.

So the following steps need to be done per site:

1. Create an AWS account
2. Create a 'deploy' IAM role for MACH to manage your resources

### IAM deploy role

Make sure the IAM role has sufficient permissions to manage your resources.

!!! TODO 
    Include example policy

# Step 4. Add AWS site account

Now that we've created our [service account](./step-3-setup-aws-services.md), we can create an AWS account for a specific site (or; MACH stack).

## Setup site-specific

For this account we will create:

1. Terraform state backend
2. `deploy` IAM role for MACH to manage your resources

### 1. Create AWS account

- In your AWS console, go to **My Organization** and choose **Add accounts**
- For your new account choose a name like `yourproject-tst`.
  
### 2. Setup your Terraform configuration

Within your `mach-account` directory [^1] create the following files:

#### `variables.tf`

```terraform
variable "aws_account_id" {}

variable "site_name" {}

variable "region" {
  default = "eu-central-1"
}
```

#### `main.tf`

```terraform
provider "aws" {
  region = var.region

  assume_role {
    role_arn = "arn:aws:iam::${var.aws_account_id}:role/sudo"
  }
}
```

#### `modules.tf`

```terraform
module "tfstate-backend" {
  source  = "cloudposse/tfstate-backend/aws"
  version = "0.33.0"
}

module "mach_account" {
  source = "git::https://github.com/labd/terraform-aws-mach-account.git"
  code_repository_name = "your-project-lambdas"  # Replace with the actual name given to the S3 bucket
}
```

[^1]: Refer to the [previous step](./step-3-setup-aws-services.md#2-setup-your-terraform-configuration) to see how we organize the two different AWS accounts

### 3. Create the first environment configuration

Create a directory called `mach-account/envs/` and create a new file `tst.tfvars`:

```terraform
aws_account_id = "<your-account-id>"
site_name      = "tst"
```

### 4. Terraform roll-out

1. Within your `mach-account` directory, run the following commands:
```bash
$ terraform init -var-file=envs/tst.tfvars 
$ terraform apply -var-file=envs/tst.tfvars 
```
2. Terraform has now createsd a `backend.tf` file which instructs Terraform to store the state on a S3 bucket.<br>
In order to move the current (local) state file to the bucket, perform this one-time command:
```bash
$ terraform init -force-copy -var-file=envs/tst.tfvars 
```
Now the state is stored in the S3 bucket, and the DynamoDB table will be used to lock the state to prevent concurrent modification.

### 5. Grant access to service account

The last step is to allow the new AWS account to access resources from the service account.

In `service/modules.tf` add the `allow_code_repo_read_access` variable to the `shared-config` module:

```terraform
module "shared-config" {
  source = "git::https://github.com/labd/terraform-aws-mach-shared.git"
  allow_code_repo_read_access = [
      "arn:aws:iam::<test-aws-account-id>:user/mach-user", # MACH Test env
  ]
```

!!! note ""
    This will make sure that this AWS account is able to read the contents of the component repository bucket.

Run `terraform apply` to apply these changes.

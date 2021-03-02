# Step 4. Add AWS site account

Now that we've created our [service account](./step-3-setup-aws-services.md), we can create an AWS account for a specific site (or; MACH stack).

## Setup site-specific

For this account we will create:

1. Terraform state backend
2. `deploy` IAM role for MACH to manage your resources

### 1. Create AWS account

- In your AWS console, go to **My Organization** and choose **Add accounts**
- For your new account choose a name like `your-project-tst`
- As **IAM role name** enter `admin`
  
### 2. Setup your Terraform configuration

Within your `mach-account` directory [^1] create the following files:

#### `variables.tf`

```terraform
variable "aws_account_id" {
    type = string
}

variable "name" {
    type = string
}

variable "region" {
  default = "eu-central-1"
}
```

#### `main.tf`

```terraform
locals {
  role_arn            = "arn:aws:iam::${var.aws_account_id}:role/admin"
  tfstate_bucket_name = "${var.name}-tfstate"
}

terraform {
  # We will uncomment this later
  # backend "s3" {}
}

provider "aws" {
  region = var.region

  assume_role {
    role_arn = local.role_arn
  }
}
```

#### `modules.tf`

```terraform
module "tfstate-backend" {
  source  = "cloudposse/tfstate-backend/aws"
  version = "0.33.0"

  s3_bucket_name = local.tfstate_bucket_name
  role_arn       = local.role_arn
}

module "mach_account" {
  source               = "git::https://github.com/labd/terraform-aws-mach-account.git"
  aws_account_alias    = var.name
  code_repository_name = "your-project-lambdas"  # Replace with the actual name given to the S3 bucket
  deploy_principle_identifiers = [
    "arn:aws:iam::000000000000:user/admin" # Specify your root account here
  ]
}
```

!!! info "`deploy_principle_identifiers`"
    We specify our root account here so it makes it easier for this tutorial to setup credentials to be able to deploy using MACH.

#### `policies.tf`

The `terraform-aws-mach-account` module will create the necessary IAM policies that allows the mach deploy user to deploy the necessary resources.

The Terraform state backend must also be used by mach, so we need to create the necessary policies that allows the mach user to read/write to that state backend:

```terraform
data "aws_iam_policy_document" "terraform_state" {
    
  statement {
    actions = [
      "s3:ListBucket"
    ]
    resources = [
      "${module.tfstate-backend.s3_bucket_arn}"
    ]
  }
  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject"

    ]
    resources = [
      "${module.tfstate-backend.s3_bucket_arn}/mach/*"
    ]
  }

  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:DeleteItem"

    ]
    resources = [
      "arn:aws:dynamodb:${var.region}:${var.aws_account_id}:table/${local.tfstate_bucket_name}-lock"
    ]
  }
}


resource "aws_iam_policy" "terraform_state" {
  name        = "terraform-state-policy"
  path        = "/"
  description = "Policy to access terraform state"

  policy = data.aws_iam_policy_document.terraform_state.json
}

resource "aws_iam_role_policy" "apigateway" {
  name   = "terraform-state-policy"
  role   = module.mach_account.mach_role_id
  policy = data.aws_iam_policy_document.terraform_state.json
}

resource "aws_iam_user_policy_attachment" "mach_user_terraform_state" {
  user       = module.mach_account.mach_user_name
  policy_arn = aws_iam_policy.terraform_state.arn
}
```
### 3. Create the first environment configuration

Create a directory called `mach-account/envs/` and create a new file `tst.tfvars`:

```terraform
aws_account_id = "<your-account-id>"
name           = "your-project-tst"
```

### 4. Terraform roll-out

Within your `mach-account` directory, run the following commands:
```bash
$ terraform init -var-file=envs/tst.tfvars 
$ terraform apply -var-file=envs/tst.tfvars 
```

### 5. Configure state backend

Terraform has now created the Terraform state backend.

We are going to store that information in a site-specific backend configuration file. This way, several backends for multiple sites can live side-by-side in the same infra repo.

1. Create a new directory `mach-account/backend-configs` and create a new file `tst.conf`:
```
region         = "eu-central-1"
bucket         = "your-project-tst-tfstate"
key            = "terraform.tfstate"
dynamodb_table = "lock"
role_arn       = "arn:aws:iam::<account-id>:role/admin"
encrypt        = "true"
```
2. Uncomment the `# backend "s3" {}` line in `main.tf`
3. Perform the following command:
```bash
$ terraform init -force-copy -var-file=envs/tst.tfvars  -backend-config=backend-configs/tst.conf 
```
Now the state is stored in the S3 bucket, and the DynamoDB table will be used to lock the state to prevent concurrent modification.

### 6. Grant access to service account

The last step is to allow the new AWS account to access resources from the service account.

In `service/modules.tf` add the `allow_code_repo_read_access` variable to the `shared-config` module:

```terraform
module "shared-config" {
  source = "git::https://github.com/labd/terraform-aws-mach-shared.git"
  allow_code_repo_read_access = [
      "arn:aws:iam::<test-aws-account-id>:user/mach", # test env
      "arn:aws:iam::<test-aws-account-id>:role/mach", # test env
  ]
```

!!! note ""
    This will make sure that this AWS account is able to read the contents of the component repository bucket.

Run `terraform apply` to apply these changes.

!!! tip "Next: step 5"
    Next we'll create our first [MACH component](./step-5-create-component.md).


[^1]: Refer to the [previous step](./step-3-setup-aws-services.md#2-setup-your-terraform-configuration) to see how we organize the two different AWS accounts

module "tfstate-backend" {
  source  = "cloudposse/tfstate-backend/aws"
  version = "0.33.0"

  s3_bucket_name = local.tfstate_bucket_name
  role_arn       = local.role_arn
}

module "mach_account" {
  source               = "git::https://github.com/labd/terraform-aws-mach-account.git"
  code_repository_name = "your-project-lambdas"  # Replace with the actual name given to the S3 bucket
  aws_account_alias    = var.name
  deploy_principle_identifiers = [
    # Specify accounts that should be able to assume the deploy role
  ]
}

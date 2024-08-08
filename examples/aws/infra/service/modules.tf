module "tfstate-backend" {
  source                             = "cloudposse/tfstate-backend/aws"
  version                            = "0.33.0"
  s3_bucket_name                     = "${var.name}-srv-tfstate"
  role_arn                           = local.role_arn
  terraform_backend_config_file_path = "."
  terraform_backend_config_file_name = "backend.tf"
}

module "shared-config" {
  source            = "git::https://github.com/labd/terraform-aws-mach-shared.git"
  code_repo_name    = "${var.name}-lambdas"
  aws_account_alias = "${var.name}-srv"
  allow_code_repo_read_access = [
    # Include sub-accounts here for each MACH environment
  ]
  allow_assume_deploy_role = [
    # Specify accounts that should be able to assume the deploy role
  ]
}

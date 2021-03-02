locals {
  role_arn = "arn:aws:iam::${var.aws_account_id}:role/admin"
}

provider "aws" {
  region = var.region

  assume_role {
    role_arn = local.role_arn
  }
}

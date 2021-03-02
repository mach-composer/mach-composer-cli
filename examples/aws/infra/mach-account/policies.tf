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
  description = "Policy to access terraform satte"

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

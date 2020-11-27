# AWS components

Unless a component is flagged as `is_software_component: false`, a component within a AWS-based MACH configuration is considered an *AWS component*.

To be able to create the resources needed, a couple of extra [Terraform variables](#terraform-variables) are set by MACH.

In addition to this, the component itself is responsible for [packaging and deploying](#packaging-and-deploying) the correct assets in case of a Lambda function.

## Terraform variables

In addition to the [base variables](./index.md#required-variables) AWS components don't require additional variables, unless an `endpoint` is expected (and set in the configuration).

### With `endpoint`

In order to support the [`endpoint`](../deployment/config/aws.md#http-routing) attribute on the component, the component needs to have the following variables defined:

- `api_gateway` - API Gateway ID to publish in (only if component has an `endpoint` defined)
- `api_gateway_execution_arn` API Gateway API Execution ARN (only if component is has an `endpoint` defined)


```terraform
variable "api_gateway" {}
variable "api_gateway_execution_arn" {}
```

## Packaging and deploying

AWS Lambda functions need to be uploaded to a S3 bucket. From there AWS Lambda will run the functions for you once instructed by the Terraform deployment.

[Read more](../deployment/components.md#on-aws) about AWS component deployments.

### Configure runtime
When defining your AWS Lambda function resource, you can reference back to the asset that is deployed:

```terraform
resource "aws_lambda_function" "example" {
  s3_bucket = "your-lambda-repo"
  s3_key    = "yourcomponent-${var.component_version}.zip"
  ...
}
```
## HTTP routing

MACH will provide the correct HTTP routing for you.<br>
To do so, the following has to be configured:

- [endpoints](../syntax.md#sites) settings in the Site configuration
- The component needs to have a [`endpoint`](../syntax.md#components) defined

More information in the [deployment section](../deployment/config/aws.md#http-routing).

## Lambda function

We recommend using the [AWS Lambda Terraform module](https://registry.terraform.io/modules/terraform-aws-modules/lambda/aws/latest) for managing a Lambda function.

```terraform
module "lambda_function" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${var.site}-${var.short_name}"
  handler       = "src/index.handler"
  runtime       = "nodejs12.x"
  memory_size   = 512
  timeout       = 10

  environment_variables = local.environment_variables
  create_package = false
  s3_existing_package = {
    bucket = local.lambda_s3_repository
    key    = local.lambda_s3_key
  }
}
```

See also [notes on using the serverless framework](../deployment/config/components.md#serverless-framework)
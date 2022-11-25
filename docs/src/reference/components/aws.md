# AWS components

All components within a AWS-based MACH Composer configuration are automatically
considered to have a 'aws' integration by default. Only if 'aws' is explicitely
omitted from the `integrations` definition, it won't require any AWS-specific
variables.

To be able to create the resources needed, a couple of extra
[Terraform variables](#terraform-variables) are set by MACH Composer.

In addition to this, the component itself is responsible for
[packaging and deploying](#packaging-and-deploying) the correct assets in case
of a Lambda function.

## Terraform variables

In addition to the [base variables](./structure.md#required-variables) AWS
components don't require additional variables, unless an `endpoint` is
expected (and set in the configuration).

### With `endpoints`

In order to support the [`endpoints`](../../topics/deployment/config/aws.md#http-routing)
attribute on the component, the component needs to define what endpoints it
expects.

For example, if the component requires two endpoints (`main` and `webhooks`) to
be set, the following variables needs to be defined:

```terraform
variable "aws_endpoint_main" {
  type = object({
    url                       = string
    api_gateway_id            = string
    api_gateway_execution_arn = string
  })
}

variable "aws_endpoint_webhooks" {
  type = object({
    url                       = string
    api_gateway_id            = string
    api_gateway_execution_arn = string
  })
}
```

## Packaging and deploying

AWS Lambda functions need to be uploaded to a S3 bucket. From there AWS Lambda
will run the functions for you once instructed by the Terraform deployment.

[Read more](../../topics/deployment/components.md#on-aws) about AWS component
deployments.

### Configure runtime
When defining your AWS Lambda function resource, you can reference back to the
asset that is deployed:

```terraform
resource "aws_lambda_function" "example" {
  s3_bucket = "your-lambda-repo"
  s3_key    = "yourcomponent-${var.component_version}.zip"
  ...
}
```
## HTTP routing

MACH Composer will provide the correct HTTP routing for you.<br>
To do so, the following has to be configured:

- [endpoints](../syntax/sites.md) settings in the Site configuration
- The component needs to have [`endpoints`](../syntax/components.md) defined

!!! tip "Default endpoint"
    If you assign `default` to one of your components endpoints, no additional
    Route53 settings are needed.

    MACH Composer will create an API Gateway for you without any custom domain.

More information in the [deployment section](../../topics/deployment/config/aws.md#http-routing).

## Lambda function

We recommend using the [AWS Lambda Terraform module](https://registry.terraform.io/modules/terraform-aws-modules/lambda/aws/latest)
for managing a Lambda function.

```terraform
module "lambda_function" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${var.site}-${var.azure_short_name}"
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

See also [notes on using the serverless framework](../../topics/deployment/config/components.md#serverless-framework)

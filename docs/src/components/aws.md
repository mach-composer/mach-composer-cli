# AWS components

Unless a component is flagged as `is_software_component: false`, a component within a AWS-based MACH configuration is considered an *AWS component*.

To be able to create the resources needed, a couple of extra [Terraform variables](#terraform-variables) are set by MACH.

In addition to this, the component itself is responsible for [packaging and deploying](#packaging-and-deploying) the correct assets in case of a Lambda function.

## Terraform variables

In addition to the [base variables](./index.md#required-variables), an AWS component expects the following:

- `api_gateway` - API Gateway ID (only if component is marked as `has_public_api`)
- `code_repository` - Name of the Lambda code repository (only if component is marked as `is_software_component`)

```terraform
variable "api_gateway" {
    default = ""
}
variable "code_repository" {
    default = ""
}
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

MACH will provide the correct HTTP routing for you.  
To do so, the following has to be configured:

- [base_url](../syntax.md#sites) settings in the Site configuration
- The component has to be marked as [`has_public_api`](../syntax.md#components)

More information in the [deployment section](../deployment/aws.md).



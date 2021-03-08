# Component deployment - second stage

The '*second stage*' deployment of a component is the process of creating and managing the resources that are defined in the [Terraform configuration](../../../reference/components/structure.md#terraform-module) of that component.

In case a component contains one or more serverless functions, the necessary Lambda Function or Azure Function App needs to be created and configured with the correct artifact.<br>
This artifact will need to be built, packaged and uploaded in the ['*first stage*' deployment](../components.md).

## Creating resources

On top of the resources that are [created by MACH composer](./index.md) by providing the correct Terraform configuration in the component any resource can be created.

Some of the components are referenced through the Terraform variables and can be used within the component to link them with the new component resources. For example when creating an route to an API endpoint in AWS:

```terraform
resource "aws_apigatewayv2_integration" "gateway" {
    api_id           = var.aws_endpoints.main.api_gateway_id
    integration_type = "AWS_PROXY"
    integration_uri  = local.lambda_function_arn
}

resource "aws_apigatewayv2_route" "app_route" {
    api_id    = var.aws_endpoints.main.api_gateway_id
    route_key = "ANY /graphql/{proxy+}"
    target    = "integrations/${aws_apigatewayv2_integration.gateway.id}"
}
```
Where the `var.api_gateway` references the [Gateway created by MACH](./aws.md#http-routing)

## Additional notes

### Serverless framework

As described in the [AWS](../../../reference/components/aws.md#lambda-function) and [Azure](../../../reference/components/azure.md#function-app) instructions of the component structure, we recommend taking full control in managing the necessary resources for a serverless function.

Although we recommend [using the Serverless framework](../../components/index.md#using-serverless-framework) for local development, we do not recommend using it for the actual deployment.

!!! info "Terraform Serverless provider"
    There is a [Terraform serverless provider](https://registry.terraform.io/providers/labd/serverless/latest) which could be used in the Terraform configuration of a component.

    However, the Serverless framework uses CloudFormation to manage the stack whereas most ideal scenario is to have everything fully managed in Terraform. Therefor we don't recommend using that at the moment.

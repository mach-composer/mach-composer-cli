# Component deployment - second stage

The '*second stage*' deployment of a component is the process of creating and
managing the resources that are defined in the
[Terraform configuration](../components/index.md) of
that component.

In case a component contains one or more serverless functions, the necessary
Lambda Function or Azure Function App needs to be created and configured with
the correct artifact.<br>
This artifact will need to be built, packaged and uploaded in the
['*first stage*' deployment](./first-stage.md).

## Creating resources

On top of the resources that are [created by MACH composer](./index.md) by
providing the correct Terraform configuration in the component any resource can
be created.

Some of the components are referenced through the Terraform variables and can be
used within the component to link them with the new component resources. For
example when creating a route to an API endpoint in AWS:

```terraform
resource "aws_apigatewayv2_integration" "gateway" {
    api_id           = var.aws_endpoint_main.api_gateway_id
    integration_type = "AWS_PROXY"
    integration_uri  = local.lambda_function_arn
}

resource "aws_apigatewayv2_route" "app_route" {
    api_id    = var.aws_endpoint_main.api_gateway_id
    route_key = "ANY /graphql/{proxy+}"
    target    = "integrations/${aws_apigatewayv2_integration.gateway.id}"
}
```

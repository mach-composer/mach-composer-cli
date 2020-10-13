{% if site.public_api_components %}

{% include 'partials/aws_api_domain.tf' %}

resource "aws_apigatewayv2_api" "main_gateway" {
  name                       = "{{ site.identifier }}-api"
  protocol_type              = "HTTP"
}

resource "aws_apigatewayv2_route" "application" {
  api_id    = aws_apigatewayv2_api.main_gateway.id
  route_key = "$default"
}

resource "aws_apigatewayv2_deployment" "latest" {
  api_id      = aws_apigatewayv2_api.main_gateway.id
  description = "Stage for latest release"

  triggers = {
    redeployment = "{{ site.public_api_components_hash }}"
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_apigatewayv2_route.application,
  ]
}

resource "aws_apigatewayv2_stage" "latest" {
  name                 = "latest"
  api_id               = aws_apigatewayv2_api.main_gateway.id
  deployment_id        = aws_apigatewayv2_deployment.latest.id


  depends_on = [aws_apigatewayv2_deployment.latest]
}

resource "aws_apigatewayv2_deployment" "primary" {
  api_id      = aws_apigatewayv2_api.main_gateway.id
  description = "Stage for primary release"

  triggers = {
    redeployment = "{{ site.public_api_components_hash }}"
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [aws_apigatewayv2_deployment.latest]
}

resource "aws_apigatewayv2_stage" "primary" {
  name                  = "primary"
  api_id                = aws_apigatewayv2_api.main_gateway.id
  deployment_id         = aws_apigatewayv2_deployment.primary.id
  
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw_primary.arn
    # needs to be one line...
    format          = "{\"requestId\":\"$context.requestId\", \"ip\": \"$context.identity.sourceIp\", \"caller\":\"$context.identity.caller\", \"user\":\"$context.identity.user\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\", \"path\":\"$context.path\", \"status\":\"$context.status\", \"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\"}"
  }

  depends_on = [aws_apigatewayv2_deployment.primary]
}

resource "aws_cloudwatch_log_group" "api_gw_primary" {
  name              = "api_gw_stage_primary_access_logs"
  retention_in_days = 30
}

resource "aws_apigatewayv2_api_mapping" "custom_domain_mapping" {
  api_id      = aws_apigatewayv2_api.main_gateway.id
  stage       = aws_apigatewayv2_stage.primary.id
  domain_name = "{{ site.base_url }}"
}
{% endif %}
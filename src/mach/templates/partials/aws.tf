{% set aws = site.aws %}
provider "aws" {
  region  = "{{ aws.region }}"
  version = "~> 3.8.0"

  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role }}"
  }
}

{% for provider in aws.extra_providers %}
provider "aws" {
  alias   = "{{ provider.name }}"
  region  = "{{ provider.region }}"
  version = "~> 3.8.0"

  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role }}"
  }
}
{% endfor %}


// Hack so the api gateway is retriggered if a component changes
resource "null_resource" "lambda_changes" {
  triggers = {
    dependency_ids = [{% for component in site.api_gateway_changed_components %}module.{{ component.name }}.component_version{% if not loop.last %}, {% endif %}{% endfor %}]
  }
}

{% if site.aws.api_gateway %}
data "aws_apigatewayv2_api" "main_gateway" {
  name = "{{ site.aws.api_gateway }}"
}

{# WIP:
resource "aws_apigatewayv2_deployment" "latest" {
  api_id            = data.aws_apigatewayv2_api.main_gateway.id
  description       = "Stage for latest release ${null_resource.lambda_changes.triggers.dependency_ids}"
  stage_description = "Stage for latest release ${null_resource.lambda_changes.triggers.dependency_ids}"

  // https://github.com/hashicorp/terraform/issues/10674#issuecomment-290767062
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_stage" "latest" {
  name                 = "latest"
  api_id               = data.aws_apigatewayv2_api.main_gateway.id
  deployment_id        = aws_apigatewayv2_deployment.latest.id


  depends_on = [aws_apigatewayv2_deployment.latest]
}

resource "aws_apigatewayv2_deployment" "primary" {
  api_id            = data.aws_apigatewayv2_api.main_gateway.id
  description       = "Stage for latest release ${null_resource.lambda_changes.triggers.dependency_ids}"
  stage_description = "Stage for latest release ${null_resource.lambda_changes.triggers.dependency_ids}"

  // https://github.com/hashicorp/terraform/issues/10674#issuecomment-290767062
  lifecycle {
    create_before_destroy = true
  }

  depends_on = [aws_apigatewayv2_deployment.latest]
}

resource "aws_apigatewayv2_stage" "primary" {
  name                  = "primary"
  api_id                = data.aws_apigatewayv2_api.main_gateway.id
  deployment_id         = aws_apigatewayv2_deployment.primary.id
  
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw_primary.arn
    # needs to be one line...
    format          = "{\"requestId\":\"$context.requestId\", \"ip\": \"$context.identity.sourceIp\", \"caller\":\"$context.identity.caller\", \"user\":\"$context.identity.user\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\", \"resourcePath\":\"$context.resourcePath\", \"status\":\"$context.status\", \"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\"}"
  }

  depends_on = [aws_apigatewayv2_deployment.primary]
}

resource "aws_cloudwatch_log_group" "api_gw_primary" {
  name              = "api_gw_stage_primary_access_logs"
  retention_in_days = 30
}

resource "aws_apigatewayv2_base_path_mapping" "custom_domain_mapping" {
  api_id      = data.aws_apigatewayv2_api.main_gateway.id
  stage_name  = "primary"
  domain_name = "{{ site.base_url }}"
}

resource "aws_apigatewayv2_method_settings" "primary" {
  api_id      = data.aws_apigatewayv2_api.main_gateway.id
  stage_name  = aws_apigatewayv2_stage.primary.stage_name
  method_path = "*/*"

  settings {
    logging_level      = "ERROR"
    data_trace_enabled = true
    metrics_enabled    = true
  }

  depends_on = [aws_apigatewayv2_deployment.primary]
}
#}
{% endif %}
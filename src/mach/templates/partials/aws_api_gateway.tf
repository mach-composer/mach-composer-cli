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

resource "aws_apigatewayv2_deployment" "default" {
  api_id      = aws_apigatewayv2_api.main_gateway.id
  description = "Stage for default release"

  triggers = {
    redeployment = sha1(join(",", list(
      {% for component in site.public_api_components %}
      module.{{ component.name }}.component_version,
      {% endfor %}
    )))
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    {% for component in site.public_api_components %}
    module.{{ component.name }},
    {% endfor %}
  ]
}

resource "aws_apigatewayv2_stage" "default" {
  name                  = "$default"
  api_id                = aws_apigatewayv2_api.main_gateway.id
  deployment_id         = aws_apigatewayv2_deployment.default.id

  depends_on = [aws_apigatewayv2_deployment.default]
}

resource "aws_apigatewayv2_api_mapping" "custom_domain_mapping" {
  api_id      = aws_apigatewayv2_api.main_gateway.id
  stage       = aws_apigatewayv2_stage.default.id
  domain_name = "{{ site.base_url }}"
}
{% endif %}
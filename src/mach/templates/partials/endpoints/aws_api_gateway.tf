resource "aws_acm_certificate" "{{ endpoint_name|slugify }}" {
  domain_name       = "{{ endpoint_url }}"
  validation_method = "DNS"
}

resource "aws_route53_record" "new_acm_validation_site" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = tolist(aws_acm_certificate.{{ endpoint_name|slugify }}.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.{{ endpoint_name|slugify }}.domain_validation_options)[0].resource_record_type
  ttl     = 60
  records = [tolist(aws_acm_certificate.{{ endpoint_name|slugify }}.domain_validation_options)[0].resource_record_value]
}

# API Gateway
resource "aws_apigatewayv2_api" "{{ endpoint_name|slugify }}_gateway" {
  name                       = "{{ site.identifier }}-api"
  protocol_type              = "HTTP"
}

resource "aws_apigatewayv2_route" "{{ endpoint_name|slugify }}_application" {
  api_id    = aws_apigatewayv2_api.{{ endpoint_name|slugify }}_gateway.id
  route_key = "$default"
}

resource "aws_apigatewayv2_deployment" "{{ endpoint_name|slugify }}_default" {
  api_id      = aws_apigatewayv2_api.{{ endpoint_name|slugify }}_gateway.id
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

resource "aws_apigatewayv2_stage" "{{ endpoint_name|slugify }}_default" {
  name                  = "$default"
  api_id                = aws_apigatewayv2_api.{{ endpoint_name|slugify }}_gateway.id
  deployment_id         = aws_apigatewayv2_deployment.default.id

  depends_on = [aws_apigatewayv2_deployment.default]
}

# Route53 mappings
resource "aws_apigatewayv2_domain_name" "{{ endpoint_name|slugify }}" {
  domain_name = "{{ endpoint_url }}"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.{{ endpoint_name|slugify }}.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_route53_record" "{{ endpoint_name|slugify }}" {
  name    = aws_apigatewayv2_domain_name.main.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.main.id

  alias {
    name                   = aws_apigatewayv2_domain_name.main.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.main.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_api_mapping" "{{ endpoint_name|slugify }}" {
  api_id      = aws_apigatewayv2_api.{{ endpoint_name|slugify }}_gateway.id
  stage       = aws_apigatewayv2_stage.default.id
  domain_name = "{{ endpoint_url }}"
}
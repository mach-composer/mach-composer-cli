resource "aws_acm_certificate" "{{ endpoint.key|slugify }}" {
  domain_name       = "{{ endpoint.url }}"
  validation_method = "DNS"
}

resource "aws_route53_record" "{{ endpoint.key|slugify }}_acm_validation" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = tolist(aws_acm_certificate.{{ endpoint.key|slugify }}.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.{{ endpoint.key|slugify }}.domain_validation_options)[0].resource_record_type
  ttl     = 60
  records = [tolist(aws_acm_certificate.{{ endpoint.key|slugify }}.domain_validation_options)[0].resource_record_value]
}

# API Gateway
resource "aws_apigatewayv2_api" "{{ endpoint.key|slugify }}_gateway" {
  name                       = "{{ site.identifier }}-api"
  protocol_type              = "HTTP"
}

resource "aws_apigatewayv2_route" "{{ endpoint.key|slugify }}_application" {
  api_id    = aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.id
  route_key = "$default"
}

resource "aws_apigatewayv2_deployment" "{{ endpoint.key|slugify }}_default" {
  api_id      = aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.id
  description = "Stage for default release"

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    {% for component in endpoint.components %}
    module.{{ component.name }},
    {% endfor %}
  ]
}

resource "aws_apigatewayv2_stage" "{{ endpoint.key|slugify }}_default" {
  name                  = "$default"
  description           = "Stage for default release"
  api_id                = aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.id
  deployment_id         = aws_apigatewayv2_deployment.{{ endpoint.key|slugify }}_default.id
}

# Route53 mappings
resource "aws_apigatewayv2_domain_name" "{{ endpoint.key|slugify }}" {
  domain_name = "{{ endpoint.url }}"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.{{ endpoint.key|slugify }}.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_route53_record" "{{ endpoint.key|slugify }}" {
  name    = aws_apigatewayv2_domain_name.{{ endpoint.key|slugify }}.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.main.id

  alias {
    name                   = aws_apigatewayv2_domain_name.{{ endpoint.key|slugify }}.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.{{ endpoint.key|slugify }}.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_api_mapping" "{{ endpoint.key|slugify }}" {
  api_id      = aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.id
  stage       = aws_apigatewayv2_stage.{{ endpoint.key|slugify }}_default.id
  domain_name = "{{ endpoint.url }}"
}

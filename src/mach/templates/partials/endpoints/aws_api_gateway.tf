# API Gateway
resource "aws_apigatewayv2_api" "{{ endpoint.key|slugify }}_gateway" {
  name                       = "{{ site.identifier }}-{{ endpoint.key|slugify }}-api"
  protocol_type              = "HTTP"
}

resource "aws_apigatewayv2_route" "{{ endpoint.key|slugify }}_application" {
  api_id    = aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.id
  route_key = "$default"
}

resource "aws_apigatewayv2_stage" "{{ endpoint.key|slugify }}_default" {
  name                  = "$default"
  description           = "Stage for default release"
  api_id                = aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.id
  auto_deploy           = true

  {% if endpoint.throttling_burst_limit or endpoint.throttling_burst_limit -%}
  default_route_settings {
    {% if endpoint.throttling_burst_limit %}
    throttling_burst_limit = {{ endpoint.throttling_burst_limit }}
    {% endif %}
    {% if endpoint.throttling_rate_limit %}
    throttling_rate_limit = {{ endpoint.throttling_rate_limit }}
    {% endif %}
  }{% endif %}

  depends_on = [
    {% for component in endpoint.components %}
    module.{{ component.name }},
    {% endfor %}
  ]
}

{% if endpoint.url %}
resource "aws_acm_certificate" "{{ endpoint.key|slugify }}" {
  domain_name       = "{{ endpoint.url }}"
  validation_method = "DNS"
}

resource "aws_route53_record" "{{ endpoint.key|slugify }}_acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.{{ endpoint.key|slugify }}.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.{{ endpoint.zone|slugify }}.zone_id
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
  zone_id = data.aws_route53_zone.{{ endpoint.zone|slugify }}.id

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
{% endif %}
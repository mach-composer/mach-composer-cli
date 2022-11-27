# API Gateway
resource "aws_apigatewayv2_api" "{{ endpoint.Key|slugify }}_gateway" {
  name                       = "{{ siteName }}-{{ endpoint.Key|slugify }}-api"
  protocol_type              = "HTTP"
}

resource "aws_apigatewayv2_route" "{{ endpoint.Key|slugify }}_application" {
  api_id    = aws_apigatewayv2_api.{{ endpoint.Key|slugify }}_gateway.id
  route_key = "$default"
}

resource "aws_apigatewayv2_stage" "{{ endpoint.Key|slugify }}_default" {
  name                  = "$default"
  description           = "Stage for default release"
  api_id                = aws_apigatewayv2_api.{{ endpoint.Key|slugify }}_gateway.id
  auto_deploy           = true

  {% if endpoint.ThrottlingRateLimit or endpoint.ThrottlingBurstLimit -%}
  default_route_settings {
    {% if endpoint.ThrottlingBurstLimit %}
    throttling_burst_limit = {{ endpoint.ThrottlingBurstLimit }}
    {%- endif -%}
    {% if endpoint.ThrottlingRateLimit %}
    throttling_rate_limit = {{ endpoint.ThrottlingRateLimit }}
    {% endif %}
  }{% endif %}

  depends_on = [
    {% for component in endpoint.AssociatedComponents %}
    module.{{ component }},
    {%- endfor %}
  ]
}

{% if endpoint.URL %}
  resource "aws_acm_certificate" "{{ endpoint.Key|slugify }}" {
    domain_name       = "{{ endpoint.URL }}"
    validation_method = "DNS"

    {% if endpoint.AWS.enable_cdn %}
      provider = aws.mach-cf-us-east-1
    {% endif %}
  }

  resource "aws_route53_record" "{{ endpoint.Key|slugify }}_acm_validation" {
    for_each = {
      for dvo in aws_acm_certificate.{{ endpoint.Key|slugify }}.domain_validation_options : dvo.domain_name => {
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
    zone_id         = data.aws_route53_zone.{{ endpoint.Zone|slugify }}.zone_id
  }

  {% if endpoint.EnableCDN %}
    resource "aws_cloudfront_distribution" "{{ endpoint.key|slugify }}" {
      origin {
        origin_id   = "api-gateway"
        domain_name = replace(aws_apigatewayv2_api.{{ endpoint.Key|slugify }}_gateway.api_endpoint, "https://", "")

        custom_origin_config {
          http_port              = 80
          https_port             = 443
          origin_protocol_policy = "https-only"
          origin_ssl_protocols   = ["TLSv1.2"]
        }
      }

      aliases             = ["{{ endpoint.URL }}"]
      enabled             = true
      wait_for_deployment = false

      default_cache_behavior {
        allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
        cached_methods   = ["GET", "HEAD"]
        target_origin_id = "api-gateway"

        forwarded_values {
          query_string = true
          cookies {
            forward = "all"
          }
        }

        viewer_protocol_policy = "redirect-to-https"
        min_ttl                = 0
        default_ttl            = 0
        max_ttl                = 0
      }

      price_class = "PriceClass_200"

      restrictions {
        geo_restriction {
          restriction_type = "none"
        }
      }

      viewer_certificate {
        acm_certificate_arn      = aws_acm_certificate.{{ endpoint.Key|slugify }}.arn
        ssl_support_method       = "sni-only"
        minimum_protocol_version = "TLSv1"

      }
    }

    resource "aws_route53_record" "{{ endpoint.Key|slugify }}_cloudfront" {
      zone_id = data.aws_route53_zone.{{ endpoint.Zone|slugify }}.zone_id
      name    = "{{ endpoint.URL }}"
      type    = "A"

      alias {
        name                   = aws_cloudfront_distribution.{{ endpoint.Key|slugify }}.domain_name
        zone_id                = aws_cloudfront_distribution.{{ endpoint.Key|slugify }}.hosted_zone_id
        evaluate_target_health = false
      }
    }

  {% else %}
  # Route53 mappings
  resource "aws_apigatewayv2_domain_name" "{{ endpoint.Key|slugify }}" {
    domain_name = "{{ endpoint.URL }}"

    domain_name_configuration {
      certificate_arn = aws_acm_certificate.{{ endpoint.Key|slugify }}.arn
      endpoint_type   = "REGIONAL"
      security_policy = "TLS_1_2"
    }
  }

  resource "aws_route53_record" "{{ endpoint.Key|slugify }}" {
    name    = aws_apigatewayv2_domain_name.{{ endpoint.Key|slugify }}.domain_name
    type    = "A"
    zone_id = data.aws_route53_zone.{{ endpoint.Zone|slugify }}.id

    alias {
      name                   = aws_apigatewayv2_domain_name.{{ endpoint.Key|slugify }}.domain_name_configuration[0].target_domain_name
      zone_id                = aws_apigatewayv2_domain_name.{{ endpoint.Key|slugify }}.domain_name_configuration[0].hosted_zone_id
      evaluate_target_health = false
    }
  }

  resource "aws_apigatewayv2_api_mapping" "{{ endpoint.Key|slugify }}" {
    api_id      = aws_apigatewayv2_api.{{ endpoint.Key|slugify }}_gateway.id
    stage       = aws_apigatewayv2_stage.{{ endpoint.Key|slugify }}_default.id
    domain_name = "{{ endpoint.URL }}"
  }

  {% endif %}
{% endif %}


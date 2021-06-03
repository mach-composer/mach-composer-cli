{% set aws = site.aws %}

provider "aws" {
  alias   = "mach-cf-us-east-1"
  region  = "us-east-1"

  {% if aws.deploy_role_name %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role_name }}"
  }
  {% endif %}
}

resource "aws_acm_certificate" "cloudfront_api_gateway" {
  domain_name       = {{ endpoint.url|tf }}
  validation_method = "DNS"

  provider = aws.mach-cf-us-east-1
}

resource "aws_route53_record" "cloudfront_api_gateway_acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.cloudfront_api_gateway.domain_validation_options : dvo.domain_name => {
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

resource "aws_acm_certificate_validation" "cloudfront_api_gateway" {
  certificate_arn         = aws_acm_certificate.cloudfront_api_gateway.arn
  validation_record_fqdns = [for record in aws_route53_record.cloudfront_api_gateway_acm_validation : record.fqdn]

  provider = aws.mach-cf-us-east-1
}

resource "aws_cloudfront_distribution" "api" {
  origin {
    origin_id   = "api-gateway"
    domain_name = replace(aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.api_endpoint, "https://", "")

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  aliases             = [{{ endpoint.url|tf }}]
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
    acm_certificate_arn      = aws_acm_certificate.cloudfront_api_gateway.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1"

  }
}

resource "aws_route53_record" "cloudfront_domain" {
  zone_id = data.aws_route53_zone.{{ endpoint.zone|slugify }}.zone_id
  name    = {{ endpoint.url|tf }}
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.api.domain_name
    zone_id                = aws_cloudfront_distribution.api.hosted_zone_id
    evaluate_target_health = false
  }
}


locals {
{% for endpoint in site.used_endpoints %}
  endpoint_url_{{ endpoint.key|slugify }} = {% if endpoint.url %}{{ endpoint.url|tf }}{% else %}aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.api_endpoint{% endif %}

{% endfor %}
}

output "endpoints" {
  value = {
  {% for endpoint in site.used_endpoints -%}
    {{ endpoint.key|slugify }}: local.endpoint_url_{{ endpoint.key|slugify }}
  {% endfor %}
  }
}

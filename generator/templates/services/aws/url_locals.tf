locals {
{% for endpoint in site.UsedEndpoints() %}
  endpoint_url_{{ endpoint.Key|slugify }} = {% if endpoint.URL %}{{ endpoint.URL|tf }}{% else %}aws_apigatewayv2_api.{{ endpoint.key|slugify }}_gateway.api_endpoint{% endif %}

{% endfor %}
}

output "endpoints" {
  value = {
  {% for endpoint in site.UsedEndpoints() -%}
    {{ endpoint.Key|slugify }}: local.endpoint_url_{{ endpoint.Key|slugify }}
  {% endfor %}
  }
}

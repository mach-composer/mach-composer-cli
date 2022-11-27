locals {
{% for endpoint in endpoints %}
  endpoint_url_{{ endpoint.Key|slugify }} = {% if endpoint.URL %}"{{ endpoint.URL }}"{% else %}aws_apigatewayv2_api.{{ endpoint.Key|slugify }}_gateway.api_endpoint{% endif %}
{% endfor %}
}

output "endpoints" {
  value = {
  {% for endpoint in endpoints -%}
    {{ endpoint.Key|slugify }}: local.endpoint_url_{{ endpoint.Key|slugify }}
  {% endfor %}
  }
}

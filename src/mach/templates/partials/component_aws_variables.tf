{% for component_endpoint, site_endpoint in component.endpoints.items() -%}
aws_endpoint_{{ component_endpoint|slugify }} = {
  url = local.endpoint_url_{{ site_endpoint|slugify }}
  api_gateway_id = aws_apigatewayv2_api.{{ site_endpoint|slugify }}_gateway.id
  api_gateway_execution_arn = aws_apigatewayv2_api.{{ site_endpoint|slugify }}_gateway.execution_arn
}
{% endfor %}

# Won't be prefixed; since in theory this could be used for other cloud integrations as well
{% if definition.artifacts %}
artifacts = {
  {% for name, artifact in definition.artifacts.items() -%}
  {{ name }} = "{{ artifact.filename }}"
  {% endfor %}
}
{% endif %}

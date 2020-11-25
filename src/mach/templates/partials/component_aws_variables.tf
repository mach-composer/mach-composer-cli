{%- if component.definition.endpoint %}
  api_gateway = aws_apigatewayv2_api.{{ component.definition.endpoint|slugify }}_gateway.id
  api_gateway_execution_arn = aws_apigatewayv2_api.{{ component.definition.endpoint|slugify }}_gateway.execution_arn
{% endif %}

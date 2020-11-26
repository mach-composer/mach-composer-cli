{%- if component.endpoint %}
  api_gateway = aws_apigatewayv2_api.{{ component.endpoint|slugify }}_gateway.id
  api_gateway_execution_arn = aws_apigatewayv2_api.{{ component.endpoint|slugify }}_gateway.execution_arn
{% endif %}

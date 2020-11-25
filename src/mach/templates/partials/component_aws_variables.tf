{%- if component.definition.endpoint %}
  api_gateway = aws_apigatewayv2_api.main_gateway.id
  api_gateway_execution_arn = aws_apigatewayv2_api.main_gateway.execution_arn
{% endif %}

{% if component.definition.has_public_api %}
api_gateway = aws_apigatewayv2_api.main_gateway.id
api_gateway_execution_arn = aws_apigatewayv2_api.main_gateway.execution_arn
{% endif %}
code_repository = "{{ site.aws.code_repository }}"

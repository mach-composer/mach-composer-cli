{% if component.definition.has_public_api %}
api_gateway = aws_apigatewayv2_api.main_gateway.id
{% endif %}
code_repository = "{{ site.aws.code_repository }}"
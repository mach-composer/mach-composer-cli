{% if component.definition.has_public_api %}
api_gateway = aws_apigatewayv2_api.main_gateway.id
{% endif %}
{% if component.definition.is_software_component %}
code_repository = "{{ site.aws.code_repository }}"
{% endif %}
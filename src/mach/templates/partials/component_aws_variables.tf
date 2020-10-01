{% if site.aws.api_gateway %}
api_gateway = aws_api_gateway_rest_api.main_gateway.arn
{% endif %}
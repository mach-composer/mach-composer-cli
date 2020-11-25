{% set aws = site.aws %}
provider "aws" {
  region  = "{{ aws.region }}"
  version = "~> 3.8.0"
  {% if aws.deploy_role %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role }}"
  }
  {% endif %}
}

{% for provider in aws.extra_providers %}
provider "aws" {
  alias   = "{{ provider.name }}"
  region  = "{{ provider.region }}"
  version = "~> 3.8.0"

  {% if aws.deploy_role %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role }}"
  }
  {% endif %}
}
{% endfor %}

{% if site.public_api_components %}
  {% include 'partials/endpoints/aws_domains.tf' %}

  {% for endpoint_name, endpoint_url in site.endpoints.items() %}
    {% include 'partials/endpoints/aws_api_gateway.tf' %}
  {% endfor %}
{% endif %}
{% set aws = site.aws %}
provider "aws" {
  region  = "{{ aws.region }}"
  version = "~> 3.8.0"
  {% if aws.deploy_role_arn %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role_arn }}"
  }
  {% endif %}
}

{% for provider in aws.extra_providers %}
provider "aws" {
  alias   = "{{ provider.name }}"
  region  = "{{ provider.region }}"
  version = "~> 3.8.0"

  {% if aws.deploy_role_arn %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role_arn }}"
  }
  {% endif %}
}
{% endfor %}

{% if site.used_endpoints %}
  {% include 'partials/endpoints/aws_domains.tf' %}

  {% for endpoint_name, endpoint_url in site.used_endpoints.items() %}
    {% include 'partials/endpoints/aws_api_gateway.tf' %}
    
  {% endfor %}
{% endif %}
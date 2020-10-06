{% set aws = site.aws %}
provider "aws" {
  region  = "{{ aws.region }}"
  version = "~> 3.8.0"

  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role }}"
  }
}

{% for provider in aws.extra_providers %}
provider "aws" {
  alias   = "{{ provider.name }}"
  region  = "{{ provider.region }}"
  version = "~> 3.8.0"

  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role }}"
  }
}
{% endfor %}

{% include 'partials/aws_api_gateway.tf' %}
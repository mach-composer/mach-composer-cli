{% set aws = site.aws %}
provider "aws" {
  region  = {{ aws.region|tf }}
  {% if aws.deploy_role_name %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role_name }}"
  }
  {% endif %}
  {% if aws.default_tags %}
  default_tags {
      tags = {{ aws.default_tags|tf }}
  }
  {% endif %}
}

{% for provider in aws.extra_providers %}
provider "aws" {
  alias   = {{ provider.name|tf }}
  region  = {{ provider.region|tf }}

  {% if aws.deploy_role_name %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role_name }}"
  }
  {% endif %}
  {% if provider.default_tags %}
  default_tags {
      tags = {{ provider.default_tags|tf }}
  }
  {% elif aws.default_tags %}
  default_tags {
      tags = {{ aws.default_tags|tf }}
  }
  {% endif %}
}
{% endfor %}

{% if site.has_cdn_endpoint %}
provider "aws" {
  alias   = "mach-cf-us-east-1"
  region  = "us-east-1"

  {% if aws.deploy_role_name %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.account_id }}:role/{{ aws.deploy_role_name }}"
  }
  {% endif %}
  {% if aws.default_tags %}
  default_tags {
      tags = {{ aws.default_tags|tf }}
  }
  {% endif %}
}
{% endif %}

{% if site.used_endpoints %}
  {% for zone in site.dns_zones %}
  data "aws_route53_zone" "{{ zone|slugify }}" {
    name = "{{ zone }}"
  }
  {% endfor %}

  {% for endpoint in site.used_endpoints %}
    {% include 'partials/endpoints/aws_endpoints.tf' %}
  {% endfor %}
  {% include 'partials/endpoints/aws_url_locals.tf' %}
{% endif %}

locals {
  tags = {
    Site = "{{ site.identifier }}"
    Environment = {{ general_config.environment|tf }}
  }
}

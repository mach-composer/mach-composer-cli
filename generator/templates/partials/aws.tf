provider "aws" {
  region  = {{ aws.Region|tf }}
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

{% if site.UsedEndpoints() %}
  {% for zone in site.dns_zones %}
  data "aws_route53_zone" "{{ zone|slugify }}" {
    name = "{{ zone }}"
  }
  {% endfor %}

  {% for endpoint in site.UsedEndpoints() %}
    {% include "./endpoints/aws_endpoints.tf" %}
  {% endfor %}
  {% include "./endpoints/aws_url_locals.tf" %}
{% endif %}

locals {
  tags = {
    Site = "{{ site.Identifier }}"
    Environment = {{ global.Environment|tf }}
  }
}

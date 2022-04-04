# Configuring AWS

provider "aws" {
  region  = {{ aws.Region|tf }}
  {% if aws.DeployRoleName %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.AccountID }}:role/{{ aws.DeployRoleName }}"
  }
  {% endif %}
  {% if aws.DefaultTags %}
  default_tags {
      tags = {{ aws.DefaultTags|tf }}
  }
  {% endif %}
}

{% for provider in aws.ExtraProviders %}
provider "aws" {
  alias   = {{ provider.Name|tf }}
  region  = {{ provider.Region|tf }}

  {% if aws.DeployRoleName %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.AccountID }}:role/{{ aws.DeployRoleName }}"
  }
  {% endif %}
  {% if provider.DefaultTags %}
  default_tags {
      tags = {{ provider.DefaultTags|tf }}
  }
  {% elif aws.DefaultTags %}
  default_tags {
      tags = {{ aws.DefaultTags|tf }}
  }
  {% endif %}
}
{% endfor %}

{% if site.HasCdnEndpoint() %}
provider "aws" {
  alias   = "mach-cf-us-east-1"
  region  = "us-east-1"

  {% if aws.DeployRoleName %}
  assume_role {
    role_arn = "arn:aws:iam::{{ aws.AccountID }}:role/{{ aws.DeployRoleName }}"
  }
  {% endif %}
  {% if aws.DefaultTags %}
  default_tags {
      tags = {{ aws.DefaultTags|tf }}
  }
  {% endif %}
}
{% endif %}

{% if site.UsedEndpoints() %}
  {% for zone in site.DnsZones() %}
  data "aws_route53_zone" "{{ zone|slugify }}" {
    name = "{{ zone }}"
  }
  {% endfor %}

  {% for endpoint in site.UsedEndpoints() %}
    {% include "./endpoints.tf" %}
  {% endfor %}
  {% include "./url_locals.tf" %}
{% endif %}

locals {
  tags = {
    Site = "{{ site.Identifier }}"
    Environment = {{ global.Environment|tf }}
  }
}

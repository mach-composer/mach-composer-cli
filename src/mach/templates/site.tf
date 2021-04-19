# This file is auto-generated by MACH composer
# Site: {{ site.identifier }}

terraform {
  {% if general_config.terraform_config.azure_remote_state %}
  {% set azure_config = general_config.terraform_config.azure_remote_state %}
  backend "azurerm" {
    resource_group_name  = "{{ azure_config.resource_group }}"
    storage_account_name = "{{ azure_config.storage_account }}"
    container_name       = "{{ azure_config.container_name }}"
    key                  = "{{ azure_config.state_folder}}/{{ site.identifier }}"
  }
  {% elif general_config.terraform_config.aws_remote_state %}
  {% set aws_config = general_config.terraform_config.aws_remote_state %}
  backend "s3" {
    bucket         = "{{ aws_config.bucket}}"
    key            = "{{ aws_config.key_prefix}}/{{ site.identifier }}"
    region         = "{{ aws_config.region }}"
    {% if aws_config.role_arn %}
    role_arn       = "{{ aws_config.role_arn }}"
    {% endif %}
    {% if aws_config.lock_table %}
    dynamodb_table = "{{ aws_config.lock_table }}"
    {% endif %}
    encrypt        = {% if aws_config.encrypt %}true{% else %}false{% endif %}

  }
  {% elif general_config.terraform_config.gcp_remote_state %}
  {% set gcp_config = general_config.terraform_config.gcp_remote_state %}
  backend "gcs" {
    bucket         = "{{ gcp_config.bucket}}"
    prefix         = "{{ gcp_config.prefix}}/{{ site.identifier }}"
  }
  {% endif %}
}

terraform {
  required_providers {
    {% if site.aws %}
    aws = {
      version = "~> {{ general_config.terraform_config.providers.aws or '3.28.0' }}"
    }
    {% endif %}
    
    {% if site.azure %}
    azurerm = {
      version = "~> {{ general_config.terraform_config.providers.azure or '2.47.0' }}"
    }
    {% endif %}
    {% if site.commercetools %}
    commercetools = {
      source = "labd/commercetools"
      version = "~> {{ general_config.terraform_config.providers.commercetools or '0.25.3' }}"
    }
    {% endif %}
    {% if site.contentful %}
    contentful = {
      source = "labd/contentful"
      version = "~> {{ general_config.terraform_config.providers.contentful or '0.1.0' }}"
    }
    {% endif %}
    {% if site.amplience %}
    amplience = {
      source = "labd/amplience"
      version = "~> {{ general_config.terraform_config.providers.amplience or '0.1.0' }}"
    }
    {% endif %}
    {% if general_config.sentry.managed %}
    sentry = {
      source = "jianyuan/sentry"
      version = "~> {{ general_config.terraform_config.providers.sentry or '0.6.0' }}"
    }
    {% endif %}
  }
}

{% if general_config.sentry.managed %}
provider "sentry" {
  token = "{{ general_config.sentry.auth_token }}"
  base_url = "{% if general_config.sentry.base_url %}{{ general_config.sentry.base_url }}{% else %}https://sentry.io/api/{% endif %}"
}
{% endif %}

{% if site.commercetools %}{% include 'partials/commercetools.tf' %}{% endif %}
{% if site.contentful %}{% include 'partials/contentful.tf' %}{% endif %}
{% if site.amplience %}{% include 'partials/amplience.tf' %}{% endif %}

{% if site.aws %}{% include 'partials/aws.tf' %}{% endif %}
{% if site.azure %}{% include 'partials/azure.tf' %}{% endif %}
{% if site.gcp %}{% include 'partials/gcp.tf' %}{% endif %}

{% include 'partials/components.tf' %}

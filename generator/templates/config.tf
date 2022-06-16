terraform {
  {% if global.TerraformConfig.AzureRemoteState %}
  {%- set azure_config = global.TerraformConfig.AzureRemoteState -%}
  backend "azurerm" {
    resource_group_name  = {{ azure_config.ResourceGroup|tf }}
    storage_account_name = {{ azure_config.StorageAccount|tf }}
    container_name       = {{ azure_config.ContainerName|tf }}
    key                  = "{{ azure_config.StateFolder}}/{{ site.Identifier }}"
  }
  {% elif global.TerraformConfig.AwsRemoteState %}
  {%- set aws_config = global.TerraformConfig.AwsRemoteState -%}
  backend "s3" {
    bucket         = {{ aws_config.Bucket|tf }}
    key            = "{{ aws_config.KeyPrefix}}/{{ site.Identifier }}"
    region         = {{ aws_config.Region|tf }}
    {% if aws_config.RoleARN -%}
    role_arn       = {{ aws_config.RoleARN|tf }}
    {% endif -%}
    {%- if aws_config.LockTable -%}
    dynamodb_table = {{ aws_config.LockTable|tf }}
    {% endif -%}
    encrypt        = {% if aws_config.Encrypt %}true{% else %}false{% endif %}
  }
  {%- endif %}

  required_providers {
    {%- if site.AWS %}
    aws = {
      version = "{{ global.TerraformConfig.Providers.AWS|render_tf_provider:"3.74.1" }}"
    }
    {% endif -%}

    {%- if site.Azure %}
    azurerm = {
      version = "{{ global.TerraformConfig.Providers.Azure|render_tf_provider:"2.99.0" }}"
    }
    {% endif -%}

    {%- if site.Commercetools %}
    commercetools = {
      source = "labd/commercetools"
      version = "{{ global.TerraformConfig.Providers.Commercetools|render_tf_provider:'0.30.0' }}"
    }
    {% endif -%}

    {%- if site.Contentful %}
    contentful = {
      source = "labd/contentful"
      version = "{{ global.TerraformConfig.Providers.Contentful|render_tf_provider:'0.1.0' }}"
    }
    {% endif -%}

    {%- if site.Amplience %}
    amplience = {
      source = "labd/amplience"
      version = "{{ global.TerraformConfig.Providers.Amplience|render_tf_provider:'0.3.7' }}"
    }
    {% endif -%}

    {%- if global.SentryConfig.AuthToken %}
    sentry = {
      source = "jianyuan/sentry"
      version = "{{ global.TerraformConfig.Providers.Sentry|render_tf_provider:'0.6.0' }}"
    }
    {% endif -%}

    {%- if variables.Encrypted %}
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
    {%- endif %}
  }
}


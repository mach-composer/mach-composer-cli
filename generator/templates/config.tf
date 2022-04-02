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

  {% include "partials/providers.tf" %}
}


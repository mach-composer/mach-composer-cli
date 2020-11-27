{% set definition = component.definition %}

{% if "sentry" in component.integrations and general_config.sentry.managed %}
resource "sentry_key" "{{ component.name }}" {
  organization = "{{ general_config.sentry.organization }}"
  project = "{{ general_config.sentry.project }}"
  name = "{{ site.identifier }}_{{ component.name }}"
  {% if component.sentry.rate_limit_window %}
  rate_limit_window = {{ component.sentry.rate_limit_window }}
  {% endif %}
  {% if component.sentry.rate_limit_count %}
  rate_limit_count = {{ component.sentry.rate_limit_count }}
  {% endif %}
}
{% endif %}

module "{{ component.name }}" {
  source            = "{{ definition.source }}{% if definition.use_version_reference %}?ref={{ definition.version }}{% endif %}"

  {% if component.has_cloud_integration %}
  component_version       = "{{ definition.version }}"
  environment             = "{{ general_config.environment }}"
  site                    = "{{ site.identifier }}"

  variables = {
    {% for key, value in component.variables.items() %}
    {{ key }} = {{ value|component_value }}
    {% endfor %}
  }

  secrets = {
    {% for key, value in component.secrets.items() %}
    {{ key }} = {{ value|component_value }}
    {% endfor %}
  }
  {% endif %}

  {% if "azure" in component.integrations %}
  {% include 'partials/component_azure_variables.tf' %}
  {% endif %}

  {% if "aws" in component.integrations %}
  {% include 'partials/component_aws_variables.tf' %}
  {% endif %}

  {% if "sentry" in component.integrations %}
  sentry_dsn              = {% if general_config.sentry.managed %}sentry_key.{{ component.name }}.dsn_secret{% else %}"{{ component.sentry.dsn }}"{% endif %}
  {% endif %}

  {% if "commercetools" in component.integrations %}
    ct_project_key    = "{{ site.commercetools.project_key }}"
    ct_api_url        = "{{ site.commercetools.api_url }}"
    ct_auth_url       = "{{ site.commercetools.token_url }}"

    stores {
      {% for store in site.commercetools.stores %}
      {{ store.key }} =  {
        key = "{{ store.key }}"
        variables = {
          {% for key, value in component.store_variables.get(store.key, {}).items() %}
          {{ key }} = {{ value|component_value }}
          {% endfor %}
        }
        secrets = {
          {% for key, value in component.store_secrets.get(store.key, {}).items() %}
          {{ key }} = {{ value|component_value }}
          {% endfor %}
        }
      }
      {% endfor %}
  }

  {% endif %}

  {% if "contentful" in component.integrations %}
    contentful_space_id = contentful_space.space.id
  {% endif %}

  {% if "amplience" in component.integrations %}
    amplience_client_id = "{{ site.amplience.client_id }}"
    amplience_client_secret = "{{ site.amplience.client_secret }}"
    amplience_hub_id = "{{ site.amplience.hub_id }}"
  {% endif %}

  providers = {
    {% if "azure" in component.integrations %}azurerm = azurerm{% endif %}
    {% if "aws" in component.integrations %}
      aws = aws
      {% for provider in site.aws.extra_providers %}
      aws.{{ provider.name }} = aws.{{ provider.name }}
      {% endfor %}
    {% endif %}
  }

  depends_on = [
    {% if site.aws and component.endpoint %}
    aws_apigatewayv2_api.{{ component.endpoint|slugify }}_gateway,
    {% endif %}
    {% if site.commercetools %}
    null_resource.commercetools,
    {% endif %}
  ]
}

{% if site.azure and component.has_cloud_integration %}
# see https://docs.microsoft.com/en-us/azure/azure-functions/functions-deployment-technologies#trigger-syncing
# this updates the functionapp in case of any changes.
data "external" "sync_triggers_{{ component.name }}" {
  program = ["bash", "-c", "az rest --method post --uri 'https://management.azure.com/subscriptions/${local.subscription_id}/resourceGroups/${local.resource_group_name}/providers/Microsoft.Web/sites/${module.{{ component.name }}.app_service_name}/syncfunctiontriggers?api-version=2016-08-01'"]

  # need to make sure this runs after the module
  depends_on = [
    module.{{ component.name }}.app_service_name
  ]
}
{% endif %}

{% set definition = component.definition %}

{% if "sentry" in component.integrations and general_config.sentry.managed %}
resource "sentry_key" "{{ component.name }}" {
  organization = {{ general_config.sentry.organization|tf }}
  project = {{ component.sentry.project|tf }}
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

  {% if component.has_cloud_integration or component.variables %}
  variables = {
    {% for key, value in component.variables.items() %}
    {{ key }} = {{ value|tfvalue }}
    {% endfor %}
  }
  {% endif %}
  {% if component.has_cloud_integration or component.secrets %}
  secrets = {
    {% for key, value in component.secrets.items() %}
    {{ key }} = {{ value|tfvalue }}
    {% endfor %}
  }
  {% endif %}

  {% if component.has_cloud_integration %}
  component_version       = {{ definition.version|tf }}
  environment             = {{ general_config.environment|tf }}
  site                    = "{{ site.identifier }}"
  tags                    = local.tags
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
    ct_project_key    = {{ site.commercetools.project_key|tf }}
    ct_api_url        = {{ site.commercetools.api_url|tf }}
    ct_auth_url       = {{ site.commercetools.token_url|tf }}

    ct_stores = {
      {% for store in site.commercetools.stores %}
      {{ store.key }} =  {
        key = {{ store.key|tf }}
        variables = {
          {% for key, value in component.store_variables.get(store.key, {}).items() %}
          {{ key }} = {{ value|tfvalue }}
          {% endfor %}
        }
        secrets = {
          {% for key, value in component.store_secrets.get(store.key, {}).items() %}
          {{ key }} = {{ value|tfvalue }}
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
    amplience_client_id = {{ site.amplience.client_id|tf }}
    amplience_client_secret = {{ site.amplience.client_secret|tf }}
    amplience_hub_id = {{ site.amplience.hub_id|tf }}
  {% endif %}

  {% if "apollo_federation" in component.integrations %}
    apollo_federation = {
      api_key       = {{ site.apollo_federation.api_key|tf }}
      graph         = {{ site.apollo_federation.graph|tf }}
      graph_variant = {{ site.apollo_federation.graph_variant|tf }}
    }
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
    {% if site.azure and component.azure.service_plan %}
    {% if component.azure.service_plan == 'default' %}
    azurerm_app_service_plan.functionapps,{% else %}
    azurerm_app_service_plan.functionapps_{{ component.azure.service_plan }},{% endif %}
    {% endif %}
    {% if site.commercetools %}
    null_resource.commercetools,
    {% endif %}
  ]
}

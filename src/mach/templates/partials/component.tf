{% set definition = component.definition %}

module "{{ component.name }}" {
  source            = "{{ definition.source }}{% if definition.use_version_reference %}?ref={{ definition.version }}{% endif %}"
  
  {% if "azure" in component.integrations %}
    {% include 'partials/component_azure_variables.tf' %}
  {% endif %}
  
  {% if "aws" in component.integrations %}
    {% include 'partials/component_aws_variables.tf' %}
  {% endif %}

  {% if component.is_software_component %}
    component_version       = "{{ definition.version }}"
    environment             = "{{ general_config.environment }}"
    site                    = "{{ site.identifier }}"
    {% if general_config.sentry %}
    sentry_dsn              = "{{ general_config.sentry.dsn }}"
    {% endif %}

    variables = {
      {% for key, value in component.variables.items() %}
      {{ key }} = {{ value|component_value }}
      {% endfor %}
      
      {% if site.azure.front_door and component.has_public_api %}
      FRONTDOOR_ID = azurerm_frontdoor.app-service.header_frontdoor_id
      {% endif %}

      {% if site.commercetools %}
        {# BACKWARDS COMPATABILITY #}
        CT_API_URL = "{{ site.commercetools.api_url }}"
        CT_AUTH_URL = "{{ site.commercetools.token_url }}"
        CT_PROJECT_KEY = "{{ site.commercetools.project_key }}"
      {% endif %}
    }

    {# TODO: See if we can merge variables and environment_variables #}
    environment_variables = {
      {% filter indent(width=4) %}
      {% for key, value in component.variables.items() %}
        {{ key }} = {{ value|component_value }}
      {% endfor %}

      {% if site.commercetools and site.commercetools.stores %}
          STORES = "{% for store in site.commercetools.stores %}{{ store.key }}{% if not loop.last %},{% endif %}{% endfor %}"
          {% if site.commercetools.stores|length == 1 %}
              DEFAULT_STORE = "{{ site.commercetools.stores[0].key }}"
              STORE         = "{{ site.commercetools.stores[0].key }}"
          {% endif %}
      {% endif %}
      {% endfilter %}
    }

    secrets = {
    {% for key, value in component.secrets.items() %}
        {{ key }} = {{ value|component_value }}
    {% endfor %}
    }

    {% if "commercetools" in component.integrations %}
      ct_project_key    = "{{ site.commercetools.project_key }}"
      # ct_api_url        = "{{ site.commercetools.api_url }}"
      # ct_auth_url       = "{{ site.commercetools.token_url }}"
    {% endif %}
  {% endif %}

  {% if "contentful" in component.integrations %}
    contentful_space_id = contentful_space.space.id
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
    {% if site.aws and component.has_public_api %}aws_apigatewayv2_api.main_gateway,{% endif %}
    {% if "commercetools" in component.integrations %}commercetools_project_settings.project,{% endif %}
  ]
}

{% if site.azure and component.is_software_component %}
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
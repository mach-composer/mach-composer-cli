{% set definition = component.definition %}

module "{{ component.name }}" {
  source            = "{{ definition.source }}{% if definition.use_version_reference %}?ref={{ definition.version }}{% endif %}"
  
  {% if component.is_software_component %}
  
  {% if site.azure %}
  {% include 'partials/component_azure_variables.tf' %}
  {% elif site.aws %}
  {% include 'partials/component_aws_variables.tf' %}
  {% endif %}

  component_version       = "{{ definition.version }}"
  environment             = "{{ general_config.environment }}"
  site                    = "{{ site.identifier }}"
  
  {% if site.commercetools %}
  ct_project_key          = "{{ site.commercetools.project_key }}"
  {% endif %}

  {% if general_config.sentry %}
  sentry_dsn              = "{{ general_config.sentry.dsn }}"
  {% endif %}

  variables = {
    {% for key, value in component.variables.items() %}
    {{ key }} = {{ value|component_value }}
    {% endfor %}

    CT_API_URL = "{{ site.commercetools.api_url }}"
    {# TODO: make token url / auth url consistent #}
    CT_AUTH_URL = "{{ site.commercetools.token_url }}"
    CT_PROJECT_KEY = "{{ site.commercetools.project_key }}"
    
    {% if site.azure.front_door and component.has_public_api %}
    FRONTDOOR_ID = azurerm_frontdoor.app-service.header_frontdoor_id
    {% endif %}
  }

  {# TODO: See if we can merge variables and environment_variables #}
  environment_variables = {
    {% filter indent(width=4) %}
    {% for key, value in component.variables.items() %}
      {{ key }} = {{ value|component_value }}
    {% endfor %}

    {% if site.stores %}
        STORES = "{% for store in site.stores %}{{ store.key }}{% if not loop.last %},{% endif %}{% endfor %}"
        {% if site.stores|length == 1 %}
            DEFAULT_STORE = "{{ site.stores[0].key }}"
            STORE         = "{{ site.stores[0].key }}"
        {% endif %}
    {% endif %}
    {% endfilter %}
  }

  secrets = {
  {% for key, value in component.secrets.items() %}
      {{ key }} = {{ value|component_value }}
  {% endfor %}
  }
  {% endif %}

  providers = {
    {% if site.commercetools %}commercetools = commercetools{% endif %}

    {% if site.azure %}azurerm = azurerm{% endif %}
    {% if site.aws %}aws = aws
    {% for provider in site.aws.extra_providers %}
    aws.{{ provider.name }} = aws.{{ provider.name }}
    {% endfor %}
    {% endif %}
  }

  {% if site.aws and component.has_public_api %}
  depends_on = [
    aws_apigatewayv2_api.main_gateway,
  ]
  {% endif %}
}

{% if site.azure and component.is_software_component %}
# see https://docs.microsoft.com/en-us/azure/azure-functions/functions-deployment-technologies#trigger-syncing
# this updates the functionapp in case of any changes.
data "external" "sync_triggers_{{ component.name }}" {
  program = ["bash", "-c", "az rest --method post --uri 'https://management.azure.com/subscriptions/${local.subscription_id}/resourceGroups/${azurerm_resource_group.main.name}/providers/Microsoft.Web/sites/${module.{{ component.name }}.app_service_name}/syncfunctiontriggers?api-version=2016-08-01'"]

  # need to make sure this runs after the module
  depends_on = [
    module.{{ component.name }}.app_service_name
  ]
}
{% endif %}
{% set definition = component.definition %}

module "{{ component.name }}" {
  source            = "{{ definition.source }}{% if definition.use_version_reference %}?ref={{ definition.version }}{% endif %}"
  
  {% if component.is_software_component %}
  # keep the same order as the module's variables!

  ### azure related
  short_name              = "{{ component.short_name }}"
  name_prefix             = local.name_prefix
  subscription_id         = local.subscription_id
  tenant_id               = local.tenant_id
  service_object_ids      = local.service_object_ids
  region                  = local.region
  resource_group_name     = azurerm_resource_group.main.name
  resource_group_location = azurerm_resource_group.main.location
  app_service_plan_id     = azurerm_app_service_plan.functionapps.id
  tags                    = local.tags

  ### functionality related
  component_version       = "{{ definition.version }}"
  environment             = "{{ general_config.environment }}"
  site                    = "{{ site.identifier }}"
  ct_project_key          = "{{ site.commercetools.project_key }}"
  
  {% if general_config.sentry %}
  sentry_dsn              = "{{ general_config.sentry.dsn }}"
  {% endif %}
  {% if site.azure.alert_group %}
  monitor_action_group_id = azurerm_monitor_action_group.alert_action_group.id
  {% endif %}
  
  # todo make simple jinja filter
  variables = {
    CT_API_URL = "{{ site.commercetools.api_url }}"
    # TODO: make token url / auth url consistent
    CT_AUTH_URL = "{{ site.commercetools.token_url }}"
    CT_PROJECT_KEY = "{{ site.commercetools.project_key }}"
    {% if site.azure.front_door and component.has_public_api %}
    FRONTDOOR_ID = azurerm_frontdoor.app-service.header_frontdoor_id
    {% endif %}
  {% for key, value in component.variables.items() %}
      {{ key }} = {{ value|component_value }}
  {% endfor %}
  }

  secrets = {
  {% for key, value in component.secrets.items() %}
      {{ key }} = {{ value|component_value }}
  {% endfor %}
  }

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
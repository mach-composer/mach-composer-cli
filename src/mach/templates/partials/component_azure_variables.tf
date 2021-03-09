### azure related
azure_short_name              = "{{ component.azure.short_name }}"
azure_name_prefix             = local.name_prefix
azure_subscription_id         = local.subscription_id
azure_tenant_id               = local.tenant_id
azure_region                  = local.region
azure_service_object_ids      = local.service_object_ids
azure_resource_group          = {
  name     = local.resource_group_name
  location = local.resource_group_location
}
{% if component.azure.service_plan -%}
azure_app_service_plan        = {
  id                  = azurerm_app_service_plan.{{ component.azure.service_plan|service_plan_resource_name }}.id
  name                = azurerm_app_service_plan.{{ component.azure.service_plan|service_plan_resource_name }}.name
  resource_group_name = azurerm_app_service_plan.{{ component.azure.service_plan|service_plan_resource_name }}.resource_group_name
}
{% endif %}
{% if site.azure.alert_group %}
azure_monitor_action_group_id = azurerm_monitor_action_group.alert_action_group.id
{% endif %}
{% for component_endpoint, site_endpoint in component.endpoints.items() -%}
azure_endpoint_{{ component_endpoint|slugify }} = {
  url = local.endpoint_url_{{ site_endpoint|slugify }}
  frontdoor_id = azurerm_frontdoor.app-service.header_frontdoor_id
}
{% endfor %}

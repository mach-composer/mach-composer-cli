### azure related
short_name              = "{{ component.short_name }}"
name_prefix             = local.name_prefix
subscription_id         = local.subscription_id
tenant_id               = local.tenant_id
region                  = local.region
service_object_ids      = local.service_object_ids
resource_group          = {
  name     = local.resource_group_name
  location = local.resource_group_location
}
{% if component.azure.service_plan -%}
app_service_plan        = {
  id = azurerm_app_service_plan.{{ component.azure.service_plan|service_plan_resource_name }}.id
  name = azurerm_app_service_plan.{{ component.azure.service_plan|service_plan_resource_name }}.name
}
{% endif %}
{% if site.azure.alert_group %}
monitor_action_group_id = azurerm_monitor_action_group.alert_action_group.id
{% endif %}
{% for component_endpoint, site_endpoint in component.endpoints.items() -%}
endpoint_{{ component_endpoint|slugify }} = {
  url = local.endpoint_url_{{ site_endpoint|slugify }}
  frontdoor_id = azurerm_frontdoor.app-service.header_frontdoor_id
}
{% endfor %}

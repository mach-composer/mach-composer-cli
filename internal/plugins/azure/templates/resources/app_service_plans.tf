{% for key, plan in site.Azure.ServicePlans %}

{% if plan.DedicatedResourceGroup %}
resource "azurerm_resource_group" "{{ key }}" {
  name     = "${local.resource_group_name}-{{ key }}"
  location = local.resource_group_location
  tags     = local.tags
}
{% endif %}

resource "azurerm_app_service_plan" "{{ key|service_plan_resource_name }}" {
  {% if key == "default" %}
  name                = "${local.name_prefix}-plan"
  {% else %}
  name                = "${local.name_prefix}-{{ key }}-plan"
  {% endif %}
  resource_group_name = {% if plan.DedicatedResourceGroup %}azurerm_resource_group.{{ key }}.name
  {% else %}local.resource_group_name
  {% endif %}
  location            = local.resource_group_location
  kind                = "{{ plan.Kind }}"
  reserved            = {% if plan.Kind|lower == 'windows' %}false{% else %}true
  {% endif %}
  per_site_scaling    = "{{ plan.PerSiteScaling }}"

  sku {
    tier = "{{ plan.Tier }}"
    size = "{{ plan.Size }}"
    {% if plan.capacity -%}
    capacity = {{ plan.Capacity }}
    {% endif %}
  }

  tags = local.tags
}
{% endfor %}

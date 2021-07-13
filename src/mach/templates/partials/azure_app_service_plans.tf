{% for key, plan in site.azure.service_plans.items() %}

{% if plan.dedicated_resource_group %}
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
  resource_group_name = {% if plan.dedicated_resource_group %}azurerm_resource_group.{{ key }}.name
  {% else %}local.resource_group_name
  {% endif %}
  location            = local.resource_group_location
  kind                = {{ plan.kind|tf }}
  reserved            = {% if plan.kind|lower == 'windows' %}false{% else %}true
  {% endif %}
  per_site_scaling    = {{ plan.per_site_scaling|tf }}

  sku {
    tier = {{ plan.tier|tf }}
    size = {{ plan.size|tf }}
    {% if plan.capacity -%}
    capacity = {{ plan.capacity|tf }}
    {% endif %}
  }

  tags = local.tags
}
{% endfor %}

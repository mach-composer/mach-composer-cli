{% for key, plan in site.azure.service_plans.items() %}

resource "azurerm_app_service_plan" "{{ key|service_plan_resource_name }}" {
  {% if key == "default" %}
  name                = "${local.name_prefix}-plan"
  {% else %}
  name                = "${local.name_prefix}-{{ key }}-plan"
  {% endif %}
  resource_group_name = local.resource_group_name
  location            = local.resource_group_location
  kind                = "{{ plan.kind }}"
  reserved            = {% if plan.kind|lower == 'windows' %}false{% else %}true{% endif %}

  sku {
    tier = "{{ plan.tier }}"
    size = "{{ plan.size }}"
    {% if plan.capacity -%}
    capacity = {{ plan.capacity }}
    {% endif %}
  }

  tags = local.tags
}
{% endfor %}

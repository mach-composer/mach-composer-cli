{% for key, plan in general_config.azure.service_plans.items() %}

{% if key == "default" %}
{% set resource_name = "functionapps" %}
{% set name_suffix = "-plan" %}
{% else %}
{% set resource_name = "functionapps_%s"|format(key) %}
{% set name_suffix = "-%s-plan"|format(key) %}
{% endif %}

resource "azurerm_app_service_plan" "{{ resource_name }}" {
  name                = "${local.name_prefix}{{ name_suffix }}"
  resource_group_name = local.resource_group_name
  location            = local.resource_group_location
  kind                = "{{ plan.kind }}"
  reserved            = {% if plan.kind|lower == 'linux' %}true{% else %}false{% endif %}

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

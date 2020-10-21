{% if site.azure and site.components %}
resource "azurerm_app_service_plan" "functionapps" {
  name                = format("%s-plan", local.name_prefix)
  resource_group_name = local.resource_group_name
  location            = local.resource_group_location
  kind                = "FunctionApp"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }

  tags = local.tags
}
{% endif %}

{% for component in site.components %}
{% include 'partials/component.tf' %}
{% endfor %}

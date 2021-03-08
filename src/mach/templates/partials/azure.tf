{% set azure = site.azure %}

provider "azurerm" {
  subscription_id = "{{ azure.subscription_id }}"
  tenant_id       = "{{ azure.tenant_id }}"
  skip_provider_registration = true
  features {}
}


locals {
  tenant_id                    = "{{ azure.tenant_id }}"
  region                       = "{{ azure.region }}"
  subscription_id              = "{{ azure.subscription_id }}"
  project_key                  = "{{ site.commercetools.project_key }}"

  region_short                 = "{{ azure.region|azure_region_short }}"
  name_prefix                  = format("{{ general_config.azure.resources_prefix }}{{ site.identifier| replace("dev", "d") | replace("tst", "t") | replace("prd", "p") }}-%s", local.region_short)

  service_object_ids           = {
      {% for key, value in azure.service_object_ids.items() %}
          {{ key }} = "{{ value }}"
      {% endfor %}
  }

  tags = {
    Site = "{{ site.identifier }}"
    Environment = "{{ general_config.environment }}"
  }
}

{% if azure.resource_group %}
data "azurerm_resource_group" "main" {
  name = "{{ azure.resource_group }}"
}  
{% else %}
resource "azurerm_resource_group" "main" {
  name     = format("%s-rg", local.name_prefix)
  location = "{{ azure.region|azure_region_long }}"
  tags = local.tags
}
{% endif %}

locals {
  {% if azure.resource_group %}
    resource_group_name = data.azurerm_resource_group.main.name
    resource_group_location = data.azurerm_resource_group.main.location
  {% else %}
    resource_group_name = azurerm_resource_group.main.name
    resource_group_location = azurerm_resource_group.main.location
  {% endif %}
}


{% if azure.alert_group %}
{% if azure.alert_group.logic_app %}
data "azurerm_logic_app_workflow" "alert_logic_app" {
  name                = "{{ azure.alert_group.logic_app_name }}"
  resource_group_name = "{{ azure.alert_group.logic_app_resource_group }}"
}
{% endif %}

resource "azurerm_monitor_action_group" "alert_action_group" {
  name                = "{{ site.identifier }}-{{ azure.alert_group.name }}"
  resource_group_name = azurerm_resource_group.main.name
  short_name          = "{{ azure.alert_group.name|replace(" ", "")|replace("-", "")|lower }}"

  {% for email in azure.alert_group.alert_emails %}
  email_receiver {
    name          = "{{ email }}"
    email_address = "{{ email }}"
  }
  {% endfor %}

  {% if azure.alert_group.logic_app %}
  logic_app_receiver {
      name                    = "Logic app receiver"
      resource_id             = data.azurerm_logic_app_workflow.alert_logic_app.id
      callback_url            = data.azurerm_logic_app_workflow.alert_logic_app.access_endpoint
      use_common_alert_schema = true
  }
  {% endif %}

  {% if azure.alert_group.webhook_url %}
  webhook_receiver {
    name                    = "alert_webhook"
    service_uri             = "{{ azure.alert_group.webhook_url }}"
    use_common_alert_schema = true
  }
  {% endif %}
}
{% endif %}

{% include 'partials/endpoints/azure_frontdoor.tf' %}
{% include 'partials/endpoints/azure_url_locals.tf' %}

{% include 'partials/azure_app_service_plans.tf' %}

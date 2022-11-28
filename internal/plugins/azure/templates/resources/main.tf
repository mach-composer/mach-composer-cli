provider "azurerm" {
  subscription_id = "{{ azure.SubscriptionID }}"
  tenant_id       = "{{ azure.TenantID }}"
  skip_provider_registration = true
  features {}
}

locals {
  tenant_id                    = "{{ azure.TenantID }}"
  region                       = "{{ azure.Region }}"
  subscription_id              = "{{ azure.SubscriptionID }}"

  region_short                 = "{{ azure.ShortRegionName() }}"
  name_prefix                  = format("{{ global.Azure.ResourcesPrefix }}{{ siteName|short_prefix }}-%s", local.region_short)

  service_object_ids           = {
      {% for key, value in azure.ServiceObjectIds %}
          {{ key }} = "{{ value }}"
      {% endfor %}
  }

  tags = {
    Site        = "{{ siteName }}"
    Environment = "{{ envName }}"
  }
}

{% if azure.ResourceGroup %}
data "azurerm_resource_group" "main" {
  name = "{{ azure.ResourceGroup }}"
}
{% else %}
resource "azurerm_resource_group" "main" {
  name     = format("%s-rg", local.name_prefix)
  location = "{{ azure.LongRegionName() }}"
  tags = local.tags
}
{% endif %}

locals {
  {% if azure.ResourceGroup %}
    resource_group_name = data.azurerm_resource_group.main.name
    resource_group_location = data.azurerm_resource_group.main.location
  {% else %}
    resource_group_name = azurerm_resource_group.main.name
    resource_group_location = azurerm_resource_group.main.location
  {% endif %}
}


{% if azure.AlertGroup %}
{% if azure.AlertGroup.LogicApp %}
data "azurerm_logic_app_workflow" "alert_logic_app" {
  name                = "{{ azure.AlertGroup.LogicAppName() }}"
  resource_group_name = "{{ azure.AlertGroup.LogicAppResourceGroup() }}"
}
{% endif %}

resource "azurerm_monitor_action_group" "alert_action_group" {
  name                = "{{ siteName }}-{{ azure.AlertGroup.Name }}"
  resource_group_name = azurerm_resource_group.main.name
  short_name          = "{{ azure.AlertGroup.Name|remove:","|remove:"-"|lower }}"

  {% for email in azure.AlertGroup.AlertEmails %}
  email_receiver {
    name          = "{{ email }}"
    email_address = "{{ email }}"
  }
  {% endfor %}

  {% if azure.AlertGroup.LogicApp %}
  logic_app_receiver {
      name                    = "Logic app receiver"
      resource_id             = data.azurerm_logic_app_workflow.alert_logic_app.id
      callback_url            = data.azurerm_logic_app_workflow.alert_logic_app.access_endpoint
      use_common_alert_schema = true
  }
  {% endif %}

  {% if azure.AlertGroup.WebhookURL %}
  webhook_receiver {
    name                    = "alert_webhook"
    service_uri             = "{{ azure.AlertGroup.WebhookURL }}"
    use_common_alert_schema = true
  }
  {% endif %}
}
{% endif %}

{% include "./frontdoor.tf" %}
{% include "./url_locals.tf" %}

{% include "./app_service_plans.tf" %}

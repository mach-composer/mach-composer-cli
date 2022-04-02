provider "azurerm" {
  subscription_id = {{ azure.SubscriptionID|tf }}
  tenant_id       = {{ azure.TenantID|tf }}
  skip_provider_registration = true
  features {}
}


locals {
  tenant_id                    = {{ azure.TenantID|tf }}
  region                       = {{ azure.Region|tf }}
  subscription_id              = {{ azure.SubscriptionID|tf }}
  project_key                  = {{ site.Commercetools.ProjectKey|tf }}

  region_short                 = "{{ azure.Region|azure_region_short }}"
  name_prefix                  = format("{{ global.Azure.ResourcesPrefix }}{{ site.Identifier|replace:"dev,d"|replace:"tst,t"|replace:"prd,p" }}-%s", local.region_short)

  service_object_ids           = {
      {% for key, value in azure.ServiceObjectIds %}
          {{ key }} = {{ value|tf }}
      {% endfor %}
  }

  tags = {
    Site = "{{ site.Identifier }}"
    Environment = {{ global.Environment|tf }}
  }
}

{% if azure.ResourceGroup %}
data "azurerm_resource_group" "main" {
  name = {{ azure.ResourceGroup|tf }}
}
{% else %}
resource "azurerm_resource_group" "main" {
  name     = format("%s-rg", local.name_prefix)
  location = "{{ azure.Region|azure_region_long }}"
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
  name                = {{ azure.AlertGroup.LogicAppName()|tf }}
  resource_group_name = {{ azure.AlertGroup.LogicAppResourceGroup()|tf }}
}
{% endif %}

resource "azurerm_monitor_action_group" "alert_action_group" {
  name                = "{{ site.Identifier }}-{{ azure.AlertGroup.Name }}"
  resource_group_name = azurerm_resource_group.main.name
  short_name          = "{{ azure.AlertGroup.Name|replace:" ,"|replace:"-,"|lower }}"

  {% for email in azure.AlertGroup.AlertEmails %}
  email_receiver {
    name          = {{ email|tf }}
    email_address = {{ email|tf }}
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
    service_uri             = {{ azure.AlertGroup.WebhookURL|tf }}
    use_common_alert_schema = true
  }
  {% endif %}
}
{% endif %}

{% include "./frontdoor.tf" %}
{% include "./url_locals.tf" %}

{% include "./app_service_plans.tf" %}

{% if site.azure.front_door  %}
data "azurerm_dns_zone" "domain" {
    name                = "{{ site.azure.front_door.dns_zone }}"
    resource_group_name = "{{ site.azure.front_door.resource_group }}"
}

data "azurerm_key_vault" "ssl" {
  name                = "{{ site.azure.front_door.ssl_key_vault_name }}"
  resource_group_name = "{{ site.azure.front_door.resource_group }}"
}

locals {
    front_door_external_domain = "{{ site.commercetools.project_key }}.{{ site.azure.front_door.dns_zone }}"
    front_door_external_domain_identifier = replace(local.front_door_external_domain, ".", "-")
}

resource "azurerm_dns_cname_record" "{{ site.commercetools.project_key }}" {
  name                = "{{ site.commercetools.project_key }}"
  zone_name           = data.azurerm_dns_zone.domain.name
  resource_group_name = "{{ site.azure.front_door.resource_group }}"
  ttl                 = 600
  record              = local.front_door_domain
}
{% endif %}

{% if site.public_api_components %}
resource "azurerm_frontdoor" "app-service" {
  name                                          = format("%s-fd", local.name_prefix)
  resource_group_name                           = local.resource_group_name
  enforce_backend_pools_certificate_name_check  = false
  tags = local.tags

  backend_pool_load_balancing {
    name = "lbSettings"
  }

  frontend_endpoint {
    name                              = local.front_door_domain_identifier
    host_name                         = local.front_door_domain
    custom_https_provisioning_enabled = false
  }

  {% if site.azure.front_door %}
  frontend_endpoint {
    name                              = local.front_door_external_domain_identifier
    host_name                         = local.front_door_external_domain
    custom_https_provisioning_enabled = true

    custom_https_configuration {
      certificate_source                         = "AzureKeyVault"
      azure_key_vault_certificate_vault_id       = data.azurerm_key_vault.ssl.id
      # no data source for certificates yet.
      azure_key_vault_certificate_secret_name    = "{{ site.azure.front_door.ssl_key_vault_secret_name }}"
      azure_key_vault_certificate_secret_version = "{{ site.azure.front_door.ssl_key_vault_secret_version }}"
    }
  }
  
  depends_on = [azurerm_dns_cname_record.{{ site.commercetools.project_key }}]
  {% endif %}


  {% for component in site.public_api_components %}
  backend_pool_health_probe {
    name = "{{ component.name }}-hpSettings"
    path = "{% if component.health_check_path %}{{ component.health_check_path }}{% else %}/{{ component.name }}/healthchecks{% endif %}"
    protocol = "Https"
    enabled = false
    probe_method = "HEAD"
  }

  routing_rule {
    name               = "http-https-redirect"
    accepted_protocols = ["Http"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.front_door_domain_identifier{% if site.azure.front_door %}, local.front_door_external_domain_identifier{% endif %}]
    redirect_configuration {
      redirect_type     = "PermanentRedirect"
      redirect_protocol = "HttpsOnly"
    }
  }

  routing_rule {
    name               = "{{ component.name }}-routing-rule"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/{{ component.name }}/*"]
    frontend_endpoints = [local.front_door_domain_identifier{% if site.azure.front_door %}, local.front_door_external_domain_identifier{% endif %}]
    forwarding_configuration {
        forwarding_protocol = "MatchRequest"
        backend_pool_name   = "{{ component.name }}"
    }
  }

  backend_pool {
    name = "{{ component.name }}"
    backend {
        host_header = "${local.name_prefix}-func-{{ component.short_name }}.azurewebsites.net"
        address     = "${local.name_prefix}-func-{{ component.short_name }}.azurewebsites.net"
        http_port   = 80
        https_port  = 443
    }

    load_balancing_name = "lbSettings"
    health_probe_name   = "{{ component.name }}-hpSettings"
  }
  {% endfor %}
}
{% endif %}
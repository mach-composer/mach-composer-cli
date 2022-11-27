locals {
  frontdoor_domain            = format("%s-fd.azurefd.net", local.name_prefix)
  frontdoor_domain_identifier = replace(local.frontdoor_domain, ".", "-")
}

{% for endpoint in endpoints %}
data "azurerm_dns_zone" "{{ endpoint.Key }}" {
    name                = "{{ endpoint.Zone }}"
    resource_group_name = "{{ site.Azure.Frontdoor.DNSResourceGroup }}"
}

{% if endpoint.IsRootDomain() %}
resource "azurerm_dns_a_record" "{{ endpoint.Key }}" {
  name                = "@"
  zone_name           = data.azurerm_dns_zone.{{ endpoint.Key }}.name
  resource_group_name = "{{ site.Azure.Frontdoor.DNSResourceGroup }}"
  ttl                 = 600
  target_resource_id  = azurerm_frontdoor.app-service.id
}
{% else %}
resource "azurerm_dns_cname_record" "{{ endpoint.Key }}" {
  name                = "{{ endpoint.Subdomain() }}"
  zone_name           = data.azurerm_dns_zone.{{ endpoint.Key }}.name
  resource_group_name = "{{ azure.Frontdoor.DNSResourceGroup }}"
  ttl                 = 600
  record              = local.frontdoor_domain
}
{% endif %}
{% endfor %}

{% if site.UsedEndpoints() %}
locals {
  {% for endpoint in site.used_custom_endpoints %}
  {% for component in endpoint.components %}
  {% set cep_key = component|component_endpoint_name:endpoint %}
  fd_{{ endpoint.key }}_{{ component.name }}_route_defs = lookup(
    module.{{ component.name }}.azure_endpoint_{{ cep_key }},
    "routes",
    [{
      patterns = ["/{{ component.name }}/*"]
    }]
  )

  fd_{{ endpoint.key }}_{{ component.name }}_routes = {
    for i in range(
      length(
        local.fd_{{ endpoint.key }}_{{ component.name }}_route_defs
      )
    ) :
    i => element(
      local.fd_{{ endpoint.key }}_{{ component.name }}_route_defs,
      i
    )
  }
  {% endfor %}
  {% endfor %}
}

resource "azurerm_frontdoor" "app-service" {
  name                                          = format("%s-fd", local.name_prefix)
  resource_group_name                           = local.resource_group_name
  enforce_backend_pools_certificate_name_check  = false
  tags = local.tags

  backend_pool_load_balancing {
    name = "lbSettings"
  }

  frontend_endpoint {
    name                              = local.frontdoor_domain_identifier
    host_name                         = local.frontdoor_domain
  }

  {% for endpoint in site.used_custom_endpoints %}
  frontend_endpoint {
    name                              = {{ endpoint|azure_frontend_endpoint_name }}
    host_name                         = "{{ endpoint.url }}"
    {% if endpoint.azure.waf_policy_id %}
    web_application_firewall_policy_link_id = endpoint.azure.waf_policy_id
    {% endif %}

    {% if endpoint.azure.session_affinity_enabled %}
    session_affinity_enabled = endpoint.azure.session_affinity_enabled
    session_affinity_ttl_seconds = endpoint.azure.session_affinity_ttl_seconds
    {% endif %}

  }
  {% endfor %}

  depends_on = [
    {% for endpoint in site.used_custom_endpoints %}
    {% if not endpoint.IsRootDomain() %}
    azurerm_dns_cname_record.{{ endpoint.key }},
    {% endif %}
    {% endfor %}
  ]

  routing_rule {
    name               = "http-https-redirect"
    accepted_protocols = ["Http"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [
      local.frontdoor_domain_identifier,
      {% for endpoint in site.used_custom_endpoints %}
      {{ endpoint|azure_frontend_endpoint_name }},
      {% endfor %}
    ]
    redirect_configuration {
      redirect_type     = "PermanentRedirect"
      redirect_protocol = "HttpsOnly"
    }
  }

  {% for endpoint in site.UsedEndpoints() %}
  {% include "./frontdoor_endpoint.tf" %}
  {% endfor %}

  {% if site.azure.frontdoor.suppress_changes %}
  # Work-around for a very annoying bug in the Azure Frontdoor API
  # causing unwanted changes in Frontdoor and raising errors.
  lifecycle {
    ignore_changes = [
      routing_rule,
      backend_pool,
      backend_pool_health_probe,
      frontend_endpoint,
    ]
  }
  {% endif %}
}

{% if site.Azure.Frontdoor.SSLKeyVault %}
data "azurerm_key_vault" "ssl" {
  name                = "{{ site.Azure.Frontdoor.SSLKeyVault.Name }}"
  resource_group_name = "{{ site.Azure.Frontdoor.SSLKeyVault.ResourceGroup }}"
}
{% endif %}

{% for endpoint in endpoints %}
  {% if endpoint.URL %}
resource "azurerm_frontdoor_custom_https_configuration" "{{ endpoint.Key|slugify }}" {
  frontend_endpoint_id              = azurerm_frontdoor.app-service.frontend_endpoints[{{ endpoint|azure_frontend_endpoint_name }}]
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    {% if site.azure.frontdoor.ssl_key_vault %}
    certificate_source                         = "AzureKeyVault"
    azure_key_vault_certificate_vault_id       = data.azurerm_key_vault.ssl.id
    azure_key_vault_certificate_secret_name    = "{{ site.Azure.Frontdoor.SSLKeyVault.SecretName }}"
    {% else %}
    certificate_source                      = "FrontDoor"
    {% endif %}
  }
}
  {% endif %}
{% endfor %}
{% endif %}

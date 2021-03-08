locals {
  frontdoor_domain            = format("%s-fd.azurefd.net", local.name_prefix)
  frontdoor_domain_identifier = replace(local.frontdoor_domain, ".", "-")
}

{% for endpoint in site.used_custom_endpoints %}
data "azurerm_dns_zone" "{{ endpoint.key }}" {
    name                = "{{ endpoint.zone }}"
    resource_group_name = "{{ site.azure.frontdoor.resource_group }}"
}

resource "azurerm_dns_cname_record" "{{ endpoint.key }}" {
  name                = "{{ endpoint.subdomain }}"
  zone_name           = data.azurerm_dns_zone.{{ endpoint.key }}.name
  resource_group_name = "{{ site.azure.frontdoor.resource_group }}"
  ttl                 = 600
  record              = local.frontdoor_domain
}
{% endfor %}

{% if site.used_endpoints %}
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
    name                              = "{{ endpoint.key }}"
    host_name                         = "{{ endpoint.url }}"

    custom_https_provisioning_enabled = true
    custom_https_configuration {
      certificate_source = "FrontDoor"
    }
  }
  {% endfor %}

  depends_on = [
    {% for endpoint in site.used_custom_endpoints %}
    azurerm_dns_cname_record.{{ endpoint.key }},
    {% endfor %}
  ]

  routing_rule {
    name               = "http-https-redirect"
    accepted_protocols = ["Http"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [
      local.frontdoor_domain_identifier,
      {% for endpoint in site.used_custom_endpoints %}
      "{{ endpoint.key }}",
      {% endfor %}
    ]
    redirect_configuration {
      redirect_type     = "PermanentRedirect"
      redirect_protocol = "HttpsOnly"
    }
  }

  {% for endpoint in site.used_endpoints %}
  {% for component in endpoint.components %}
  backend_pool_health_probe {
    name = "{{ endpoint.key }}-{{ component.name }}-hpSettings"
    path = "{% if component.health_check_path %}{{ component.health_check_path }}{% else %}/{{ component.name }}/healthchecks{% endif %}"
    protocol = "Https"
    enabled = false
    probe_method = "HEAD"
  }

  routing_rule {
    name               = "{{ endpoint.key }}-{{ component.name }}-routing"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/{{ component.name }}/*"]
    frontend_endpoints = [
      local.frontdoor_domain_identifier,
      {% if endpoint.url %}
      "{{ endpoint.key }}",
      {% endif %}
    ]
    forwarding_configuration {
        forwarding_protocol = "MatchRequest"
        backend_pool_name   = "{{ endpoint.key }}-{{ component.name }}"
    }
  }

  backend_pool {
    name = "{{ endpoint.key }}-{{ component.name }}"
    backend {
        host_header = "${local.name_prefix}-func-{{ component.azure.short_name }}.azurewebsites.net"
        address     = "${local.name_prefix}-func-{{ component.azure.short_name }}.azurewebsites.net"
        http_port   = 80
        https_port  = 443
    }

    load_balancing_name = "lbSettings"
    health_probe_name   = "{{ endpoint.key }}-{{ component.name }}-hpSettings"
  }
  {% endfor %}
  {% endfor %}

  {% if site.azure.frontdoor.suppress_changes %}
  # Work-around for a very annoying bug in the Azure Frontdoor API
  # causing unwanted changes in Frontdoor and raising errors.
  lifecycle {
    ignore_changes = [
      routing_rule,
      backend_pool,
      backend_pool_health_probe,
    ]
  }
  {% endif %}
}
{% endif %}

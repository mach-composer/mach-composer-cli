locals {
  frontdoor_domain            = format("%s-fd.azurefd.net", local.name_prefix)
  frontdoor_domain_identifier = replace(local.frontdoor_domain, ".", "-")
}

{% if site.azure.frontdoor  %}
data "azurerm_dns_zone" "domain" {
    name                = "{{ site.azure.frontdoor.dns_zone }}"
    resource_group_name = "{{ site.azure.frontdoor.resource_group }}"
}

locals {
    frontdoor_external_domain = "{{ site.commercetools.project_key }}.{{ site.azure.frontdoor.dns_zone }}"
    frontdoor_external_domain_identifier = replace(local.frontdoor_external_domain, ".", "-")
}

resource "azurerm_dns_cname_record" "{{ site.commercetools.project_key }}" {
  name                = "{{ site.commercetools.project_key }}"
  zone_name           = data.azurerm_dns_zone.domain.name
  resource_group_name = "{{ site.azure.frontdoor.resource_group }}"
  ttl                 = 600
  record              = local.frontdoor_domain
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
    name                              = local.frontdoor_domain_identifier
    host_name                         = local.frontdoor_domain
  }

  {% if site.azure.frontdoor %}
  frontend_endpoint {
    name                              = local.frontdoor_external_domain_identifier
    host_name                         = local.frontdoor_external_domain

    custom_https_configuration {
      certificate_source = "FrontDoor"
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
    frontend_endpoints = [local.frontdoor_domain_identifier{% if site.azure.frontdoor %}, local.frontdoor_external_domain_identifier{% endif %}]
    redirect_configuration {
      redirect_type     = "PermanentRedirect"
      redirect_protocol = "HttpsOnly"
    }
  }

  routing_rule {
    name               = "{{ component.name }}-routing-rule"
    accepted_protocols = ["Https"]
    patterns_to_match  = ["/{{ component.name }}/*"]
    frontend_endpoints = [local.frontdoor_domain_identifier{% if site.azure.frontdoor %}, local.frontdoor_external_domain_identifier{% endif %}]
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
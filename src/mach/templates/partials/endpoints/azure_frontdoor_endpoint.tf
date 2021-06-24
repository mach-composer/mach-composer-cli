{% for component in endpoint.components %}
backend_pool_health_probe {
  name = "{{ endpoint.key }}-{{ component.name }}-hpSettings"
  path = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "health_probe_path", "/")
  protocol = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "health_probe_protocol", "Https")
  enabled = contains(keys(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}), "health_probe_path")
  probe_method = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "health_probe_method", "GET")
}

routing_rule {
  name               = "{{ endpoint.key }}-{{ component.name }}-routing"
  accepted_protocols = ["Https"]
  patterns_to_match  = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "routing_patterns", ["/{{ component.name }}/*"])
  frontend_endpoints = [
    local.frontdoor_domain_identifier,
    {% if endpoint.url %}
    {{ endpoint.key|tf }},
    {% endif %}
  ]
  forwarding_configuration {
      forwarding_protocol            = "MatchRequest"
      backend_pool_name              = "{{ endpoint.key }}-{{ component.name }}"
      cache_enabled                  = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "cache_enabled", false)
      custom_forwarding_path         = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "custom_forwarding_path", null)
  }
}

backend_pool {
  name = "{{ endpoint.key }}-{{ component.name }}"
  backend {
      host_header = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "host_header", module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}.address)
      address     = module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}.address
      http_port   = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "http_port", 80)
      https_port  = lookup(module.{{ component.name }}.azure_endpoint_{{ endpoint.key }}, "https_port", 443)
  }

  load_balancing_name = "lbSettings"
  health_probe_name   = "{{ endpoint.key }}-{{ component.name }}-hpSettings"
}
{% endfor %}

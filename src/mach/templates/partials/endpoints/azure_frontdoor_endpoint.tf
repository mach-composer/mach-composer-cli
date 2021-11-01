{% for component in endpoint.components %}

{# set the component endpoint key #}
{% set cep_key = component|component_endpoint_name(endpoint) %}
backend_pool_health_probe {
  name = "{{ endpoint.key }}-{{ component.name }}-hpSettings"
  path = lookup(module.{{ component.name }}.azure_endpoint_{{ cep_key }}, "health_probe_path", "/")
  protocol = lookup(module.{{ component.name }}.azure_endpoint_{{ cep_key }}, "health_probe_protocol", "Https")
  enabled = contains(keys(module.{{ component.name }}.azure_endpoint_{{ cep_key }}), "health_probe_path")
  probe_method = lookup(module.{{ component.name }}.azure_endpoint_{{ cep_key }}, "health_probe_method", "GET")
}

dynamic "routing_rule" {
  for_each = local.fd_{{ endpoint.key }}_{{ component.name }}_routes
  content {
    name = "{{ endpoint.key }}-{{ component.name }}-routing-${lookup(routing_rule.value, "name", routing_rule.key)}"

    accepted_protocols = ["Https"]
    patterns_to_match  = routing_rule.value.patterns
    frontend_endpoints = [
      {% if endpoint.url %}
      {{ endpoint|azure_frontend_endpoint_name }},
      {% else %}
      local.frontdoor_domain_identifier,
      {% endif %}
    ]
    forwarding_configuration {
        forwarding_protocol            = "MatchRequest"
        backend_pool_name              = "{{ endpoint.key }}-{{ component.name }}"
        cache_enabled                  = lookup(routing_rule.value, "cache_enabled", false)
        custom_forwarding_path         = lookup(routing_rule.value, "custom_forwarding_path", null)
    }
  }
}

backend_pool {
  name = "{{ endpoint.key }}-{{ component.name }}"
  backend {
      host_header = lookup(module.{{ component.name }}.azure_endpoint_{{ cep_key }}, "host_header", module.{{ component.name }}.azure_endpoint_{{ cep_key }}.address)
      address     = module.{{ component.name }}.azure_endpoint_{{ cep_key }}.address
      http_port   = lookup(module.{{ component.name }}.azure_endpoint_{{ cep_key }}, "http_port", 80)
      https_port  = lookup(module.{{ component.name }}.azure_endpoint_{{ cep_key }}, "https_port", 443)
  }

  load_balancing_name = "lbSettings"
  health_probe_name   = "{{ endpoint.key }}-{{ component.name }}-hpSettings"
}
{% endfor %}

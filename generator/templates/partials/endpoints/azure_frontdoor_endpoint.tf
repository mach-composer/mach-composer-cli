{% for component in endpoint.Components %}

{# set the component endpoint key #}
{% set cep_key = component|component_endpoint_name:endpoint %}
backend_pool_health_probe {
  name = "{{ endpoint.Key }}-{{ component.Name }}-hpSettings"
  path = lookup(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}, "health_probe_path", "/")
  protocol = lookup(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}, "health_probe_protocol", "Https")
  enabled = contains(keys(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}), "health_probe_path")
  probe_method = lookup(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}, "health_probe_method", "GET")
}

dynamic "routing_rule" {
  for_each = local.fd_{{ endpoint.key }}_{{ component.Name }}_routes
  content {
    name = "{{ endpoint.Key }}-{{ component.Name }}-routing-${lookup(routing_rule.value, "name", routing_rule.key)}"

    accepted_protocols = ["Https"]
    patterns_to_match  = routing_rule.value.patterns
    frontend_endpoints = [
      {% if endpoint.URL %}
      {{ endpoint|azure_frontend_endpoint_name }},
      {% else %}
      local.frontdoor_domain_identifier,
      {% endif %}
    ]
    forwarding_configuration {
        forwarding_protocol            = "MatchRequest"
        backend_pool_name              = "{{ endpoint.Key }}-{{ component.Name }}"
        cache_enabled                  = lookup(routing_rule.value, "cache_enabled", false)
        custom_forwarding_path         = lookup(routing_rule.value, "custom_forwarding_path", null)
    }
  }
}

backend_pool {
  name = "{{ endpoint.Key }}-{{ component.Name }}"
  backend {
      host_header = lookup(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}, "host_header", module.{{ component.Name }}.azure_endpoint_{{ cep_key }}.address)
      address     = module.{{ component.Name }}.azure_endpoint_{{ cep_key }}.address
      http_port   = lookup(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}, "http_port", 80)
      https_port  = lookup(module.{{ component.Name }}.azure_endpoint_{{ cep_key }}, "https_port", 443)
  }

  load_balancing_name = "lbSettings"
  health_probe_name   = "{{ endpoint.Key }}-{{ component.Name }}-hpSettings"
}
{% endfor %}

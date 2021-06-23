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
    {{ endpoint.key|tf }},
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

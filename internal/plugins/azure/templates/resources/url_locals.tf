locals {
{% for endpoint in site.UsedEndpoints() %}
endpoint_url_{{ endpoint.Key|slugify }} = {% if endpoint.URL %}"{{ endpoint.URL }}"{% else %}local.frontdoor_domain{% endif %}
{% endfor %}
}

output "endpoints" {
  value = {
  {% for endpoint in site.UsedEndpoints() -%}
    {{ endpoint.Key|slugify }}: local.endpoint_url_{{ endpoint.Key|slugify }}
  {% endfor %}
  }
}

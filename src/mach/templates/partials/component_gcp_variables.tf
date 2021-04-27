{% for component_endpoint, site_endpoint in component.endpoints.items() -%}
gcp_endpoint_{{ component_endpoint|slugify }} = {
  url = local.endpoint_url_{{ site_endpoint|slugify }}
  api_id = google_api_gateway_api.{{ site_endpoint|slugify }}_gateway.api_id
}
{% endfor %}

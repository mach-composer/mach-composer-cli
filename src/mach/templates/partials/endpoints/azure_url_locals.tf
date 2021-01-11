locals {
{% for endpoint in site.used_endpoints %}
endpoint_url_{{ endpoint.key|slugify }} = "{{ endpoint.url }}"
{% endfor %}
}
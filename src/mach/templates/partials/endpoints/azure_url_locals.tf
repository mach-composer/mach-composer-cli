locals {
{% for endpoint in site.used_endpoints %}
endpoint_url_{{ endpoint.key|slugify }} = {% if endpoint.url %}"{{ endpoint.url }}"{% else %}local.frontdoor_domain{% endif %}
{% endfor %}

}
{% set gcp = site.gcp %}
provider "google" {
  project = "{{ gcp.project_id }}"
  region  = "{{ gcp.region }}"
}

provider "google-beta" {
  project = "{{ gcp.project_id }}"
  region  = "{{ gcp.region }}"
}

{% for provider in gcp.extra_providers %}
provider "google" {
  project = "{{ gcp.project_id }}"
  region  = "{{ gcp.region }}"
}

provider "google-beta" {
  project = "{{ gcp.project_id }}"
  region  = "{{ gcp.region }}"
}
{% endfor %}

{% if site.used_endpoints %}
  {% for zone in site.dns_zones %}
  data "google_dns_managed_zone" "{{ zone|slugify }}" {
    name = "{{ zone|slugify }}"
  }
  {% endfor %}

  {% for endpoint in site.used_endpoints %}
    {% include 'partials/endpoints/gcp_api_gateway.tf' %}
  {% endfor %}
  {% include 'partials/endpoints/gcp_url_locals.tf' %}
{% endif %}

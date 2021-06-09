{% set gcp = site.gcp %}
provider "google" {
  project = "{{ gcp.project_id }}"
  region  = "{{ gcp.region }}"
}

provider "google-beta" {
  project = "{{ gcp.project_id }}"
  region  = "{{ gcp.region }}"
}

{% if site.used_endpoints %}
  {% for zone in site.dns_zones %}
  data "google_dns_managed_zone" "{{ zone|replace('.', '-')|slugify }}" {
    name = "{{ zone|replace('.', '-')|slugify('-') }}"
  }
  {% endfor %}

  {% for endpoint in site.used_endpoints %}
    {% include 'partials/endpoints/gcp_api_gateway.tf' %}
    
  {% endfor %}

  {% include 'partials/endpoints/gcp_url_locals.tf' %}
{% endif %}

locals {
  tags = {
    Site = "{{ site.identifier }}"
    Environment = "{{ general_config.environment }}"
  }
}

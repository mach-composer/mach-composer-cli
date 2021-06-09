

resource "google_api_gateway_api" "{{ endpoint.key|slugify }}_gateway" {
  provider = google-beta
  api_id = "{{ endpoint.key|slugify }}"
}

locals {
  api_spec_base_{{ endpoint.key|slugify }} = <<EOT
swagger: '2.0'
info:
  title: {{ endpoint.key }} API
  description: {{ endpoint.key }} API
  version: 0.1.0
schemes:
  - https
produces:
  - application/json
EOT
}


data "utils_deep_merge_yaml" "{{ endpoint.key|slugify }}_api_spec" {
  input = [
    {% for component in endpoint.components %}
    module.{{ component.name }}.gcp_api_spec_{{ endpoint.key|slugify }},
    {% endfor %}
    local.api_spec_base_{{ endpoint.key|slugify }},
  ]
}

locals {
  api_spec_{{ endpoint.key|slugify }} = base64encode(data.utils_deep_merge_yaml.{{ endpoint.key|slugify }}_api_spec.output)
  api_spec_{{ endpoint.key|slugify }}_hash = lower(substr(md5(local.api_spec_{{ endpoint.key|slugify }}), 0, 5))
}

resource "google_api_gateway_api_config" "{{ endpoint.key|slugify }}" {
  provider = google-beta
  api = google_api_gateway_api.{{ endpoint.key|slugify }}_gateway.api_id
  api_config_id = "{{ endpoint.key|slugify}}-cfg-${local.api_spec_{{ endpoint.key|slugify }}_hash}"

  openapi_documents {
    document {
      path = "spec.yaml"
      contents = local.api_spec_{{ endpoint.key|slugify }}
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "{{ endpoint.key|slugify }}" {
  provider = google-beta
  api_config = google_api_gateway_api_config.{{ endpoint.key|slugify }}.id
  gateway_id = "{{ endpoint.key|slugify }}-gw-${local.api_spec_{{ endpoint.key|slugify }}_hash}"

  # TODO: Hard coded for now; we might want to make this configurable in the endpoint settings
  # so we can overwrite for only the API gateway (since for example europe-west3 is not supported 
  # yet)
  region = "europe-west2"

  lifecycle {
    create_before_destroy = true
  }
}

# TODO: This doesn't work yet; it needs to have a load balancer in front of it
resource "google_dns_record_set" "{{ endpoint.key|slugify }}_api" {
  provider     = "google-beta"
  managed_zone = data.google_dns_managed_zone.{{ endpoint.zone|replace('.', '-')|slugify }}.name
  name         = "{{ endpoint.url }}."
  type         = "CNAME"
  rrdatas      = ["${google_api_gateway_gateway.{{ endpoint.key|slugify }}.default_hostname}."]
  ttl          = 600
}

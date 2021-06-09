

resource "google_api_gateway_api" "{{ endpoint.key|slugify }}_gateway" {
  provider = google-beta
  api_id = "{{ endpoint.key|slugify }}"
}


data "utils_deep_merge_yaml" "{{ endpoint.key|slugify }}_api_spec" {
  input = [
    {% for component in endpoint.components %}
    module.{{ component.name }}.gcp_api_spec_{{ endpoint.key|slugify }},
    {% endfor %}
  ]
}

resource "google_api_gateway_api_config" "{{ endpoint.key|slugify }}" {
  provider = google-beta
  api = google_api_gateway_api.{{ endpoint.key|slugify }}_gateway.api_id
  api_config_id = "{{ endpoint.key|slugify}}-cfg"

  openapi_documents {
    document {
      path = "spec.yaml"
      contents = base64encode(data.utils_deep_merge_yaml.{{ endpoint.key|slugify }}_api_spec.output)
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "{{ endpoint.key|slugify }}" {
  provider = google-beta
  api_config = google_api_gateway_api_config.{{ endpoint.key|slugify }}.id
  gateway_id = "{{ endpoint.key|slugify }}-gw"

  # TODO: Hard coded for now; we might want to make this configurable in the endpoint settings
  # so we can overwrite for only the API gateway (since for example europe-west3 is not supported 
  # yet)
  region = "europe-west2"
}


output "deep_merge_output" {
  value = data.utils_deep_merge_yaml.{{ endpoint.key|slugify }}_api_spec.output
}


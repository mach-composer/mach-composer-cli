

resource "google_api_gateway_api" "{{ endpoint.key|slugify }}_gateway" {
  provider = google-beta
  api_id = "{{ endpoint.key|slugify }}"
}


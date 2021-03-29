{% if commercetools.stores %}
{% for store in commercetools.stores %}
resource "commercetools_api_client" "frontend_credentials_{{ store.key }}" {
  name = "frontend_credentials_terraform_{{ store.key }}"
  scope = {{ commercetools.frontend.permission_scopes|render_commercetools_scopes(commercetools.project_key, store.key) }}

  {% if store.managed %}
  depends_on = [commercetools_store.{{ store.key }}]
  {% endif %}
}

output "frontend_client_scope_{{ store.key }}" {
  value = commercetools_api_client.frontend_credentials_{{ store.key }}.scope
}

output "frontend_client_id_{{ store.key }}" {
  value = commercetools_api_client.frontend_credentials_{{ store.key }}.id
}

output "frontend_client_secret_{{ store.key }}" {
  value = commercetools_api_client.frontend_credentials_{{ store.key }}.secret
}
{% endfor %}
{% else %}
resource "commercetools_api_client" "frontend_credentials" {
  name = "frontend_credentials_terraform"
  scope = {{ commercetools.frontend.permission_scopes|render_commercetools_scopes(commercetools.project_key) }}
}

output "frontend_client_scope" {
    value = commercetools_api_client.frontend_credentials.scope
}

output "frontend_client_id" {
    value = commercetools_api_client.frontend_credentials.id
}

output "frontend_client_secret" {
    value = commercetools_api_client.frontend_credentials.secret
}
{% endif %}

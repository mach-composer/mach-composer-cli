{% if commercetools.Stores %}
{% for store in commercetools.Stores %}
resource "commercetools_api_client" "frontend_credentials_{{ store.Key }}" {
  name = "frontend_credentials_terraform_{{ store.Key }}"
  {# scope = {{ commercetools.Frontend.PermissionScopes|render_commercetools_scopes(commercetools.ProjectKey, store.Key) }} #}
  scope = {{ commercetools.Frontend.PermissionScopes|render_commercetools_scopes:commercetools.ProjectKey }}

  {% if store.Managed %}
  depends_on = [commercetools_store.{{ store.Key }}]
  {% endif %}
}

output "frontend_client_scope_{{ store.Key }}" {
  value = commercetools_api_client.frontend_credentials_{{ store.Key }}.scope
}

output "frontend_client_id_{{ store.Key }}" {
  value = commercetools_api_client.frontend_credentials_{{ store.Key }}.id
}

output "frontend_client_secret_{{ store.Key }}" {
  value = commercetools_api_client.frontend_credentials_{{ store.Key }}.secret
}
{% endfor %}
{% else %}
resource "commercetools_api_client" "frontend_credentials" {
  name = "frontend_credentials_terraform"
  scope = {{ commercetools.Frontend.PermissionScopes|render_commercetools_scopes:commercetools.ProjectKey }}
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

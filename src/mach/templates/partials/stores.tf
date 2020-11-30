{% if commercetools.stores %}
{% set stores = commercetools.stores %}

{% for store in stores %}
resource "commercetools_store" "{{ store.key }}" {
  key  = "{{ store.key }}"
  name = {
    {% for language, localized_name in store.name.items() %}
        {{ language  }} = "{{ localized_name }}"
    {% endfor %}
  }
  {% if store.languages %}
  languages  = [{% for language in store.languages %}"{{ language }}"{% if not loop.last %},{% endif %}{% endfor %}]
  {% endif %}

  {% if store.distribution_channels %}
  distribution_channels = [{% for dc in store.distribution_channels %}commercetools_channel.{{ dc }}.key{% if not loop.last %},{% endif %}{% endfor %}]
  {% endif %}
  {% if store.supply_channels %}
  supply_channels = [{% for sc in store.supply_channels %}commercetools_channel.{{ sc }}.key{% if not loop.last %},{% endif %}{% endfor %}]
  {% endif %}
}

{% if commercetools.create_frontend_credentials %}
resource "commercetools_api_client" "frontend_credentials_{{ store.key }}" {
  name = "frontend_credentials_terraform_{{ store.key }}"
  scope = [
    "manage_my_profile:{{ commercetools.project_key }}",
    "manage_my_orders:{{ commercetools.project_key }}",
    "manage_my_orders:{{ commercetools.project_key}}:{{ store.key }}",
    "view_states:{{ commercetools.project_key }}",
    "manage_my_shopping_lists:{{ commercetools.project_key }}",
    "view_products:{{ commercetools.project_key }}",
    "manage_my_payments:{{ commercetools.project_key}}",
    "create_anonymous_token:{{ commercetools.project_key }}",
    "view_project_settings:{{ commercetools.project_key }}"
  ]

  depends_on = [commercetools_store.{{ store.key }}]
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
{% endif %}
{% endfor %}

{% elif commercetools.create_frontend_credentials %}
{# note: No stores definied, create 1 set of credentials #}

resource "commercetools_api_client" "frontend_credentials" {
  name = "frontend_credentials_terraform"
  scope = [
    "manage_my_profile:{{ commercetools.project_key }}",
    "manage_my_orders:{{ commercetools.project_key }}",
    "view_states:{{ commercetools.project_key }}",
    "manage_my_shopping_lists:{{ commercetools.project_key }}",
    "view_products:{{ commercetools.project_key }}",
    "manage_my_payments:{{ commercetools.project_key}}",
    "create_anonymous_token:{{ commercetools.project_key }}",
    "view_project_settings:{{ commercetools.project_key }}"
  ]
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
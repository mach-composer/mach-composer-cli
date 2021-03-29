{% set commercetools = site.commercetools %}

provider "commercetools" {
    client_id     = "{{ commercetools.client_id }}"
    client_secret = "{{ commercetools.client_secret }}"
    project_key   = "{{ commercetools.project_key }}"
    scopes        = "{{ commercetools.scopes }}"
    token_url     = "{{ commercetools.token_url }}"
    api_url       = "{{ commercetools.api_url }}"
}

{% if commercetools.project_settings %}
resource "commercetools_project_settings" "project" {
    name       = "{{ commercetools.project_key }}"
    countries  = [{% for country in commercetools.project_settings.countries %}"{{ country }}"{% if not loop.last %},{% endif %}{% endfor %}]
    currencies = [{% for currency in commercetools.project_settings.currencies %}"{{ currency }}"{% if not loop.last %},{% endif %}{% endfor %}]
    languages  = [{% for language in commercetools.project_settings.languages %}"{{ language }}"{% if not loop.last %},{% endif %}{% endfor %}]
    messages   = {
        enabled = {{ commercetools.project_settings.messages_enabled | string | lower }}
    }
}
{% endif %}

{% for channel in commercetools.channels %}
resource "commercetools_channel" "{{ channel.key }}" {
    key = "{{ channel.key }}"
    roles = [{% for role in channel.roles %}"{{ role }}"{% if not loop.last %}, {% endif %}{% endfor %}]

    {% if channel.name %}
    name = {
        {% for language, localized_name in channel.name.items() %}
        {{ language }} = "{{ localized_name }}"
        {% endfor %}
    }
    {% endif %}

    {% if channel.description %}
    description = {
        {% for language, localized_name in channel.description.items() %}
        {{ language }} = "{{ localized_name}}"
        {% endfor %}
    }
    {% endif %}
}
{% endfor %}

{% if commercetools.taxes %}
resource "commercetools_tax_category" "standard" {
  name = "Standard tax category"
  key  = "standard"
}

{% for tax in commercetools.taxes %}
resource "commercetools_tax_category_rate" "{{ tax.country|lower }}_vat" {
  tax_category_id = commercetools_tax_category.standard.id
  name = "{{ tax.name }}"
  amount = {{ tax.amount }}
  country = "{{ tax.country }}"
  included_in_price = true
}
{% endfor %}
{% endif %}

{% for zone in commercetools.zones %}
resource "commercetools_shipping_zone" "{{ zone.name|slugify }}" {
  name = "{{ zone.name }}"
  description = "{{ zone.description }}"
  {% for location in zone.locations %}
  location {
      country = "{{ location.country }}"
      {% if location.state %}
      state = "{{ location.state }}"
      {% endif %}
  }
  {% endfor %}
}
{% endfor %}

output "frontend_channels" {
    value = [
        {% for channel in commercetools.channels %}commercetools_channel.{{ channel.key }}.id,{% endfor %}
    ]
}

resource "null_resource" "commercetools" {
  depends_on = [
    {% if commercetools.project_settings %}
    commercetools_project_settings.project,
    {% endif %}
    {% for channel in commercetools.channels %}
    commercetools_channel.{{ channel.key }},
    {% endfor %}
    {% if commercetools.taxes %}
    commercetools_tax_category.standard,
    {% endif %}
    {% for store in stores %}
    commercetools_store.{{ store.key }},
    {% endfor %}
  ]
}

{% if commercetools.stores %}
{% include 'partials/stores.tf' %}
{% endif %}

{% if commercetools.frontend.create_credentials %}
{% include 'partials/commercetools_frontend.tf' %}
{% endif %}

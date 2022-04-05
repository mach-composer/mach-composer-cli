provider "commercetools" {
    client_id     = {{ commercetools.ClientID|tf }}
    client_secret = {{ commercetools.ClientSecret|tf }}
    project_key   = {{ commercetools.ProjectKey|tf }}
    scopes        = {{ commercetools.Scopes|tf }}
    token_url     = {{ commercetools.TokenURL|tf }}
    api_url       = {{ commercetools.ApiURL|tf }}
}

{% if commercetools.ProjectSettings %}
resource "commercetools_project_settings" "project" {
    name       = {{ commercetools.ProjectKey|tf }}
    countries  = [{% for country in commercetools.ProjectSettings.Countries %}{{ country|tf }}{% if not forloop.Last %},{% endif %}{% endfor %}]
    currencies = [{% for currency in commercetools.ProjectSettings.Currencies %}{{ currency|tf }}{% if not forloop.Last %},{% endif %}{% endfor %}]
    languages  = [{% for language in commercetools.ProjectSettings.Languages %}{{ language|tf }}{% if not forloop.Last %},{% endif %}{% endfor %}]
    messages   = {
        enabled = {{ commercetools.ProjectSettings.MessagesEnabled | string | lower }}
    }
}
{% endif %}

{%- for channel in commercetools.Channels %}
resource "commercetools_channel" "{{ channel.Key }}" {
    key = "{{ channel.Key }}"
    roles = {{ channel.Roles|tf }}

    {%- if channel.Name %}
    name = {
        {% for language, localized_name in channel.Name %}
        {{- language }} = {{ localized_name|tf }}
        {%- endfor %}
    }
    {%- endif %}

    {%- if channel.Description %}
    description = {
        {% for language, localized_name in channel.Description %}
        {{ language }} = {{ localized_name|tf }}
        {% endfor %}
    }
    {%- endif %}
}
{% endfor %}

{% for tax_category in commercetools.TaxCategories %}
resource "commercetools_tax_category" "{{ tax_category.Key|lower }}" {
  name = {{ tax_category.Name|tf }}
  key = {{ tax_category.Key|tf }}
}

  {% for rate in tax_category.Rates %}
resource "commercetools_tax_category_rate" "{{ rate.Name|slugify }}" {
  tax_category_id = commercetools_tax_category.{{ tax_category.Key|lower }}.id
  name = {{ rate.Name|tf }}
  amount = {{ rate.Amount|tf }}
  country = "{{ rate.Country }}"
  included_in_price = {{ rate.IncludedInPrice|tf }}
}
  {% endfor %}
{% endfor %}

{% if commercetools.Taxes %}
  resource "commercetools_tax_category" "standard" {
    name = "Standard tax category"
    key  = "standard"
  }

  {% for tax in commercetools.Taxes %}
  resource "commercetools_tax_category_rate" "{{ tax.Country|lower }}_vat" {
    tax_category_id = commercetools_tax_category.standard.id
    name = {{ tax.Name|tf }}
    amount = {{ tax.Amount|tf }}
    country = "{{ tax.Country }}"
    included_in_price = true
  }
  {% endfor %}
{% endif %}

{%- for zone in commercetools.Zones %}
resource "commercetools_shipping_zone" "{{ zone.Name|slugify }}" {
  name = "{{ zone.Name }}"
  description = {{ zone.Description|tf }}
  {% for location in zone.Locations %}
  location {
      country = {{ location.Country|tf }}
      {% if location.State %}
      state = {{ location.State|tf }}
      {% endif %}
  }
  {% endfor %}
}
{% endfor %}

output "frontend_channels" {
    value = [
        {% for channel in commercetools.Channels %}commercetools_channel.{{ channel.Key }}.id,{% endfor %}
    ]
}

resource "null_resource" "commercetools" {
  depends_on = [
    {%- if commercetools.ProjectSettings %}
    commercetools_project_settings.project,
    {%- endif -%}
    {%- for channel in commercetools.Channels %}
    commercetools_channel.{{ channel.Key }},
    {%- endfor -%}
    {% if commercetools.Taxes -%}
    commercetools_tax_category.standard,
    {%- endif %}
    {%- for store in commercetools.Stores %}
    commercetools_store.{{ store.Key }},
    {%- endfor -%}
  ]
}

{% if commercetools.Stores %}
  {% include "./stores.tf" %}
{% endif %}

{% if commercetools.Frontend.CreateCredentials %}
  {% include "./commercetools_frontend.tf" %}
{% endif %}

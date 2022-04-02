{% for store in commercetools.managed_stores %}
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
{% endfor %}

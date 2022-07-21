{% for store in commercetools.ManagedStores %}
resource "commercetools_store" "{{ store.Key }}" {
  key  = "{{ store.Key }}"
  name = {
    {% for language, localizedName in store.Name %}
        {{ language  }} = {{ localizedName|tf }}
    {% endfor %}
  }
  {% if store.Languages %}
  languages  = [{% for language in store.Languages %}"{{ language }}"{% if not forloop.Last %},{% endif %}{% endfor %}]
  {% endif %}

  {% if store.DistributionChannels %}
  distribution_channels = [{% for dc in store.DistributionChannels %}commercetools_channel.{{ dc }}.key{% if not forloop.Last %},{% endif %}{% endfor %}]
  {% endif %}
  {% if store.SupplyChannels %}
  supply_channels = [{% for sc in store.SupplyChannels %}commercetools_channel.{{ sc }}.key{% if not forloop.Last %},{% endif %}{% endfor %}]
  {% endif %}
}
{% endfor %}

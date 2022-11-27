# Module: {{ component.Name }}
{% for v in pluginResources %}
  {{ v|safe }}
{% endfor %}

module "{{ component.Name }}" {
  source            = "{{ definition.Source|safe }}{% if definition.UseVersionReference() %}?ref={{ definition.Version }}{% endif %}"

  {% if component.HasCloudIntegration() || component.Variables %}
  variables = {
    {% for key, value in component.Variables -%}
    {{ key }} = {{ value|interpolate }}
    {% endfor %}
  }
  {% endif %}

  {% if component.HasCloudIntegration() || component.Secrets -%}
  secrets = {
    {% for key, value in component.Secrets -%}
    {{ key }} = {{ value|interpolate }}
    {% endfor %}
  }
  {%- endif %}

  {% if component.HasCloudIntegration() -%}
  component_version       = "{{ definition.Version }}"
  environment             = "{{ global.Environment }}"
  site                    = "{{ site.Identifier }}"
  tags                    = local.tags
  {%- endif %}



  {% for v in pluginVariables %}
    {{ v|safe }}
  {% endfor %}

  providers = {
    {% for v in pluginProviders %}
      {{ v|safe }},
    {% endfor %}
  }

  depends_on = [
    {% for v in pluginDependsOn %}
      {{ v|safe }},
    {% endfor %}
  ]
}

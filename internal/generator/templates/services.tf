{% if variables.Encrypted %}
data "local_file" "variables" {
  filename = "{{ variables.Filepath }}"
}

data "sops_external" "variables" {
  source     = data.local_file.variables.content
  input_type = "yaml"
}
{% endif %}

# Plugins
{% for resource in resources %}
  {{ resource|safe }}
{% endfor %}

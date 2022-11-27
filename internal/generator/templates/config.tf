terraform {

  {{ backend_config|safe }}

  required_providers {
    {% for provider in providers %}
      {{ provider|safe }}
    {% endfor %}

    {%- if variables.Encrypted %}
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
    {%- endif %}
  }
}


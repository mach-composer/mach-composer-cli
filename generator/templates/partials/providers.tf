terraform {
  required_providers {
    {%- if site.AWS %}
    aws = {
      version = "~> {{ global.TerraformConfig.providers.aws|default:"3.66.0" }}"
    }
    {% endif -%}

    {%- if site.Azure %}
    azurerm = {
      version = "~> {{ global.TerraformConfig.providers.azure|default:"2.86.0" }}"
    }
    {% endif -%}

    {%- if site.Commercetools %}
    commercetools = {
      source = "labd/commercetools"
      version = "~> {{ global.TerraformConfig.Providers.Commercetools|default:'0.29.3' }}"
    }
    {% endif -%}

    {%- if site.contentful %}
    contentful = {
      source = "labd/contentful"
      version = "~> {{ global.TerraformConfig.Providers.Contentful|default:'0.1.0' }}"
    }
    {% endif -%}

    {%- if site.Amplience %}
    amplience = {
      source = "labd/amplience"
      version = "~> {{ global.TerraformConfig.Providers.Amplience|default:'0.2.2' }}"
    }
    {% endif -%}

    {%- if global.SentryConfig.AuthToken %}
    sentry = {
      source = "jianyuan/sentry"
      version = "~> {{ global.TerraformConfig.Providers.Sentry|default:'0.6.0' }}"
    }
    {% endif -%}

    {%- if config.variables_encrypted %}
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
    {%- endif %}
  }
}

{%- set definition = component.Definition -%}

{% if "sentry" in definition.Integrations and global.SentryConfig.AuthToken -%}
resource "sentry_key" "{{ component.Name }}" {
  organization = {{ global.SentryConfig.Organization|tf }}
  project = {{ component.Sentry.Project|tf }}
  name = "{{ site.Identifier }}_{{ component.Name }}"
  {% if component.Sentry.RateLimitWindow -%}
  rate_limit_window = {{ component.Sentry.RateLimitWindow }}
  {%- endif %}
  {% if component.Sentry.RateLimitCount -%}
  rate_limit_count = {{ component.Sentry.RateLimitCount }}
  {%- endif %}
}
{%- endif %}

module "{{ component.Name }}" {
  source            = "{{ definition.Source|safe }}{% if definition.UseVersionReference() %}?ref={{ definition.Version }}{% endif %}"

  {% if component|has_cloud_integration or component.Variables %}
  variables = {
    {% for key, value in component.Variables -%}
    {{ key }} = {{ value|tfvalue }}
    {% endfor %}
  }
  {% endif %}

  {% if component|has_cloud_integration or component.Secrets -%}
  secrets = {
    {% for key, value in component.Secrets -%}
    {{ key }} = {{ value|tfvalue }}
    {% endfor %}
  }
  {%- endif %}

  {% if component|has_cloud_integration -%}
  component_version       = {{ definition.Version|tf }}
  environment             = {{ global.Environment|tf }}
  site                    = "{{ site.Identifier }}"
  tags                    = local.tags
  {%- endif %}

  {% if "azure" in definition.Integrations -%}
  {% include "./component_azure_variables.tf" with definition=definition %}
  {%- endif %}

  {% if "aws" in definition.Integrations -%}
  {% include "./component_aws_variables.tf" with definition=definition %}
  {%- endif %}

  {% if "sentry" in definition.Integrations %}
  sentry_dsn              = {% if global.SentryConfig.AuthToken %}sentry_key.{{ component.Name }}.dsn_secret{% else %}"{{ component.Sentry.DSN }}"{% endif %}
  {% endif %}

  {% if "commercetools" in definition.Integrations %}
    ct_project_key    = {{ site.Commercetools.ProjectKey|tf }}
    ct_api_url        = {{ site.Commercetools.ApiURL|tf }}
    ct_auth_url       = {{ site.Commercetools.TokenURL|tf }}

    ct_stores = {
      {% for store in site.Commercetools.Stores %}
      {{ store.Key }} =  {
        key = {{ store.Key|tf }}
        variables = {
          {% for key, value in component.store_variables|get:store.key %}
          {{ key }} = {{ value|tfvalue }}
          {% endfor %}
        }
        secrets = {
          {% for key, value in component.store_secrets|get:store.key %}
          {{ key }} = {{ value|tfvalue }}
          {% endfor %}
        }
      }
      {% endfor %}
  }

  {% endif %}

  {% if "contentful" in component.Integrations %}
    contentful_space_id = contentful_space.space.id
  {% endif %}

  {% if "amplience" in component.Integrations %}
    amplience_client_id = {{ site.amplience.client_id|tf }}
    amplience_client_secret = {{ site.amplience.client_secret|tf }}
    amplience_hub_id = {{ site.amplience.hub_id|tf }}
  {% endif %}

  {% if "apollo_federation" in component.Integrations %}
    apollo_federation = {
      api_key       = {{ site.apollo_federation.api_key|tf }}
      graph         = {{ site.apollo_federation.graph|tf }}
      graph_variant = {{ site.apollo_federation.graph_variant|tf }}
    }
  {% endif %}

  providers = {
    {% if "azure" in definition.Integrations -%}azurerm = azurerm{%- endif %}
    {% if "aws" in definition.Integrations -%}
      aws = aws
      {% for provider in site.AWS.ExtraProviders -%}
      aws.{{ provider.Name }} = aws.{{ provider.Name }}
      {% endfor %}
    {%- endif %}
  }

  depends_on = [
    {% if site.AWS and component.Endpoint %}
    aws_apigatewayv2_api.{{ component.Endpoint|slugify }}_gateway,
    {% endif %}
    {% if site.azure and component.Azure.service_plan %}
    {% if component.azure.service_plan == 'default' %}
    azurerm_app_service_plan.functionapps,{% else %}
    azurerm_app_service_plan.functionapps_{{ component.azure.service_plan }},{% endif %}
    {% endif %}
    {% if site.Commercetools %}
    null_resource.commercetools,
    {% endif %}
  ]
}

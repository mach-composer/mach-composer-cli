# Module: {{ component.Name }}
{% if "sentry" in definition.Integrations and global.SentryConfig.AuthToken -%}
resource "sentry_key" "{{ component.Name }}" {
  organization      = {{ global.SentryConfig.Organization|tf }}
  project           = {{ component.Sentry.Project|tf }}
  name              = "{{ site.Identifier }}_{{ component.Name }}"
  {% if component.Sentry.RateLimitWindow -%}
  rate_limit_window = {{ component.Sentry.RateLimitWindow }}
  {%- endif %}
  {% if component.Sentry.RateLimitCount -%}
  rate_limit_count  = {{ component.Sentry.RateLimitCount }}
  {%- endif %}
}
{%- endif %}

module "{{ component.Name }}" {
  source            = "{{ definition.Source|safe }}{% if definition.UseVersionReference() %}?ref={{ definition.Version }}{% endif %}"

  {% if component.HasCloudIntegration() || component.Variables %}
  variables = {
    {% for key, value in component.Variables -%}
    {{ key }} = {{ value|tfvalue }}
    {% endfor %}
  }
  {% endif %}

  {% if component.HasCloudIntegration() || component.Secrets -%}
  secrets = {
    {% for key, value in component.Secrets -%}
    {{ key }} = {{ value|tfvalue }}
    {% endfor %}
  }
  {%- endif %}

  {% if component.HasCloudIntegration() -%}
  component_version       = {{ definition.Version|tf }}
  environment             = {{ global.Environment|tf }}
  site                    = "{{ site.Identifier }}"
  tags                    = local.tags
  {%- endif %}

  {% if "azure" in definition.Integrations %}{% include "./component/azure_variables.tf" %}{% endif %}
  {% if "aws" in definition.Integrations %}{% include "./component/aws_variables.tf" %}{% endif %}
  {% if "commercetools" in definition.Integrations %}{% include "./component/commercetools_variables.tf" %}{% endif %}

  {% if "sentry" in definition.Integrations %}
  sentry_dsn              = {% if global.SentryConfig.AuthToken %}sentry_key.{{ component.Name }}.dsn_secret{% else %}"{{ component.Sentry.DSN }}"{% endif %}
  {% endif %}

  {% if "contentful" in definition.Integrations %}
    contentful_space_id = contentful_space.space.id
  {% endif %}

  {% if "amplience" in definition.Integrations %}
    amplience_client_id = {{ site.Amplience.ClientID|tf }}
    amplience_client_secret = {{ site.Amplience.ClientSecret|tf }}
    amplience_hub_id = {{ site.Amplience.HubID|tf }}
  {% endif %}

  {% if "apollo_federation" in definition.Integrations %}
    apollo_federation = {
      api_key       = {{ site.ApolloFederation.APIKey|tf }}
      graph         = {{ site.ApolloFederation.Graph|tf }}
      graph_variant = {{ site.ApolloFederation.GraphVariant|tf }}
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
    {% if site.Azure and component.Azure.ServicePlan %}
    {% if component.Azure.ServicePlan == "default" %}
    azurerm_app_service_plan.functionapps,{% else %}
    azurerm_app_service_plan.functionapps_{{ component.Azure.ServicePlan }},{% endif %}
    {% endif %}
    {% if site.Commercetools %}
    null_resource.commercetools,
    {% endif %}
  ]
}

{% if variables.Encrypted %}
data "local_file" "variables" {
  filename = "{{ variables.Filepath }}"
}

data "sops_external" "variables" {
  source     = data.local_file.variables.content
  input_type = "yaml"
}
{% endif %}

{%- if global.SentryConfig.AuthToken %}
provider "sentry" {
  token = {{ global.SentryConfig.AuthToken|tf }}
  base_url = {% if global.SentryConfig.BaseURL %}{{ global.SentryConfig.BaseURL|tf }}{% else %}"https://sentry.io/api/"{% endif %}
}
{%- endif %}


{% if site.Commercetools %}{% include "services/commercetools/main.tf" with commercetools=site.Commercetools %}{% endif %}
{% if site.Contentful %}{% include "services/contentful.tf" with contentful=site.Contentful %}{% endif %}
{% if site.Amplience %}{% include "services/amplience.tf" with amplience=site.Amplience %}{% endif %}
{% if site.AWS %}{% include "services/aws/main.tf" with aws=site.AWS %}{% endif %}
{% if site.Azure %}{% include "services/azure/main.tf" with azure=site.Azure %}{% endif %}

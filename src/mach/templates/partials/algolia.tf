{% set algolia = site.algolia %}

provider "algolia" {
  app_id  = {{ algolia.application_id|tf }}
  api_key = {{ algolia.api_key|tf }}
}

{% set contentful = site.contentful %}

provider "contentful" {
  cma_token       = {{ contentful.cma_token|tf }}
  organization_id = {{ contentful.organization_id|tf }}
}

resource "contentful_space" "space" {
  name           = {{ contentful.space|tf }}
  default_locale = {{ contentful.default_locale|tf }}
}

resource "contentful_apikey" "apikey" {
  space_id = contentful_space.space.id

  name        = "frontend"
  description = "MACH generated frontend API key"
}

output "contentful_space_id" {
  value = contentful_space.space.id
}

output "contentful_apikey_access_token" {
  value = contentful_apikey.apikey.access_token
}

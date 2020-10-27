{% set contentful = site.contentful %}

provider "contentful" {
  cma_token       = "{{ contentful.cma_token }}"
  organization_id = "{{ contentful.organization_id }}"
}

resource "contentful_space" "space" {
  name           = "{{ contentful.space }}"
  default_locale = "{{ contentful.default_locale }}"
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

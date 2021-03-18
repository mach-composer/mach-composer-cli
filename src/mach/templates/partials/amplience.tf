{% set amplience = site.amplience %}

provider "amplience" {
  client_id        = amplience.client_id
  client_secret    = amplience.client_secret
  hub_id           = amplience.hub_id
}

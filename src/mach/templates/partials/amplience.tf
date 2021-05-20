{% set amplience = site.amplience %}

provider "amplience" {
  client_id        = {{ amplience.client_id|tf }}
  client_secret    = {{ amplience.client_secret|tf }}
  hub_id           = {{ amplience.hub_id|tf }}
}

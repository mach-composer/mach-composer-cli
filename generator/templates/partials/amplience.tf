provider "amplience" {
  client_id        = {{ amplience.ClientID|tf }}
  client_secret    = {{ amplience.ClientSecret|tf }}
  hub_id           = {{ amplience.HubID|tf }}
}

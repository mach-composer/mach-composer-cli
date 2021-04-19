output "shared_resource_group" {
  value = module.shared_infra.resource_group_name
}

output "ssl_key_vault_name" {
  value = module.shared_infra.ssl_key_vault_name
}

output "ssl_key_vault_secret_name" {
  value = module.shared_infra.ssl_key_vault_secret_name
}

output "ssl_key_vault_secret_version" {
  value = module.shared_infra.ssl_key_vault_secret_version
}
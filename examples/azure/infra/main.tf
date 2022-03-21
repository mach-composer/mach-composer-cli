provider "azurerm" {
  version                    = "=2.37.0"
  subscription_id            = var.subscription_id
  tenant_id                  = var.tenant_id
  skip_provider_registration = true
  features {}
}

module "shared_infra" {
  source                        = "git::https://github.com/labd/terraform-azure-mach-shared.git"
  name_prefix                   = "mach-shared-we"
  region                        = "westeurope"
  dns_zone_name                 = var.dns_zone_name
  certificate_access_object_ids = [
    # Fill in any object IDs of users, groups or service principles that should be able to
    # access the certificates
  ]
}
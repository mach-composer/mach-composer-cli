provider "google" {
  project = var.project_id
  region = var.region
}

module "shared_infra" {
  source                        = "git::https://github.com/labd/terraform-gcp-mach-shared.git"
  name_prefix                   = "mach-shared-we"
  region                        = "westeurope"
  dns_zone_name                 = var.dns_zone_name
  certificate_access_object_ids = [
    # Fill in any object IDs of users, groups or service principles that should be able to
    # access the certificates
  ]
}

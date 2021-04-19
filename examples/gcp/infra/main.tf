provider "google" {
  project     = var.project_id
  region      = var.region
}

module "shared_infra" {
  source      = "git::https://github.com/labd/terraform-gcp-mach-shared.git"
  region      = var.region
  name_prefix = "mach-playground"
}

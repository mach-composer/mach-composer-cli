{% set gcp = site.gcp %}

provider "google" {
  project = "{{ gcp.project_id }}"
}


locals {
  project_id  = "{{ gcp.project_id }}"
  region      = "{{ gcp.region }}"
  project_key = "{{ site.commercetools.project_key }}"

  tags = {
    Site = "{{ site.identifier }}"
    Environment = "{{ general_config.environment }}"
  }
}

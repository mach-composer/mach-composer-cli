# This file is auto-generated by MACH composer
# SiteComponent: test-1/component-2
terraform {
  backend "local" {
    path = "./states/test-1/component-2.tfstate"
  }
  required_providers {}
}

# File sources
# Resources
data "terraform_remote_state" "test-1" {
  backend = "local"
  config = {
    path = "./states/test-1.tfstate"
  }
}

# Component: component-2
module "component-2" {
  source = "/home/thomas/Projects/mach-composer/mach-composer-cli/internal/cmd/testdata/modules/application"
  variables = {
    parent_names = [data.terraform_remote_state.test-1.outputs.component-1.name]
  }
}

output "component-2" {
  description = "The module outputs for component-2"
  sensitive   = true
  value       = module.component-2
}

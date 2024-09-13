locals {
  module_path_elements = split("/", path.module)
  module_name          = element(local.module_path_elements, length(local.module_path_elements) - 1)
}

resource "time_sleep" "sleep" {
  create_duration = var.variables.sleep
}

variable "variables" {
  type = object({
    parent_names = list(string)
    sleep        = optional(string, "0s")
    fail         = optional(bool, false)
    random_value = optional(string, "default")
  })
}

data "http" "example" {
  count = var.variables.fail == true ? 1 : 0
  url   = "fails"
}

resource "local_file" "child" {
  content = jsonencode({
    name         = local.module_name
    parent_names = var.variables.parent_names
  })
  filename = "${path.cwd}/outputs/${local.module_name}.json"

  depends_on = [
    time_sleep.sleep
  ]
}

output "name" {
  value = local.module_name
}

output "random_value" {
  value = var.variables.random_value
}


output "array" {
  value = []
}

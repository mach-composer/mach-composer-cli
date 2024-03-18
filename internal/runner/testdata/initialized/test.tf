terraform {
  backend "local" {
    path = "./terraform.tfstate"
  }
  required_providers {}
}

output "some-output" {
  value = "hello-world"
}

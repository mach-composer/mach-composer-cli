target "docker-metadata-action" {}

target "cli" {
  inherits   = ["docker-metadata-action"]
  context    = "./"
  dockerfile = "docker/cli.Dockerfile"
  platforms = [
    "linux/amd64",
  ]
}

target "base" {
  inherits   = ["docker-metadata-action"]
  context    = "./"
  dockerfile = "docker/base.Dockerfile"
  platforms = [
    "linux/amd64",
  ]
  target = "base"
}

target "default-all" {
  inherits = ["docker-metadata-action"]
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/all.Dockerfile"
}

target "default-aws" {
  inherits = ["docker-metadata-action"]
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/aws.Dockerfile"
}

target "default-azure" {
  inherits = ["docker-metadata-action"]
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/azure.Dockerfile"
}

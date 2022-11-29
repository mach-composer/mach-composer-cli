
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
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/all.Dockerfile"
}

target "default-aws" {
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/aws.Dockerfile"
}

target "default-azure" {
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/azure.Dockerfile"
}

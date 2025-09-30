target "docker-metadata-action" {}

group "cli" {
  targets = ["cli-arm64", "cli-amd64"]
}

target "cli-arm64" {
  inherits   = ["docker-metadata-action"]
  context    = "./"
  dockerfile = "docker/cli.Dockerfile"
  platforms = [
    "linux/arm64",
  ]
  args = {
    GOOS   = "linux"
    GOARCH = "arm64_v8.0"
  }
}
target "cli-amd64" {
  inherits   = ["docker-metadata-action"]
  context    = "./"
  dockerfile = "docker/cli.Dockerfile"
  platforms = [
    "linux/amd64",
  ]
  args = {
    GOOS   = "linux"
    GOARCH = "amd64_v1"
  }
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

target "default-gcp" {
  inherits = ["docker-metadata-action"]
  contexts = {
    base = "target:base"
  }
  dockerfile = "docker/gcp.Dockerfile"
}

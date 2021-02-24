# Apollo Federation

You can define a site specific configuration for [managed Apollo federation](https://go.apollo.dev/s/managed-federation).

## Retrieve credentials

Go to Apollo Studio and you will be prompted with a popup for the credentials.

The *graph* setting can be determined from the **APOLLO_KEY** string, use the name between the two `:` characters.

## Example site configuration block

An example site block:

```yaml
apollo_federation:
  api_key: service:mach-poc-123:Abc00kHbB89h
  graph: mach-poc-123
  graph_variant: current
```
      
## Expose Apollo Federation to components

MACH needs to know what components want to use the Apollo Federation configuration.<br>
For this you need to include `apollo_federation` to the list of integrations.<br>
When doing so, MACH expects the component to have one variable `apollo_federation` defined ([more info](../../reference/components/structure.md#apollo federation))

## Managed Apollo Federation gateway setup

Make sure you set the following environment variables:

- APOLLO_KEY
- APOLLO_GRAPH_VARIANT

## Downstream federated services apollo integration

Make sure the `apollo` cli executable is available in the context where you're running mach.

While there are some [edge cases](https://www.apollographql.com/docs/federation/managed-federation/deployment/#pushing-configuration-updates-safely),
in general you want to push after the new service is completely deployed. An example on how to do this:

```terraform
locals {
  # URL which the gateway uses to access this service
  service_url       = "http://internal.some-gateway.com/some-component/graphql"
  # Exported schema of the service that is being deployed
  local_schema_file = "${path.module}/generated/schema.graphql"
  # identifier of the service for Apollo Studio
  service_name      = "some-component"
}

resource "null_resource" "apollo_push_schema_after_deploy" {
  triggers = {
    schema_hash  = filesha256(local.local_schema_file)
    service_name = local.service_name
    graph        = var.apollo_federation.graph
    variant      = var.apollo_federation.variant
  }

  provisioner "local-exec" {
    command = "apollo service:push --localSchemaFile=${local.local_schema_file} --key=${var.apollo_federation.api_key} --graph=${var.apollo_federation.graph} --variant=${var.apollo_federation.graph_variant} --serviceName=${local.serviceName} --serviceURL=${local.service_url}"
  }

  depends_on = [
    # list dependencies to make sure this is run before or after deployment of the service, f.e. an ecs service
    aws_ecs_service.some_component,
  ]
}
```

### Deleting a downstream service

Deleting a service automatically is possible with some [caveats](https://github.com/apollographql/apollo-tooling/issues/2115):

```terraform
resource "null_resource" "apollo_push_schema_after_deploy" {
  triggers = {
    schema_hash = filesha256(local.local_schema_file)
    # for destroy we can't use local variables, so include them in triggers
    service_name = local.service_name
    api_key = var.apollo_federation.api_key
    graph = var.apollo_federation.graph
    variant = var.apollo_federation.variant
  }

  provisioner "local-exec" {
    command = "apollo service:push --localSchemaFile=${local.local_schema_file} --key=${var.apollo_federation.api_key} --graph=${var.apollo_federation.graph} --variant=${var.apollo_federation.graph_variant} --serviceName=${local.serviceName} --serviceURL=${local.service_url}"
  }

  # note this command fails if it's the last federated service
  provisioner "local-exec" {
    when = destroy
    command = "apollo service:delete --yes --key=${self.triggers.api_key} --serviceName=${self.triggers.service_name} --variant=${self.triggers.variant} || true"
  }
}
```

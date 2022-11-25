# Apply Apollo Federation

A common use-case when having multiple components (and/or external SaaS
services) with a GraphQL API is to expose the APIs using a
[federated gateway](https://www.apollographql.com/docs/federation/).

This setup consists of two things:

- A GraphQL (federated) Gateway component
- One or more components exposing their own GraphQL API


## Managed Apollo Federation Gateway setup

The Gateway components needs to implement the `ApolloGateway` and be connected
to the Apollo Studio configured in the [Site
configuration](../reference/syntax/sites.md#apollo_federation).

Any component that uses the [Apollo Federation integration](../reference/components/structure.md#apollo-federation)
must have the following Terraform variable defined:

```terraform
variable "apollo_federation" {
  type = object({
    api_key       = string
    graph         = string
    graph_variant = string
  })
}
```

Make sure that in your runtime settings you set the following environment variables defined:

```terraform
APOLLO_KEY = var.apollo_federation.api_key
APOLLO_GRAPH_VARIANT = var.apollo_federation.graph_variant
```

## Downstream federated services apollo integration

Once a component is deployed (either added to your MACH stack or updated) it
needs to update its schema in the Apollo Federation.

To do so we need to make this update call part of our deployment processes.

!!! important "Apollo CLI"
    Make sure the [Apollo CLI](https://www.apollographql.com/docs/devtools/cli/)
    is available in the context where you're running mach.

### Updating a downstream service

While there are some [edge cases](https://www.apollographql.com/docs/federation/managed-federation/deployment/#pushing-configuration-updates-safely),
in general you want to push after the new service is completely deployed. An
example on how to do this:

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

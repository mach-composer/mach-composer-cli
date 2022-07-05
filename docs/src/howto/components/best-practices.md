# Best practices

Here's a list of best practices regarding MACH components.

## Check artifacts existence 

We often see builds fail with bad error messages (403 for Lambda functions) or silently for Docker if the artifact (zip or Docker container)
does not exist. This often happens because a CI step failed (linting or testing).

We can fail early in Terraform by explicitly checking for the artifact, and doing a `depends_on` on this data resource:

For AWS ECR (Docker):

```
data "aws_ecr_image" "mach_graphql" {
  registry_id     = local.registry_id
  repository_name = local.graphql_image_name
  image_tag       = var.component_version
}
```

For AWS lambda:

```
data "aws_s3_bucket_object" "component_exists" {
  bucket = <your lambda s3 bucket>
  key    = <component name.zip>
}
```

For Azure blob storage there is no data resource, but one can script it and hook it up to an [external data source](https://registry.terraform.io/providers/hashicorp/external/latest/docs/data-sources/data_source):

```bash
az storage blob exists --account-key $1 --account-name $2 --container-name $3 --name $4 --query exists
```


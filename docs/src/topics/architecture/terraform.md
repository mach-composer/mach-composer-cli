# Everything-as-configuration-as-code

As a principle, MACH composer encourages declarative defining your
configuration as code, as much as possible. This will enable you to manage much
of your infrastructure in an automated fashion, which in turn will enable you to
scale your MACH ecosystem to many instances and use-cases.


## Terraform & infrastructure-as-code

In most cases, this boils down to managing your SaaS and Cloud resources using
Terraform, which MACH composer uses underneath as the infrastructure-as-code
engine. As a rule-of-thumb we can say: *if it can be managed by Terraform, it
can be managed by MACH composer*.

Find supported services on these pages:

- [Official providers](https://www.terraform.io/docs/providers/index.html)
- [Community providers](https://www.terraform.io/docs/providers/type/community-index.html)


There are services however, that don't have a Terraform provider available. So
you could go on and create a Terraform provider yourself, [just like Lab Digital
did with the Terraform Provider for Commercetools](https://blog.labdigital.nl/commercetools-terraform-a-match-made-in-heaven-1d7a48e4931b).


!!! info "Creating your own Terraform Provider"
    Read [this blogpost by Hashicorp](https://www.hashicorp.com/resources/creating-terraform-provider-for-anything)
    to find out how to start with that. The article also points to
    [this Tutorial](https://learn.hashicorp.com/collections/terraform/providers),
    which explains how to build your own Terraform Provider.


When you don't have the time or the resources to create your own Terraform
provider, there are other options (that we're not shy of using in production
ourselves) that allow you to apply infra-as-code principles:

### 1. Use the generic REST API Terraform provider

Mastercard has made available their open source provider, for managing generic
RESTful webservices. We've used it in production for SaaS services that don't
have their own Terraform provider.

Find it here: [https://github.com/Mastercard/terraform-provider-restapi](https://github.com/Mastercard/terraform-provider-restapi)


### 2. Use the generic GraphQL Terraform provider

Same as the generic REST API provider, a generic GraphQL Terraform provider is
available, that you can use to manage GraphQL resources.

Find it here: [https://github.com/sullivtr/terraform-provider-graphql](https://github.com/sullivtr/terraform-provider-graphql)


### 3. Invoke an existing CLI from Terraform

Using Terraform's [null_resource](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource),
you can invoke local commands.

```terraform
resource "null_resource" "my_resource" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = format("%s/my-cli command --parameter=%s", path.module, var.my_parameter)
  }
}
```

!!! warning "Last resort"
    Using the `null_resource` option is a considered a last resort. Consider
    using one of the other options, as most likely these CLIs will use the API
    as well underneath.

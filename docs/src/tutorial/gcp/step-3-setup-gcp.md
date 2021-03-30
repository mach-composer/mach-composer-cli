# Step 3. Prepare your GCP environment

!!! todo
    GCP is not yet supported by MACH composer.<br>
    We are happy to accept contributions to add this! And think with you if you plan to work on it. Feel free to reach out to us via [opensource@labdigital.nl](mailto:opensource@labdigital.nl).

## What needs to happen to support GCP?

Luckily, adding GCP should not be a lot of work. Most of the implementation boils down to adding the right terraform code in MACH composer, to generate the nessesary resources. This is primarily about setting up an API gateway that ties together many services, through [GCP's Terraform support](https://registry.terraform.io/providers/hashicorp/google/latest/docs).

- [ ] Decide on what [services](#infra-decisions) to use
- [ ] Add support in MACH yaml to support GCP cloud
- [ ] Add terraform templates to MACH composer, to support GCP resources
- [ ] Add support for storing terraform state remotely in GCP
- [ ] Ideally: extend component bootstrapper and component cookiecutter to include Google Cloud setup
- [ ] Expand documentation with GCP

## Infra decisions

For the following solutions we have in place for AWS and Azure, we need to decide on the service we want to use in Google Cloud:

- [x] Terraform state
- [x] HTTP routing / API gateway
- [x] Custom domains / DNS

For reference implementations of components:

- [x] Secrets management
- [x] Serverless functions

!!! note
    For all services, we need to take into account that it should be deployed through terraform completely.<br>
    Though solutions exist when there is no 'full blown' terrafrom support, in which case you could fall back to a CLI. We recently did this to implement [Apollo Studio support](https://github.com/labd/mach-composer/pull/78).

### Terraform state backend

Use [Google Cloud Storage](https://www.terraform.io/docs/language/settings/backends/gcs.html) as documented by Terraform.

### HTTP routing / API gateway

Use [API Gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/api_gateway_api).

- [Quickstart](https://cloud.google.com/api-gateway/docs/get-started-cloud-functions)
- [Beta support for it in terraform](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/api_gateway_api)


**Custom domains**

Custom domain names are not supported by API Gateway. 
For custom domains, we need to create a load balancer and direct requests to the `gateway.dev` domain of the deployed API.

- [Setting up load balancing 'the hard way'](https://cloud.google.com/blog/topics/developers-practitioners/serverless-load-balancing-terraform-hard-way)
- [API gateway behind load balancer](https://medium.com/swlh/google-api-gateway-and-load-balancer-cdn-9692b7a976df)
- [Load balancing with Terraform](https://cloud.google.com/community/tutorials/modular-load-balancing-with-terraform)

### Custom domains / DNS

[Google Cloud DNS](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dns_record_set)

### Secrets management

[GCP secrets manager](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret).

### Serverless function

- Serverless functions: [Google Cloud Functions](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudfunctions_function)
    - Example for commercetools api extensions
    - Example for commercetools subscriptions
    - Example for generic API
- Serverless Docker containers: [Google CloudRun](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service)
    - Example for generic API

Interesting read: [https://dev.to/didil/gcp-api-gateway-demo-with-terraform-go-cloud-run-3o9e](https://dev.to/didil/gcp-api-gateway-demo-with-terraform-go-cloud-run-3o9e)

## Multi tenancy

We need to determine what 'project structure' will be used in GCP. 
In other clouds we use multiple **resource groups** (Azure) and **accounts** (AWS) to provide multi-tenancy/platform partitioning between sites. 
The same should be achieved with GCP.
### Useful links:

  - [https://cloud.google.com/solutions/best-practices-vpc-design](https://cloud.google.com/solutions/best-practices-vpc-design)
  - [https://cloud.google.com/docs/enterprise/best-practices-for-enterprise-organizations](https://cloud.google.com/docs/enterprise/best-practices-for-enterprise-organizations)

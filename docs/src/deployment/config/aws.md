# AWS deployments

## API Gateway

Only when a MACH stack contains components that are configures as [`has_public_api`](../../syntax.md#components), MACH will setup the necessary resources to be able to route traffic to that component:

- API Gateway
- API Gateway default routing
- Route53 record on the zone configured with `base_url`
- ACM Certificate (with DNS validation)

The information needed for components to add custom routes to that API Gateway are provided through [Terraform variables](../../components/aws.md#terraform-variables).

!!! info "Route53 zone"
    MACH will not create and manage the Route53 zone itself but expects it to be created already as described in the [prerequisites](../../prerequisites/aws.md#route53-zone) section.<br>
    It will try to lookup that zone using the `base_url` setting.
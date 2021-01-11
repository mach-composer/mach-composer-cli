# AWS deployments

## HTTP routing

If any component used in a MACH stack is defined with an [`endpoint`](../../syntax.md#components), MACH will create the necessary resources to be able to route traffic to that components.

The information needed for components to add custom routes to that API Gateway are provided through [Terraform variables](../../components/aws.md#terraform-variables).

### Default endpoint

If you have defined your component with a `default` endpoint, MACH will create an API Gateway for you which includes a default AWS API Gateway domain.

```
components:
  - name: payment
    source: git::ssh://git@github.com/your-project/components/payment-component.git//terraform
    endpoints: 
      public: default
    version: ....
```

!!! note ""
    This `default` endpoint doesnt need to be defined in your [endpoints definition](../../syntax.md#endpoints).

### Custom endpoint

A site might have a couple of [endpoints](../../syntax.md#endpoints) defined and for each endpoint MACH will create:

- API Gateway + default routing
- ACM Certificate (with DNS validation)
- Route53 record on the zone auto-detected or configured on the endpoint


!!! info "Route53 zone"
    MACH will not create and manage the Route53 zone itself but expects it to be created already as described in the [prerequisites](../../prerequisites/aws.md#route53-zone) section.<br>
    It will try to lookup that zone using the `route53_zone_name` setting.
# AWS deployments

## HTTP routing

{% include-markdown "./endpoints_deprecated.md" %}

If any component used in a MACH composer stack is defined with an
[`endpoint`](../../../reference/syntax/component.md), MACH composer will create
the necessary resources to be able to route traffic to that components.

The information needed for components to add custom routes to that API Gateway
are provided through [Terraform variables](../../../reference/syntax/index.md#variables).

### Default endpoint

If you have defined your component with a `default` endpoint, MACH composer will
create an API Gateway for you which includes a default AWS API Gateway domain.

```
components:
  - name: payment
    source: git::ssh://git@github.com/your-project/components/payment-component.git//terraform
    endpoints:
      public: default
    version: ....
```

!!! note ""
    This `default` endpoint doesn't need to be defined in your [endpoints' definition](../../../reference/syntax/site.md#nested-schema-for-endpoints).

### Custom endpoint

A site might have a couple of [endpoints](../../../reference/syntax/site.md#nested-schema-for-endpoints)
defined and for each endpoint MACH composer will create:S

- API Gateway + default routing
- ACM Certificate (with DNS validation)
- Route53 record on the zone auto-detected or configured on the endpoint


!!! info "Route53 zone"
    MACH composer will not create and manage the Route53 zone itself but expects
    it to be created already as described in the
    [prerequisites](../../../tutorial/aws/step-4-setup-aws-site.md) section.<br>
    It will try to look up that zone using the `route53_zone_name` setting.

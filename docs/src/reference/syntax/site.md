# Site

A site describes which available components make up a specific site, and
what configurations they have.

## Schema

## Required

- `identifier` - (Required) Unique identifier for this site. Will be used for
  the Terraform state and naming all cloud resources.
- `components` (List of Block)
  [Component configurations](#nested-schema-for-components)

## Optional

- `deployment` (Block) [Deployment configuration](#nested-schema-for-deployment)
- `endpoints` (Map of String, _deprecated_)
  [Endpoint definitions](#nested-schema-for-endpoints) to be used in the
  API Gateway or Frontdoor routing

### Dynamic

{% include-markdown "./dynamic.md" %}

## Nested schema for `components`

Configures the components for the site. They must reference a defined component
(defined in the [component definitions](component.md)). A site component is
the smallest deployable unit in mach-composer. It is a combination of a
component and specific site configuration.

### Example

```yaml
components:
  - name: api-extensions
    variables:
      ORDER_PREFIX: mysitetst
    depends_on:
      - order-mailer
  - name: order-mailer
    variables:
      FROM_EMAIL: mach@example.com
    secrets:
      SENDGRID_API_KEY: my-api-token
    store_variables:
      brand-a:
        FROM_EMAIL: mach@brand-a.com
      other-brand:
        FROM_EMAIL: mach@other-brand.com

```

### Required

- `name` (String) Reference to a [component](./component.md) definition

### Optional

- `variables` (Block) Variables for this component
- `secrets` (Block) Variables for this component that should be stored
  in an encrypted key-value store
- `deployment` (Block) [Deployment configuration](#nested-schema-for-deployment)
- `depends_on` (List of String) allows for the explicit setting of dependencies
  between components. This is useful when a component depends on another
  component, but the dependency cannot be inferred from the component
  definition. For example, when a component depends on a component that is not
  deployed to the same cloud provider. The value of `depends_on` is the name of
  the component it depends on. This will overload any inferred relations.
  See [deployment](../../concepts/deployment/index.md) for more information.

### Dynamic

{% include-markdown "./dynamic.md" %}

## Nested schema for `deployment`

{% include-markdown "./deployment.md" %}

## Nested schema for `endpoints`

Endpoint definitions to be used in the API Gateway or Frontdoor routing.

Each component might require a different endpoint. In the
[component definition](./component.md) it can be defined which endpoint it
expects. The actual endpoint can be defined here using the unique key.

### Example

```yaml
endpoints:
  main: api.tst.mach-example.net
  services: services.tst.mach-example.net
```

Complex example:

```yaml
endpoints:
  internal:
    url: internal-api.tst.mach-example.net
    zone: tst.mach-example.net
    aws:
      throttling_burst_limit: 5000
      throttling_rate_limit: 10000
      enable_cdn: true
```

### Required

- `url` (String) url of the endpoint

### Optional

- `zone` (String)  DNS zone to use, if missing it's determined based on the
  given `url`
- `aws` (Block)  Configuration block
  for [AWS-specific settings](#nested-schema-for-aws)
- `azure` (Block)  Configuration block
  for [Azure-specific settings](#nested-schema-for-azure)

### Nested schema for `aws`

### Example

```yaml
aws:
  throttling_burst_limit: 5000
  throttling_rate_limit: 10000
  enable_cdn: true
```

### Optional

- `throttling_burst_limit` - Set burst limit for API Gateway endpoints
- `throttling_rate_limit` - Set burst limit for API Gateway endpoints
- `enable_cdn` - Defaults to false. Sets a CDN in front of this endpoint for
  better global availability. For AWS creates a CloudFront distribution

### Nested schema for `azure`

### Example

```yaml
azure:
  session_affinity_enabled: True
  session_affinity_ttl_seconds: 3600
  waf_policy_id: "string"
```

### Optional

- `session_affinity_enabled` - Whether to allow session affinity on this host
- `session_affinity_ttl_seconds` - The TTL to use in seconds for session
  affinity, if applicable.
- `waf_policy_id` - Defines the Web Application Firewall policy `ID` for the
  host.


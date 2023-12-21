# Public components

Components that have a public git repository and are deployed to a public
registry can be used in any project.

Lab Digital is also committed into publishing several MACH components which can
be used.

!!! todo ""
    At this point, we have one public component. Over time, we will update this
    list with more components published either by Lab Digital or any other we
    select or gets recommended to us.

## AWS

#### commercetools token refresher
:fa-github: [https://github.com/labd/mach-component-aws-commercetools-token-refresher](https://github.com/labd/mach-component-aws-commercetools-token-refresher)

Creates a lambda function that can be set on a SecretsManager value as
auto-rotate lambda.

```yaml
sites:
  - identifier: some site
    components:
    - name: ct-refresher
...

components:
- name: ct-refresher
  source: git::https://github.com/labd/mach-component-aws-commercetools-token-refresher.git//terraform
  version: <git hash of version you want to release>
  integrations: ["aws", "commercetools", "sentry"]
```

## Azure

Nothing yet

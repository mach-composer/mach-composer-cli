# Overview

!!! info "New since MACH composer 2.14"

With projects increasing in size and complexity, it is often necessary to start
splitting out state for multiple components to speed up deployments and to
reduce the risk of deployment failures due to errors. This is inherent to the
way the underlying Terraform code processes changes.

## Terraform and state

In short, Terraform will process all resources in a configuration and check
their
required state against the current state of the target systems. If there are
changes, Terraform will update the state file and apply the changes. This means
that the more resources there are in a state file, the longer it will take to
process the changes. This even is the case if the update is only within a
single resource, as Terraform will need to fetch the current state for all the
configurations regardless. For a more in-depth explanation of why this is, see
the [Terraform docs](https://developer.hashicorp.com/terraform/language/state)

## Mach Composer and state

By default, Mach Composer will create a configuration file per site to manage
its dependent components. Updates to the state will then be processed in
parallel. For small projects this is perfectly acceptable as Terraform will only
have limited resources it needs to check and update. However, as more components
are added to a project applies will take longer.

To speed up deployments, Mach Composer allows you to split out the state for
separate components into independent configurations. This will allow you to run
multiple component deployments in parallel, speeding up the overall deployment
time. Mach Composer will also recognize when no changes occurred, skipping an
apply altogether. In this way it is possible to move out the configurations for
components that are unlikely to change often, reducing the overall deployment
time.

[//]: <> (@formatter:off)
!!! tip "Splitting state is optional"
    This feature is fully optional and can be enabled on a per component basis.
    If migrating from an existing codebase check
    our [howto on migrating existing component state](../../howto/state/migration.md)
[//]: <> (@formatter:on)

In addition, Mach Composer will also recognize the relationships between
components (based on various factors such as required dependent inputs or
explicit configurations) and will make sure that the order in which components
are applied fits with the component interdependencies. In this way you don't run
the risk of accidentally deploying a frontend component when the required
backend is not yet updated.
See [managing dependencies](./managing-dependencies.md) for more information.

In this way you can taylor your state to your project needs, by moving out
components on an as-needed basis as you notice that deployments are slowing
down. See below for how to configure your components to use separate state, as
well as tools to inspect your environment, and an explanation of how Mach
Composer will apply these changes.

### SEE ALSO

* [`mach-composer apply`](../../reference/cli/mach-composer_apply.md) - The
  apply command
* [component life cycle](../../concepts/components/lifecycle/index.md) - How a component lifecycle
  works
* [configuration](./configuration.md) - Configuring a site component to use
  separate state
* [dependency graph](./managing-dependencies.md) - How component dependencies
  are resolved
* [applying changes](./applying-changes.md) - What happens during an apply
  command

# Configuration

Mach Composer allows for several different configurations that together define
how components will be applied.

## `deployment`

The `deployment` sections define how the component will actually be deployed.
This configuration can be set at the global level, the site level or the
component level, depending on your needs. In this case, component configuration
will be checked first, then site and finally the global level.

```yaml
...
global:
  # Global deployment type set to site
  deployment:
    type: site
sites:
  - identifier: my-site
    # No site deployment type set, so inherited from global
    components:
      - name: my-component
        # deployment type set to site-component, so component will be deployed 
        # independently
        deployment:
          type: site-component
      - name: my-other-component
        # No component deployment type set, so inherited from site
components:
  - name: my-component
  - name: my-other-component

...
```

[//]: <> (@formatter:off)
!!! tip "Set deployment type at the global level"
    It is recommended to set the deployment type to `site-component` at the global
    level when starting a new project. This will ensure that all components are
    deployed in the same way, and no later migration is needed. The current default
    of `site` is only for backwards compatibility.
[//]: <> (@formatter:on)

## `dependes_on`

Although in most cases Mach Composer will be able to determine the correct
dependency order based on the component configuration, in some cases it is
necessary to explicitly define the dependencies between components. This can be
done using the `depends_on` configuration.

[//]: <> (@formatter:off)
!!! warning "Note that setting this configuration will override the automatic dependency resolution completely."
    This means that if you set this configuration, you will need to make sure that
    the dependencies are correct. If you are unsure, it is recommended to use the
    automatic dependency resolution.
[//]: <> (@formatter:on)

```yaml
sites:
  - identifier: my-site
    components:
      - name: my-component
      - name: my-other-component
        depends_on:
          # This tells Mach Composer that my-other-component depends on my-component explicitly
          - my-component
```

## The `--workers` parameter

Finally, the `--workers` parameter allows you to set the number of workers that
will process updates concurrently. This can be useful if you have a large number
of independent components that need to be updated, and you want to speed up the
process further

Note that this parameter is structurally limited by the automated batch
processing of changes done by Mach Composer.
See [applying changes](applying-changes.md) for more information on how batches
are determined.

!!! info "By default, a single worker is used"

[//]: <> (@formatter:off)
!!! warning "API Rate Limits"
    Take care with setting this parameter too high, as it can cause API rate 
    limits issues with the underlying SAAS vendors.
[//]: <> (@formatter:on)

```bash
mach-composer apply -f my-site.yaml --workers 10
```

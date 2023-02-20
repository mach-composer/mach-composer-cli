# Plug-in system
!!! info "New since MACH composer 2.6"

In order to make it more flexible and easier to add support for new MACH
services to MACH composer, we've introduced a plug-in system that allows anyone
to create a plugin.

If a plug-in is not availale on the `$PATH`, MACH composer will try to fetch the
plugin from the location that's given in the configuration. We use GitHub
releases for hosting plug-ins. Any GitHub organisation and repository can be
used to host a plugin.

The officially supported plug-ins are hosted at the [MACH composer github
organization](https://github.com/mach-composer?q=plugin&type=all&language=&sort=).

In order to add the AWS and Sentry plugins to your MACH composer configuration,
add the the below configuration to your MACH composer YAML configuration.

```yaml
mach_composer:
  version: 1
  plugins:
    aws:
      source: mach-composer/aws
      version: 0.1.0
    sentry:
      source: mach-composer/sentry
      version: 0.1.2
```

In order to be backwards compatible, MACH composer itself currently ships the
plug-ins that were previously part of MACH composer itself. When you specify
specific versions of plug-ins, these will alsways prevail and MACH composer will
try to download these.

## Overriding plugin behaviour

As the plug-ins are open source and hosted on GitHub Releases, overriding
behaviour of a MACH composer plug-in is quite easy. Just fork the repository,
implement your changes, and configure MACH composer to fetch the plug-in from
your own repository! Just make sure that the release pipelines have resulted in
an installable artifact (`goreleaser` will take care of this).

```yaml
mach_composer:
  plugins:
    aws:
      source: my-org/my-fork-of-aws
      version: 0.1.0
```

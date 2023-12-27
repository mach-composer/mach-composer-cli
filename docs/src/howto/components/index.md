# Overview

## Creating a component

The easiest way to get started with a component is by looking at the cookiecutter:

[`mach-component-cookiecutter`](https://github.com/labd/mach-component-cookiecutter)

We'll add a GitHub template scaffold project in the near future.

## Adding it to your stack

A component can be added to your MACH stack by including it in your MACH
Composer configuration.

It should be present in your:

1. [Component definitions](../../reference/syntax/component.md) so that MACH
   Composer knows where to find your component
2. [Site component configuration](../../reference/syntax/site.md#nested-schema-for-components) to
    include it in your MACH composer stack add site-specific configuration

The [tutorial](../../tutorial/aws/step-6-create-mach-stack.md) includes an
example of a configuration file with a component implemented.

## Using Serverless framework

The 'function'-part of a MACH component can easily be integrated with the
[serverless framework](https://www.serverless.com) for "*zero-friction serverless development*".

This gives you a couple of features:

- Easy local development of your function code including mocked infrastructure
  that might be needed for your setup
- Build & package your function by calling `sls package`.
  [More info about Packaging & Deployment](../../concepts/components/lifecycle/index.md)

!!! info "Serverless framework in your MACH deployment"
    Although we do encourage the usage of the serverless framework for
    development and packaging, we don't recommend using it for the actual
    MACH composer deployment itself.<br>
    More info about this in the [MACH configuration deployment notes](../../concepts/components/lifecycle/index.md).

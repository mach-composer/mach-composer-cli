# MACH concepts

MACH composer is intended as the 'center piece' of your MACH architecture, and
coordinates everything in terms of infrastructure, MACH services, sites and
components configuration, and deploying all of that. Underneath Terraform is
used to manage all of this.

A couple of concepts that are applied in MACH composer.

## Single YAML configuration file for an entire MACH stack
For your MACH composer project, a single YAML configuration file is used to
manage one or multiple MACH site configurations, your cloud infrastructure
required for components, and the integration between MACH services and these
cloud services.

## Composition through serverless components
With MACH composer, you can develop serverless (microservice architecture)
components and (re-)use and configure them across multiple environments or
projects; each environment can have their own composition of those reusable
components.

Within one YAML configuration, you can configure multiple sites that each have a
different composition of components.

## First class support for a number of MACH services and cloud platforms
MACH composer offers first-class support for a number of MACH services such as
commercetools, Amplience and contentful. Also it offers native support for the
major cloud platforms, AWS, Azure and soon GCP.

First-class support means that the YAML configuration may contain attributes
that are specific to the MACH service, as well as the ability to configure
components to integrate with a particular MACH or cloud service.

## Easy management of multiple MACH environments
MACH composer is intended to power multiple environments/sites from a single
source. By copy/pasting an existing configuration, and making the necessary
adjustments, MACH composer can rollout your new environment or sites.

## Integration with CI/CD process
MACH composer is intended to run in CI/CD context and can be integrated into
your (existing) CI/CD pipeline, for automated deployments of your entire stack.

This documentation offers guidance for implementing CI/CD workflows in GitHub
Actions, Gitlab and Azure DevOps

!!! tip "Further reading"
    - [What makes a component](./components/index.md)
    - [Architectural Guidance](./architecture/index.md)
    - [Getting started](../tutorial/step-1-installation.md)

# MACH composer

This documentation describes the workings of the MACH composer. Intended to setup and manage a **M**icroservices based, **A**PI-first, **C**loud-native SaaS and **H**eadless platform.

## How does it work?

The MACH composer takes a [YAML configuration](./syntax.md) as input, and will translate this into a Terraform configuration. It will then execute the terraform configuration, which will deploy all resources for the site architecture.

[![MACH diagram](./_img/mach.png)](./_img/mach.png)

The MACH composer is intended for managing multiple instances of the architecture.

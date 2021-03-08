# Configuration syntax

A configuration file can contain several sites with all different configurations and all using a different mix of re-usable serverless microservice components.

It is common to have a single configuration file per environment since they usually share the same general configurations.

The configuration file has the following structure:

- **[global](./global.md)**
    - **[environment](./global.md)**
    - **[terraform_config](./global.md#terraform_config)**
    - **[cloud](./global.md)**
    - [azure](./global.md#azure)
    - [sentry](./global.md#sentry)
    - [contentful](./global.md#contentful)
    - [amplience](./global.md#amplience)
- **[sites](./sites.md)**
    - **[identifier](./sites.md)**
    - [commercetools](./sites.md#commercetools)
    - [contentful](./sites.md#contentful)
    - [amplience](./sites.md#amplience)
    - [azure](./sites.md#azure)
    - [aws](./sites.md#aws)
    - [stores](./sites.md#stores)
    - [components](./sites.md#components)
- [components](./components.md)


!!! tip "JSON schema"
    A JSON schema for the syntax is [available on GitHub](https://github.com/labd/mach-composer/blob/master/schema.json). This can be used to configure IntelliSense autocompletion support in VSCode.

## Full example

!!! TODO
    Add full example

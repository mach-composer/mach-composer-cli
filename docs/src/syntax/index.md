# Configuration syntax

A configuration file can contain several sites with all different configurations and all using a different mix of re-usable serverless microservice components.

It is common to have a single configuration file per environment since they usually share the same general configurations.

The configuration file has the following structure:

- **[general_config](./general_config.md)**
    - **[environment](./general_config.md)**
    - **[terraform_config](./general_config.md#terraform_config)**
    - **[cloud](./general_config.md)**
    - [azure](./general_config.md#azure)
    - [sentry](./general_config.md#sentry)
    - [contentful](./general_config.md#contentful)
    - [amplience](./general_config.md#amplience)
- **[sites](./sites.md)**
    - **[identifier](./sites.md)**
    - [commercetools](./sites.md#commercetools)
    - [contentful](./sites.md#contentful)
    - [amplience](./sites.md#amplience)
    - [azure](./sites.md#azure)
    - [aws](./sites.md#aws)
    - [stores](./sites.md#stores)
    - [components](./sites.md#component-configurations)
- [components](./components.md)


!!! tip "JSON schema"
    A JSON schema for the syntax is [available on GitHub](https://github.com/labd/mach-composer/blob/master/schema.json). This can be used to configure IntelliSense autocompletion support in VSCode.

## Full example

!!! TODO
    Add full example
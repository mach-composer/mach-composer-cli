# Contentful

## Plugin usage

{{ external_markdown('https://raw.githubusercontent.com/mach-composer/mach-composer-plugin-contentful/main/README.md', '## Usage') }}

When configured, MACH can create and manage a [space](https://www.contentful.com/help/spaces-and-organizations/) per site.

For this you need to define an **organization ID** and **CMA token** in your
[general config](../../reference/syntax/global.md#contentful).

## Space configuration

Each site can have their own Contentful Space.

You can define the name of the space that needs to be created in your
[site configuration](../reference/syntax/site.md#dynamic) by giving it a
**name** and optionally set a custom **default locale**.

## Integrate with components

When `contentful` is set as a [component integration](../reference/syntax/site.md#dynamic-1),
the component should have the following Terraform variables defined:

- `contentful_space_id`

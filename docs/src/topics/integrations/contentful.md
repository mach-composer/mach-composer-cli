# Contentful

When configured, MACH can create and manage a [space](https://www.contentful.com/help/spaces-and-organizations/) per site.

For this you need to define a **organization ID** and **CMA token** in your [general config](../../reference/syntax/general_config.md#contentful).

## Space configuration

Each site can have their own Contentful Space.

You can define the name of the space that needs to be created in your [site configuration](../../reference/syntax/sites.md#contentful) by giving it a **name** and optionally set a custom **default locale**.

## Expose Space ID to components

MACH needs to know what components want to use the Space ID.<br>
For this you need to include `contentful` to the list of integrations.<br>
When doing so, MACH expects the component to have one variable `contentful_space_id` defined ([more info](../../reference/components/structure.md#contentful)).
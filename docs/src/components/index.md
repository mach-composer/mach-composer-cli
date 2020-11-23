# What makes a component

A MACH component should be a single repository specifically to provide **domain-specific** functionality.

A single component can consist of one or more of the following things:

- custom commercetools objects
- one or more serverless functions
- a commercetools API extension
- other cloud specific infrastructure

## One component per domain

In order to make your MACH stack truly composable, components should be designed so that it serves to provide functionality for a specific domain within your entire commerce environment.

#### Simple components
For some components, the *domain* and thus the implementation is quite straight-forward and clear, like for example a component responsible for **handling payments** or a component that should **sends a confirmation email** upon a new order.

#### Complex components
Some components might be a bit more complex, but still a good fit for a single component;  
For example an **ERP component** that imports data from an ERP system and exports commercetools order data to be imported in the ERP.  
This might include multiple serverless functions that receive different triggers for example:

- New object is upload to a S3 bucket or Storage account
- New order is received through a commercetools Order

Instead of provided these seperate pieces of functionality in two seperate components, it would make more sense to include that in one component so that when including that component in your MACH stack, all ERP-related functionality is taken care of.

## Creating a component

The easiest way to get started with a component is by using the [`mach bootstrap`](../workflow/cli.md#bootstrap) command:

```bash
$ mach bootstrap component
```

This will generate a new project including a README file on how to get started and include it in your MACH configuration.

## Adding it to your stack

A component can be added to your MACH stack by including it in your MACH configuration.

It should be present in your:

1. [Component definitions](../syntax.md#components) so that MACH knows where to find your component
2. [Site component configuration](../syntax#component-configurations.md) to include it in your MACH stack add site-specific configuration

The [getting started guide](../gettingstarted.md) includes an example of a configuration file with a component implemented.

## Using Serverless framework

The 'function'-part of a MACH component can easily be integrated with the [serverless framework](https://www.serverless.com) for "*zero-friction serverless development*".

This gives you a couple of features:

- Easy local development of your function code including mocked infrastructure that might be needed for your setup
- Build & package your function by calling `sls package`. [More info about Packaging & Deployment](../deployment/components.md##using-serverless)

!!! info "Serverless framework in your MACH deployment"
    Altho we do encourage the usage of the serverless framework for development and packaging, we don't recommend using it for the actual MACH deployment itself.  
    More info about this in the [MACH configuration deployment notes](../deployment/config.md#serverless-framework).

## Further reading

Continue to [component structure](./structure.md) for an explenation of a component's internals.
# Introduction

**MACH composer Cloud** is a platform and API to facilitate and coordinate work
across teams that build composable architectures using MACH technology.

What we've noticed is when the number of teams on a project scales, certain
issues arise that can only be solved by having a centralised API for MACH
composer itself, which acts as a metadata service and provides the ability to
work in a distributed way more effectively. This is why we've started work on
MACH composer cloud.

!!! info "Getting started with MACH composer ☁️" 
    MACH composer ☁️ is an extension to MACH composer and uses GitHub for
    authentication. So if you have [installed the latest
    version](../tutorial/step-1-installation.md#install-mach-composer) of MACH
    composer, you can simply log in or register with MACH composer ☁️ by
    executing this command:

    ``` console
    $ mach-composer cloud login
    ```

    *While in private Beta, we will need to manually enable your account and create 
    or add you to an organisation.*

    [Get started](getting-started.md){ .md-button .md-button--primary }

!!! tip "MACH composer cloud is in private beta"  
    MACH composer ☁️ is very early stage. We try to release new features as often as
    possible and work with these in our MACH projects. So take into account a
    level of instability and backwards incompatibility as long as we're in beta.
    If you would like to know more about it, reach out to us via
    [mach@labdigital.nl](mailto:mach@labdigital.nl?subject=MACH composer cloud)

## Features

Currently, we focus on building features that improve the way of working at
larger scale, with many teams working on the same composable architecture
project.

1. **[Component registry ➡️](component-registry.md)**<br/>
    Increase stability of component deployments, by introducing a centralised
    component registry.
2. **[Autonomous component deployments ➡️](autonomous-deployments.md)** <br/>
    Introduce the ability for components to deploy independently of other
    components, so that teams are less likely to block each other.

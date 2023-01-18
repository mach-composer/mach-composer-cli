!!! warning "MACH composer is undergoing a major, backwards compatible overhaul"
    MACH composer has been rewritten in Golang, coming from a Python version.
    We've done this because the Golang ecosystem is more suitable for a tool
    like this, as it is also the language in which Hashicorp builds Terraform
    itself.

    Because of this, some commands might not be available yet in the Golang version; particularly the
    `mach-composer bootstrap` command is not available, which makes it a bit
    harder to start up. However, based on the [examples in the Github repository](https://github.com/labd/mach-composer/tree/main/examples),
    you should be able to get started without the `bootstrap` command.

    For creating componens, please instead of `mach-composer bootstrap`, use [mach-composer-cookiecutter](https://github.com/labd/mach-component-cookiecutter).

    **We do currently recommend to use the latest version (2.5.x).**

    If you are still on the Python version, the below tutorial should still work.


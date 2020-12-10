# Getting started

## 1. Installation

The easiest way to get started is by installing the MACH composer locally to be able to perform a MACH deploy.

```bash
pipx install mach-composer
```

[More information](./workflow/environment.md#installing-the-cli) about installation.

## 2. Preparations

Before you can deploy a MACH stack, you have to;

1. Setup the basics in your cloud provider, and
2. Create a commercetools project

Instructions can be found in the [prerequisites section](./prerequisites/index.md).

!!! tip
    Keep your created component repository information at hand for later: you'll need it when creating a new component.

    This means; the resource group, storage account and container name in Azure or the AWS S3 bucket name.

## 3. Create a component

Create a MACH component which we can use in our MACH stack later on.

```bash
mach bootstrap component
```

And follow the wizard to create a new component.

This component can now be pushed to a Git repository.

!!! TODO
    Add links to example components on GitHub which can be forked

## 4. Deploy your component

It depends on what component you have, but if you've created a component containing a serverless function, that function needs to be built, packaged and uploaded to a **component registry**.

For a component created with the `mach bootstrap` command or one of the provided example components, these commands will be enough:

```bash
./build.sh package
./build.sh upload
```

!!! info "Component deployment"
    The deployment process of a component can vary.<br>
    [Read more](./deployment/components.md) about component deployments.

## 5. Setup configuration

To create a new MACH configuration file, run

```bash
mach bootstrap config
```

A configuration will be created and can be used as input for the MACH composer.

An example:

=== "AWS"
      ```yaml
      ---
      general_config:
        environment: test
        cloud: aws
        terraform_config:
          aws_remote_state:
            bucket: mach-tfstate-tst
              key_prefix: mach-composer-tst
              region: eu-central-1
      sites:
        - identifier: my-site
          aws:
            account_id: 1234567890
            region: eu-central-1
            route53_zone_name: tst.mach-example.net
          endpoints:
            main: api.tst.mach-example.net
          commercetools:
            project_key: my-site-tst
            client_id: ...
            client_secret: ...
            scopes: manage_project:my-site-tst manage_api_clients:my-site-tst view_api_clients:my-site-tst
            languages:
              - en-GB
              - nl-NL
            currencies:
              - GBP
              - EUR
            countries:
              - GB
              - NL
          components:
            - name: payment
              variables:
                STRIPE_ACCOUNT_ID: 0123456789
              secrets:
                STRIPE_SECRET_KEY: secret-value
      components:
        - name: payment
          source: git::ssh://git@github.com/your-project/components/payment-component.git//terraform
          endpoint: main
          version: e638e57
      ```
=== "Azure"
      ```yaml
      ---
      general_config:
        environment: test
        cloud: azure
        terraform_config:
          azure_remote_state:
            resource_group: my-shared-rg
            storage_account: mysharedsaterra
            container_name: tfstate
            state_folder: test
        azure:
          tenant_id: e180345a-b3e1-421f-b448-672ab50d8502
          subscription_id: 086bd7e7-0755-44ab-a730-7a0b8ad4883f
          region: westeurope
      sites:
        - identifier: my-site
          commercetools:
          project_key: my-site-tst
          client_id: ...
          client_secret: ...
          scopes: manage_project:my-site-tst manage_api_clients:my-site-tst view_api_clients:my-site-tst
          languages:
            - en-GB
            - nl-NL
          currencies:
            - GBP
            - EUR
          countries:
            - GB
            - NL
          components:
            - name: payment
              variables:
                STRIPE_ACCOUNT_ID: 0123456789
              secrets:
                STRIPE_SECRET_KEY: secret-value
      components:
        - name: payment
          source: git::ssh://git@github.com/your-project/components/payment-component.git//terraform
          endpoint: main
          version: e638e57
      ```

See [Syntax](./syntax.md) for all configuration options.

## 6. Deploy

You can deploy your current configuration by running

```bash
$ mach apply
```

If you wish to review the changes before applying them, run

```bash
$ mach plan
```

!!! tip "Using Docker image"
    You can invoke MACH by running the Docker image:<br>
    `$ docker run --rm --volume $(pwd):/code docker.pkg.github.com/labd/mach-composer/mach apply`

    You do need to provide the docker container with the necessary environment variables to be able to authenticate with the cloud provider. More info on that in the [deployment section](./deployment/config/index.md#providing-credentials)


## Example files

You can find example files needed for preparing the infrastructure and a configuration file [on GitHub](https://github.com/labd/mach-composer/tree/master/examples/) in the [/examples](https://github.com/labd/mach-composer/tree/master/examples/) directory

## Further reading

- See the [CLI reference](./workflow/cli.md#apply) for more deployment options.
- Setup your CI/CD pipeline on [GitLab](./deployment/ci/gitlab.md), [GitHub](./deployment/ci/github.md) or [Azure DevOps](./deployment/ci/devops.md)
- [Encrypting your configuration](./security.md#encrypt-your-mach-configuration) with SOPS
- How to create a [new MACH component](./components/index.md)
- [Architectural Guidance](./guidance/index.md)


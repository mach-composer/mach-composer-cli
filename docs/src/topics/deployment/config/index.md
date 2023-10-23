# MACH composer configuration deployment

A MACH composer configuration deployment (or simply put: **MACH deployment**)
will generate and apply a Terraform configuration **per site** so that it can
deploy.

#### 1. MACH-Composer managed resources
The resources that are managed by MACH composer depend on the cloud integration:

- AWS (See [AWS deployments](./aws.md))
- Azure (See [Azure deployments](./azure.md))
#### 2. Integration resources
Resources needed for the integrations such as

- [commercetools](./integrations.md#commercetools)
- [sentry](./integrations.md#sentry)
- [contentful](./integrations.md#contentful)

#### 3. Components
Since components are loaded into the configuration as
[Terraform modules](../../../reference/components/structure.md#terraform-module),
during a MACH composer deployment the resources defined in the component will
get created.

1. The [**first stage**](../components.md) of a component deployment (uploading
   the assets to a component repository) is done before a component is deployed as
   part of a MACH composer stack.

2. The [**second stage**](./components.md) is getting the previously deployed
   component assets actually up and running in your MACH composer stack and to
   create other necessary resources.

More info about the [second stage deployment](./components.md).

!!! info "Component deployment - first and second stage"
    Not all components have a '*first stage*' which means: some components might
    **just** have a Terraform configuration to be applied and no serverless
    *function assets.<br>
    In that case, there is no need of a '*first stage*' component deployment.


## Providing credentials

MACH composer needs to be able to access:

- The components repositories
- The AWS account / Azure subscription it needs to manage resources in

When running MACH composer directly **from the command line**, whenever you have
been authenticated (either by setting the correct AWS environment variables or
on Azure using `az login`) you should be able to deploy using MACH composer
without any issues.

When running the **MACH Docker image**, the necessary environment variables need
to be passed on to the docker container:

=== "AWS"
    ```bash
    docker run --rm \
        --volume $(pwd):/code \
        --volume $SSH_AUTH_SOCK:/ssh-agent \
        -e SSH_AUTH_SOCK=/ssh-agent \
        -e AWS_DEFAULT_REGION=<your-region> \
        -e AWS_ACCESS_KEY_ID=<your-access-key-id> \
        -e AWS_SECRET_ACCESS_KEY=<your-secret-access-key> \
        docker.pkg.github.com/mach-composer/mach-composer-cli/mach:latest \
        apply
    ```
=== "Azure"
    ```bash
    docker run --rm \
        --volume $(pwd):/code \
        --volume $SSH_AUTH_SOCK:/ssh-agent \
        -e SSH_AUTH_SOCK=/ssh-agent \
        -e ARM_CLIENT_ID=<your-client-id> \
        -e ARM_CLIENT_SECRET=<your-client-secret> \
        -e ARM_SUBSCRIPTION_ID=<your-subscription-id> \
        -e ARM_TENANT_ID=<your-tenant-id> \
        docker.pkg.github.com/mach-composer/mach-composer-cli/mach:latest \
        apply --with-sp-login
    ```

    For Azure you'll need to run it with the `--with-sp-login` option let MACH
    Composer perform an `az login` command.<br>
    [More info](../../../reference/cli.md#apply).


## Cache Terraform providers

MACH composer comes with Terraform providers pre-installed in the Docker image.

If you're overwriting these versions with in your
[`terraform_config` block](../../../reference/syntax/global.md#terraform_config),
these providers will be downloaded.

To avoid having to re-download it everytime you run MACH composer through the Docker
image, make sure you mount the [plugin cache](https://www.terraform.io/docs/commands/cli-config.html#provider-plugin-cache)
directory;

```bash
docker run --rm \
    --volume $(pwd):/code \
    --volume $(pwd)/.terraform_cache:/root/.terraform.d/plugin-cache \
    docker.pkg.github.com/mach-composer/mach-composer-cli/mach:latest \
    apply
```

!!! tip "Caching in CI/CD"
    For an example on how to set up the Terraform plugin cache, see the examples in the how-tos for:

    - [GitLab](../../../howto/ci/gitlab.md#terraform-plugin-cache)
    - GitHub Actions (todo)
    - Azure DevOps (todo)

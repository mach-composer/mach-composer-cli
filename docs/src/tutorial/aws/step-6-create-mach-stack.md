# Step 5. Create your MACH stack

To create a new MACH configuration file, run

```bash
mach bootstrap config
```

A configuration will be created and can be used as input for the MACH composer.

An example:

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
    endpoints: 
    main: main
    version: e638e57
```

See [Syntax](../../reference/syntax/index.md) for all configuration options.

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

    You do need to provide the docker container with the necessary environment variables to be able to authenticate with the cloud provider. More info on that in the [deployment section](../../topics/deployment/config/index.md#providing-credentials)


## Example files

You can find example files needed for preparing the infrastructure and a configuration file [on GitHub](https://github.com/labd/mach-composer/tree/master/examples/) in the [/examples](https://github.com/labd/mach-composer/tree/master/examples/) directory

## Further reading

- See the [CLI reference](../../reference/cli.md#apply) for more deployment options.
- Setup your CI/CD pipeline on [GitLab](../../howto/ci/gitlab.md), [GitHub](../../howto/ci/github.md) or [Azure DevOps](../../howto/ci/devops.md)
- [Encrypting your configuration](../../howto/encrypt.md) with SOPS
- How to create a [new MACH component](../../howto/create-component.md)
- [Architectural Guidance](../../topics/architecture/index.md)


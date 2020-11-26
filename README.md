[![MACH composer](https://github.com/labd/mach-composer/blob/master/docs/src/_img/logo.png?raw=true)](https://docs.machcomposer.io)

**Documentation:** [docs.machcomposer.io](https://docs.machcomposer.io)

MACH composer is a framework that you use to orchestrate and extend modern digital commerce & experience platforms, based on MACH technologies and cloud native services. It provides a standards-based, future-proof tool-set and methodology to hand to your teams when building these types of platforms.

It includes:

- A configuration framework for managing MACH-services configuration, using infrastructure-as-code underneath (powered by Terraform)
- A microservices architecture based on modern serverless technology (AWS Lambda and Azure Functions), including (alpha) support for building your microservices with the Serverless Framework
- Multi-tenancy support for managing many instances of your platform, that share the same library of micro services
- CI/CD tools for automating the delivery of your MACH ecosystem
- Tight integration with AWS an Azure, including an (opinionated) setup of these cloud environments

The framework is intended as the 'center piece' of your MACH architecture and incorporates industry best practises such as the 12 Factor Methodology, Infrastrucure-as-code, DevOps, immutable deployments, FAAS, etc.

With combining (and requiring) these practises, using the framework has significant impact on your engineering methodology and organisation. On the other hand, by combining those practises we believe it offers an accelerated 'way in' in terms of embracing modern engineering practises in your organisation.

## Getting started

Read our [getting started guide](https://docs.machcomposer.io/gettingstarted.html) on how to deploy your MACH stack with the MACH composer.

## Example yaml file

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

### Installation

- Make sure you have Terraform 0.13 installed 
tip: Use [tfenv](https://github.com/tfutils/tfenv) to support multiple versions of Terraform)
- Create a new virtualenv with Python 3.8
- Run the following commands:

```
$ make install
``` 

### Running MACH

To generate the files:

`mach generate # generates all available configs.` 

`mach generate -f main.yml`

To plan Terraform:

`mach plan`

To apply Terraform config:

`mach apply`

### Checking for updates

MACH can check your components for available updates.

To do this, run:

`mach update -f main.yml`


## Contributing

### Code style
The Python source code should be formatted using
[black](https://github.com/python/black) and the JavaScript code should be
formatted using [prettier](https://prettier.io/). You can use `make format`
to run both black and prettier.

This project uses [pre-commit](https://pre-commit.com) to validate the changed
files before they are committed. You can install it on MacOS using brew:

    $ brew install pre-commit

In the repository you need to register the hooks in git the first time using:

    $ pre-commit install

The pre-commit config (`.pre-commit-config.yaml`) currently runs black and
flake8.
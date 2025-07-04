[![MACH composer](https://github.com/mach-composer/mach-composer-cli/blob/main/docs/src/_img/logo.png?raw=true)](https://docs.machcomposer.io)

**Documentation:** [docs.machcomposer.io](https://docs.machcomposer.io)

**Plug-ins overview:** [MACH composer plug-ins README](https://github.com/mach-composer#cli-plugins)

MACH composer is a framework that you use to orchestrate and extend modern
digital commerce & experience platforms, based on MACH technologies and cloud
native services. It provides a standards-based, future-proof tool-set and
methodology to hand to your teams when building these types of platforms.

It includes:

- A configuration framework for managing MACH-services configuration, using
  infrastructure-as-code underneath (powered by Terraform)
- A microservices architecture based on modern serverless technology (AWS
  Lambda and Azure Functions), including (alpha) support for building your
  microservices with the Serverless Framework
- Multi-tenancy support for managing many instances of your platform, that
  share the same library of micro services
- CI/CD tools for automating the delivery of your MACH ecosystem
- Tight integration with AWS an Azure, including an (opinionated) setup of
  these cloud environments

The framework is intended as the 'center piece' of your MACH architecture and
incorporates industry best practises such as the 12 Factor Methodology,
Infrastrucure-as-code, DevOps, immutable deployments, FAAS, etc.

With combining (and requiring) these practises, using the framework has
significant impact on your engineering methodology and organisation. On the
other hand, by combining those practises we believe it offers an accelerated
'way in' in terms of embracing modern engineering practises in your
organisation.

## Installation

### Shell script

Run the following command to install MACH composer on your system. The script will download the latest release of MACH composer and install it in your
`$HOME/bin` directory. Make sure that this directory is in your `$PATH`.

```
$ curl -sfL https://raw.githubusercontent.com/mach-composer/mach-composer-cli/f08424b1bc38086696767a1ce05e1b0fbb199326/scripts/install-mach-composer.sh | bash
```

If you want to install a specific version, you can specify the `VERSION` environment variable:

```
$ export VERSION=v2.20.0
$ curl -sfL https://raw.githubusercontent.com/mach-composer/mach-composer-cli/refs/heads/main/scripts/install-mach-composer.sh | bash
```

### Brew

```bash
brew tap mach-composer/mach-composer
brew install mach-composer
```

### Chocolatey (Windows)

Windows installation through Chocolatery is currently unstable. We recommend to [download the latest release from GitHub
Releases](https://github.com/mach-composer/mach-composer-cli/releases/latest). Also, it is recommended to run MACH composer through [WSL](https://learn.microsoft.com/en-us/windows/wsl/install).

<!--```ps
choco install mach-composer --version=2.5.0
``` -->

### Nix

Add the flake input

```
  inputs.mach-composer = {
    url = "github:mach-composer/nix-mach-composer";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  # in outputs, pass mach-composer into your configuration
```

Then add to your package list

```
  mach-composer.packages.${system}.mach-composer
```

#### making an overlay
```
let
  mach-composer-overlay = final: prev: {
    mach-composer = mach-composer.packages.${system}.mach-composer;
  };
  pkgs = import nixpkgs {
    inherit system;
    config = {
      allowUnfree = true;
    };

    overlays = [
      mach-composer-overlay
    ];
  };
```

> if using [numtide/devshell](https://github.com/numtide/devshell/), you can then put `mach-composer` in your `devshell.toml` packages list.

## Getting started

Read our [getting started guide](https://docs.machcomposer.io/gettingstarted.html)
on how to deploy your MACH stack with MACH composer.

## Example yaml file

```yaml
---
mach_composer:
  version: 1
global:
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
      public: api.tst.mach-example.net
    commercetools:
      project_key: my-site-tst
      client_id: ...
      client_secret: ...
      scopes: manage_project:my-site-tst manage_api_clients:my-site-tst view_api_clients:my-site-tst
      token_url: https://auth.europe-west1.gcp.commercetools.com
      api_url: https://api.europe-west1.gcp.commercetools.com
      project_settings:
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
      main: public
    version: e638e57
```

### Running MACH

To generate the files:

```console
mach-composer generate # generates config for main.yml
mach-composer generate -f other-file.yml
```

To plan Terraform:

```console
mach-composer plan
```

To apply Terraform config:

```console
mach-composer apply
```

Optionally you can run a terraform init without taking any action:

```console
mach-composer terraform init
```


### Checking for updates

MACH can check your components for available updates.

To do this, run:

```console
mach-composer update -f main.yml
```

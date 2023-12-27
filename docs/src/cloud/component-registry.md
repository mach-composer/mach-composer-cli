# Component registry

The MACH composer cloud **component registry** provides an API where you can
register and retrieve components and their versions, so that you can guarantee
that a component-version and its accompanying artifact, actually exist. This
will make your MACH composer deployment pipelines more stable.

## Usage

Pushing a new component version can be done by executing the below command from
the git clone of your component (or from within a CI/CD pipeline).

```bash
$ mach-composer cloud register-component-version
--organization $MCC_ORGANISATION --project $MCC_PROJECT $MCC_COMPONENT --auto
``` 

!!! warning "Only register a version, if there is an artifact."
    Pushing a new version should only be done after an artifact was
    actually created and published to your artifact repository.

If your using GitHub Actions, we have created a [convenient GitHub Action for
registering
versions](https://github.com/mach-composer/new-component-version-action), that
you can simply include in an existing component pipeline.

## How does it work?

We introduce an additional step to the [component deployment
sequence](../concepts/components/lifecycle/index.md), by connecting that sequence to MACH
composer cloud.

<div class="mermaid">
sequenceDiagram
    participant D as Developer
    participant S as SCM
    participant CI as CI/CD
    participant C as Artifact Repository
    participant MCC as MACH Composer☁️ 

    D->>S: Merges new code
    S->>CI: run build and tests
    opt if master branch
        CI->>C: Push artifact to artifact repository
    end
    opt if using MACH composer ☁️
        CI->>MCC: push component version to MACH composer ☁️
    end
</div>

Then when running the `mach-composer update` command, instead of fetching the
latest component versions, MACH composer will fetch the latest versions from
MACH composer cloud.

<div class="mermaid">
sequenceDiagram
    participant Dev as Developer
    participant MC as MACH composer CLI
    participant MCC as MACH Composer☁️ 
    participant GIT as Git repository
    Dev->>MC: runs update
    alt using MACH composer ☁️
        MC->>MCC: fetch component version
    else using git
        MC->>GIT: fetch latest version from Git
    end
    MC->>Dev: updates versions in YAML
</div>

## Why do we need a registry?

MACH composer is in the business of deploying components as part of a coherent
composable architecture. It provides tools to coordinate this and make this
easier.

Components are in fact simple Terraform modules that are usually containing a
microservice that's deployed in for example AWS Lambda or Fargate.

In your MACH composer configuration you configure the components that you use in
you project, as well as the version of components. Just like the example below.

```yaml
components:
  - name: api-service
    source: git::https://github.com/mach-composer/mcc-api-service.git//terraform
    version: "66e74f4"
    branch: main

  - name: auth-service
    source: git::https://github.com/mach-composer/mcc-auth-service.git//terraform
    version: "7e746bb"
    branch: main
```

Then, using the `mach-composer update` command, MACH composer will check the
source repository for new versions that might be available.

```console
$ mach-composer update

Updates for auth-service (19ffb8b...7e746bb)
7e746bb: Only add *.go source files to archive (Michael van Tellingen michaelvantellingen@gmail.com)
92f14ba: Remove unused function (Michael van Tellingen michaelvantellingen@gmail.com)

Updates for api-service (0771f35...66e74f4)
66e74f4: Only add *.go files to artifact zip (Michael van Tellingen michaelvantellingen@gmail.com)
```

### The problem

After updating and running `mach-composer apply`, usually some Terraform code
will naively take this version information and use it to deploy either a Lambda
or Docker container by building an artifact name like this:
`my-component-{version}.zip`.

Usually this works fine, but it occasionally happens that the zip file does not
exist, in case a CI/CD build has failed to produce the artifact. This will cause
the `mach-composer apply` execution to fail - sometimes silently.

### The solution

Our solution to this is to provide a way to register 'deployable versions' with
MACH composer cloud, as part of the component CI/CD workflow: **after
successfully building and storing the component artifact, it must be registered
with the MACH composer API**.

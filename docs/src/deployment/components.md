# Component deployment

As described in the [components section](../components/index.md#deployment-process), the component itself is responsible for the following steps:

- Packaging the function
- Deploying it to the code registry

In this section a couple of methods – for various cloud providers – will be described.

## On Azure

### Package
=== "Python"
    ```bash
    VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    NAME=yourcomponent-$VERSION
    ARTIFACT_NAME=$NAME.zip

    func pack --build-native-deps --python
    mv $BASENAME.zip $ARTIFACT_NAME
    ```

=== "Node"
    ```bash
    TODO
    ```

=== ".NET"
    ```bash
    TODO
    ```

### Upload
```bash
STORAGE_ACCOUNT_KEY=`az storage account keys list -g my-shared-we-rg -n mysharedwesacomponents --query [0].value -o tsv`
az storage blob upload --account-name mysharedwesacomponents --account-key $STORAGE_ACCOUNT_KEY -c code -f yourcomponent-0.1.0.zip -n yourcomponent-0.1.0.zip
```

!!! tip ""
    An example build and deploy script is provided in the [component cookiecutter](https://git.labdigital.nl/mach/component-cookiecutter)


## On AWS

### Package

=== "Python"
    ```bash
    VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    NAME=yourcomponent-$VERSION
    ARTIFACT_NAME=$NAME.zip

    python -m pip install dist/*.whl -t ./build
    cp handler.py ./build
    cd build && zip -9 -r $ARTIFACT_NAME .
    ```

### Upload
```bash
aws s3 cp build/$ARTIFACT_NAME s3://your-lambda-repo/
```
## Using serverless

TODO

## Setting up CI

Refer to the CI/CD section for instructions on how to setup your Continuous Integration pipeline for component deployments:

- [GitLab](./ci/gitlab.md#components)
- [Azure DevOps](./ci/devops.md#components)
- [Jenkins](./ci/jenkins.md#components)
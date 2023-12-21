# Component deployment - first stage

As described in the [components section](../../concepts/components/index.md)
whenever a component contains a **serverless function** that needs to be
uploaded to a registry, it needs a '*first stage*' deployment.

In that case, the component itself is responsible for the following steps:

- Packaging the function
- Deploying it to the code registry

In this section a couple of methods – for various cloud providers – will be described.

!!! tip "Second stage deployment"
    The second stage deployment is described in the [MACH configuration deployment section](./second-stage.md)

## Package

### Using serverless

When using [serverless](https://www.serverless.com) for local development, the
easiest way to package your function would be to invoke

```bash
$ serverless package
```

### Python project

Sample bash script to package your function to a ZIP file with the correct version tag.

=== "AWS"
    ```bash
    VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    NAME=your-component-$VERSION
    ARTIFACT_NAME=$NAME.zip

    python -m pip install dist/*.whl -t ./build
    cp handler.py ./build
    cd build && zip -9 -r $ARTIFACT_NAME .
    ```
=== "Azure"
    ```bash
    VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    NAME=your-component-$VERSION
    ARTIFACT_NAME=$NAME.zip

    func pack --build-native-deps --python
    mv $BASENAME.zip $ARTIFACT_NAME
    ```


## Package & Upload script

An example of a bash script that can be used in any CI/CD pipeline.
The content of the `upload` and `package` functions can be taken from the examples mentioned above.

=== "Azure"
    ```bash
    #!/bin/bash

    VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    BASENAME=component-name
    NAME=$BASENAME-$VERSION
    BUILD_NAME=$NAME
    ARTIFACT_NAME="${NAME}.zip"

    package () {
        # Package instructions
    }

    upload () {
        src="build/${ARTIFACT_NAME}"
        dest=$ARTIFACT_NAME
        STORAGE_ACCOUNT_KEY=`az storage account keys list -g my-shared-we-rg -n mysharedwesacomponents --query [0].value -o tsv`
        az storage blob upload --account-name mysharedwesacomponents --account-key ${STORAGE_ACCOUNT_KEY} -c code -f $src -c code -n $dest
    }

    version () {
        echo "Version: '${VERSION}'"
    	echo "Name: '${NAME}'"
    	echo "Artifect name: '${ARTIFACT_NAME}'"
    }

    case $1 in
        package)
            package $2 $3
        ;;
        upload)
            upload $2
        ;;
        version)
            version
        ;;
    esac
    ```
=== "AWS"
    ```bash
    #!/bin/bash

    VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    BASENAME=component-name
    NAME=$BASENAME-$VERSION
    BUILD_NAME=$NAME
    ARTIFACT_NAME="${NAME}.zip"

    package () {
        # Package instructions
    }

    upload () {
        src="build/${ARTIFACT_NAME}"
        dest=$ARTIFACT_NAME
        aws s3 cp $src s3://your-lambda-repo/$dest
    }

    version () {
        echo "Version: '${VERSION}'"
    	echo "Name: '${NAME}'"
    	echo "Artifect name: '${ARTIFACT_NAME}'"
    }

    case $1 in
        package)
            package $2 $3
        ;;
        upload)
            upload $2
        ;;
        version)
            version
        ;;
    esac
    ```

**Usage:**
```bash
$ ./build.sh package
$ ./build.sh upload
```
## Setting up CI

Refer to the CI/CD section for instructions on how to set up your Continuous
Integration pipeline for component deployments:

- [GitLab](../../howto/ci/gitlab.md#component-deployment)
- [GitHub Actions](../../howto/ci/github.md#component-deployment)
- [Azure DevOps](../../howto/ci/azure_devops.md#component-deployment)


## Second stage - MACH deployment

As part of the [MACH deployment](./index.md), the component will be
loaded as a Terraform module so that the component can deploy the necessary
resources needed for the specific site it is configured for.

This is what we call the [second stage deployment](./second-stage.md).

###

# Component deployment

As described in the [components section](../components/structure.md#deployment-process), the component itself is responsible for the following steps:

- Packaging the function
- Deploying it to the code registry

In this section a couple of methods – for various cloud providers – will be described.

## Package

### Using serverless

When using [serverless](https://www.serverless.com) for local development, the easiest way to package your function would be to invoke

```bash
$ serverless package
```

### Python project

Sample bash script to package your function to a ZIP file with the correct version tag.

=== "AWS"
    ```bash
    VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    NAME=yourcomponent-$VERSION
    ARTIFACT_NAME=$NAME.zip

    python -m pip install dist/*.whl -t ./build
    cp handler.py ./build
    cd build && zip -9 -r $ARTIFACT_NAME .
    ```
=== "Azure"
    ```bash
    VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
    NAME=yourcomponent-$VERSION
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

Refer to the CI/CD section for instructions on how to setup your Continuous Integration pipeline for component deployments:

- [GitLab](./ci/gitlab.md#components)
- [GitHub actions](./ci/github.md#components)
- [Azure DevOps](./ci/devops.md#components)
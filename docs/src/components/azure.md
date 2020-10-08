# Azure components

## Packaging and deploying

For Azure functions, the deployment process constist of two steps:

- Packaging the function
- Deploying it to the [function app storage](../prerequisites/azure.md#create-function-app-storage)

### Package
```bash
VERSION=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev" )
NAME=yourcomponent-$VERSION
ARTIFACT_NAME=$NAME.zip

func pack --build-native-deps --python
mv $BASENAME.zip $ARTIFACT_NAME
```

### Upload
```bash
az storage blob upload --account-name mysharedwesacomponents --account-key $STORAGE_ACCOUNT_KEY -c code -f yourcomponent-0.1.0.zip -n yourcomponent-0.1.0.zip
```

!!! tip ""
    An example build and deploy script is provided in the [component cookiecutter](https://git.labdigital.nl/mach/component-cookiecutter)

## HTTP routing

MACH will provide the correct HTTP routing for you.  
To do so, the following has to be configured:

- [Frontdoor](../syntax.md#front_door) settings in the Azure configuration
- The component has to be marked as [`has_public_api`](../syntax.md#components)

More information in the [deployment section](../deployment/azure.md#http-routing).

# Azure dashboard configuration

!!! Todo
    Future implementation
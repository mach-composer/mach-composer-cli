# Getting started with MACH composer ☁️ in an existing MACH composer project

MACH composer ☁️ is an extension to MACH composer and can be used in existing
MACH composer projects. It's main entrypoint is therefore still the CLI, which
has been expanded with a `cloud` subsection with a [number of
actions](../reference/cli/mach-composer_cloud.md).

!!! info "The below steps can be achieved through the dashboard" 
    While MACH composer ☁️ is an API-first product, we also provide a dashboard that
    you can use to manage your projects and components. All the actions above
     can also be achieved through it.

    You can login to the dashboard here: [https://app.mach.cloud](https://app.mach.cloud)


In order to start using MACH composer ☁️, we recommend to follow the below steps
using the CLI. These steps describe how to start registering new components with
MACH composer ☁️.

## 1. Create a MACH composer ☁️ account, organization and project

1. **Create a MACH composer cloud account by logging in with GitHub**

    ```bash
    # this will open a browser in which you can login through a GitHub account
    % mach-composer cloud login
    Successfully authenticated to mach composer cloud
    ```

    !!! warning "Currently we're in private beta and need to manually approve new accounts."


1. **Create an organization**
    ```bash
    % mach-composer cloud create-organization --key my-org --name "My org name"
    Created new organization: my-org
    ```

    !!! info "You can be part of multiple organizations and also invite other users to your organization by using the `mach-composer cloud add-organization-user` command."

2. **Create a project**
    ```bash
    % mach-composer cloud create-project my-project "My Project" --organization my-org
    Created new project: my-project
    ```


## 2. Register components in MACH composer ☁️ through CI/CD

Now that you've set up your account, it is now time to start registering
components with MACH composer ☁️. 

1. **Create components in MACH composer ☁️**

    For each of the components you have, execute the below command. Each new
    component should be registered like this with MACH composer ☁️.

    ```bash
    % mach-composer cloud create-component my-comp --organization my-org --project my-project
    Created new component: my-comp
    ```

2. **Now you can push new versions of a component.**

    Run the following command from a GitHub repository.

    ```bash
    $ mach-composer cloud register-component-version
    --organization my-org --project my-project --auto
    ``` 

    Using `--auto` will detect the latest version from git and will also read
    metadata (i.e. the commits that are part of the new version) and push that
    to MACH composer ☁️. For more info about using the component registry, [take
    a look at its documentation](component-registry.md).

3. **In CI/CD context you need to create API Client credentials**
    
    When you are logged in from the commandline the MACH composer CLI will
    manage your API credentials for you through GitHub authentication. but in
    CI/CD context, you need to create a set of API Client credentials to make it
    work. `mach-composer` will detect the `MCC_CLIENT_ID` and
    `MCC_CLIENT_SECRET` environment variables. 

    You can create those like this:

    ```bash
    $ mach-composer cloud create-api-client --organization my-org --project my-project
    Client ID: LWAxfy1gCwS5kfXXXXXXXXXXXXXX
    Client Secret: 1IFau1MhyETzPzZVQ5Vp3Uw6hMOfFOyXXXXXXXXXXXXXXXXXX
    Scopes: project:manage
    ```

    !!! info "You can list existing API clients (without secret)"

        ```bash
        $ mach-composer cloud list-api-clients --organization my-org --project my-project
        --------------------------------------------------------------------------------------------------------------
        CREATED AT         	CLIENT ID                       	CLIENT SECRET	LAST USED	DESCRIPTION	SCOPES        	
        --------------------------------------------------------------------------------------------------------------
        2023-01-13 17:52:06	LWAxfy1gCwS5kf1bjad2HqOofOEfKo0n	*****uNOOl   	never    	           	project:manage	
        --------------------------------------------------------------------------------------------------------------
        ```

## 3. Integrate MACH composer ☁️ in your MACH composer configuration

!!! tip "This feature is available since MACH composer 2.7.x"

The last step is to integrate MACH composer ☁️ in your existing MACH composer
YAML configuration. You can achieve this by adding the following configuration
to the `mach_composer` section.

```yaml
mach_composer:
    organization_id: "lab-digital"
    project_id: "my-aws-example"
```

When running `mach-composer update --cloud`, MACH composer will fetch the latest
versions from MACH composer ☁️ instead of the components GitHub repository. As a
result, this process will be much faster as well as more robust, as component
versions are guaranteed to exist.

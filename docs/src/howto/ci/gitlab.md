# GitLab CI/CD

This section will describe how to setup your CI/CD pipeline using GitLab including some examples.

## MACH stack

How to set up the deployment process for your MACH configuration.

### Providing credentials

For an deployment we have to make sure the following variables set in the GitLab CI/CD settings;

- Personal access token for the MACH docker image
- Azure or AWS credentials

These will be explained further:

- [Access to MACH docker image](#mach-docker-image)
- [Access to component repositories](#component-repositories)
  
![CI/CD variables](../../_img/deployment/gitlab/variables.png)

#### MACH docker image

Since the MACH docker image is hosted on a private GitHub repository, we need to add credentials to be able to pull the image.

1. Create a [personal access token](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)<br>
   Make sure this as the `read:packages` permission
2. Determine your [`DOCKER_AUTH_CONFIG` data](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html#determining-your-docker_auth_config-data)
 <br>
Following the [example from GitLab](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html#determining-your-docker_auth_config-data), you can create the auth-token as such:

```bash
# The use of "-n" - prevents encoding a newline in the password.
echo -n "my_username:my_access_token" | base64

# Example output to copy
bXlfdXNlcm5hbWU6bXlfcGFzc3dvcmQ=
```

And then give `DOCKER_AUTH_CONFIG` the following value:

```
{
    "auths": {
        "registry.example.com:5000": {
            "auth": "(Base64 content from above)"
        }
    }
}
```


#### Component repositories

When MACH is applied it will have to download the various components from their Git repositories.<br>
We have to make sure the current runner has access to those.

Most probably you'll have the CI for the MACH configuration running under the same GitLab account as the components itself.<br>
In that case you can use the [`CI_JOB_TOKEN`](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html) variable and place it in a [`.netrc`](https://docs.gitlab.com/ee/user/project/new_ci_build_permissions_model.html#dependent-repositories) file so that other repositories can be accessed (see [example](#example)).

### Example

=== "AWS"
      ```yaml
      ---
      image: docker.pkg.github.com/labd/mach-composer/mach:0.4

      variables:
        AWS_DEFAULT_REGION: $AWS_DEFAULT_REGION
        AWS_ACCESS_KEY_ID: $AWS_ACCESS_KEY_ID
        AWS_SECRET_ACCESS_KEY: $AWS_SECRET_ACCESS_KEY

      before_script:
        - mkdir -p ~/.ssh
        - chmod 700 ~/.ssh
        # Replace with your custom GitLab domain<br>
        - ssh-keyscan your.gitlab-domain.com >> ~/.ssh/known_hosts
        - chmod 644 ~/.ssh/known_hosts
        - echo -e "machine your.gitlab-domain.com\nlogin gitlab-ci-token\npassword ${CI_JOB_TOKEN}" > ~/.netrc

      deploy:
        script:
          - mach apply --auto-approve -f $CI_PROJECT_DIR/main.yml
      ```
=== "Azure"
      ```yaml
      ---
      image: docker.pkg.github.com/labd/mach-composer/mach:0.4

      variables:
        ARM_CLIENT_ID: $AZURE_SP_CLIENT_ID
        ARM_CLIENT_SECRET: $AZURE_SP_CLIENT_SECRET
        ARM_SUBSCRIPTION_ID: $AZURE_SP_SUBSCRIPTION_ID
        ARM_TENANT_ID: $AZURE_SP_TENANT_ID

      before_script:
        - mkdir -p ~/.ssh
        - chmod 700 ~/.ssh
        # Replace with your custom GitLab domain<br>
        - ssh-keyscan your.gitlab-domain.com >> ~/.ssh/known_hosts
        - chmod 644 ~/.ssh/known_hosts
        - echo -e "machine your.gitlab-domain.com\nlogin gitlab-ci-token\npassword ${CI_JOB_TOKEN}" > ~/.netrc

      deploy:
        script:
          - mach apply --auto-approve --with-sp-login -f $CI_PROJECT_DIR/main.yml
      ```

## Components

Example GitLab CI configuration

=== "Python"
  ```yaml
  stages:
    - test
    - build
    - deploy

  image: mcr.microsoft.com/azure-functions/python:3.0-python3.8-core-tools

  variables:
    PIP_CACHE_DIR: "$CI_PROJECT_DIR/pip-cache"

  cache:
    paths:
      - "$CI_PROJECT_DIR/pip-cache"
    key: "$CI_PROJECT_ID"

  test:
    image: python:3.7.5
    stage: test
    script: 
      - pip install -r requirements_dev.txt
      - py.test tests/ --cov=. --cov-report=term-missing --cov-report=xml:reports/coverage.xml --junit-xml=reports/junit.xml
    artifacts:
      reports:
        junit: reports/junit.xml
        cobertura: reports/coverage.xml
      paths:
        - reports/junit.xml
        - reports/coverage.xml

  build:
    stage: build
    script: 
      - make pack
    artifacts:
      paths:
        - build/*.zip
      expire_in: 1 day

  deploy:
    stage: deploy
    script: 
      - az login --service-principal -u $AZURE_SP_CLIENT_ID -p $AZURE_SP_CLIENT_SECRET --tenant $AZURE_SP_TENANT_ID
      - az account set --subscription $AZURE_SP_SUBSCRIPTION_ID
      - make upload
  ```
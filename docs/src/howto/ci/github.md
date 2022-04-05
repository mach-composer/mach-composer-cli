# Deploy using GitHub Actions

This section will describe how to setup your CI/CD pipeline using GitHub Actions
including some examples.

## MACH stack deployment

How to set up the deployment process for your MACH configuration.

### Providing credentials

For an deployment we have to make sure the following variables set in the GitLab
CI/CD settings;

- [Personal access token](#create-access-token)
- AWS or Azure credentials

#### Create access token
1. Create a [personal access token](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)<br>
   Make sure this has the `repo` permission
2. Set the personal access token credentials as secrets in your MACH
   configuration repo settings.

!!! note "Permissions needed"
      We need `repo` to have access to any private repositories so that MACH can pull in the components during deployment.

### Example
=== "AWS"
      ```yaml
      name: MACH rollout

      on:
        push:
          branches:
            - master

      jobs:
        mach:
          runs-on: ubuntu-latest
          container:
            image: docker.pkg.github.com/labd/mach-composer/mach:2.0.0
            credentials:
              username: ${{ secrets.GITHUB_USER }}
              password: ${{ secrets.GITHUB_TOKEN }}
          steps:
          - uses: actions/checkout@v2
          - run: |
              echo -e "machine github.com\nlogin ${{ secrets.GITHUB_USER }}\npassword ${{ secrets.GITHUB_TOKEN }}" > ~/.netrc
            name: Prepare credentials
          - run: mach-composer apply --auto-approve
            name: Apply
            env:
              AWS_DEFAULT_REGION: eu-central-1
              AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
              AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
      ```
=== "Azure"
      ```yaml
      name: MACH rollout

      on:
        push:
          branches:
            - master

      jobs:
        mach:
          runs-on: ubuntu-latest
          container:
            image: docker.pkg.github.com/labd/mach-composer/mach:2.0.0
            credentials:
              username: ${{ secrets.GITHUB_USER }}
              password: ${{ secrets.GITHUB_TOKEN }}
          steps:
          - uses: actions/checkout@v2
          - run: |
              echo -e "machine github.com\nlogin ${{ secrets.GITHUB_USER }}\npassword ${{ secrets.GITHUB_TOKEN }}" > ~/.netrc
            name: Prepare credentials
          - run: mach-composer apply --auto-approve --with-sp-login
            name: Apply
            env:
              ARM_CLIENT_ID: ${{ secrets.ARM_CLIENT_ID }}
              ARM_CLIENT_SECRET: ${{ secrets.ARM_CLIENT_SECRET }}
              ARM_SUBSCRIPTION_ID: ${{ secrets.ARM_SUBSCRIPTION_ID }}
              ARM_TENANT_ID: ${{ secrets.ARM_TENANT_ID }}
      ```

## Component deployment

For the component CI pipeline we need to be able to test, package and upload the
function app ZIP file.

### Example

Example GitHub action to package and deploy a component on AWS.

=== "Node"
  ```yaml
   name: Package and upload

   on:
     push:
       branches:
         - main

   env:
     PACKAGE_NAME: my-component
     AWS_BUCKET_NAME: my-lambda-bucket

   jobs:
     package:
       name: Package Lambda function
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v2

         - name: Artifact Name
           id: artifact-name
           run: echo "::set-output name=artifact::$(echo $PACKAGE_NAME-${GITHUB_SHA:0:7}.zip)"

         - name: Use Node.js
           uses: actions/setup-node@v1
           with:
             node-version: 12.x

         - name: Cache modules
           uses: actions/cache@v2
           with:
             path: '**/node_modules'
             key: ${{ runner.os }}-modules-${{ hashFiles('**/yarn.lock') }}

         - name: Install dependencies
           run: yarn

         - name: Package
           uses: dragonraid/sls-action@v1.2
           with:
             args: --stage prod package

         - name: Configure AWS Credentials
           uses: aws-actions/configure-aws-credentials@v1
           with:
             aws-access-key-id: ${{ secrets.MACH_ARTIFACT_AWS_ACCESS_KEY_ID }}
             aws-secret-access-key: ${{ secrets.MACH_ARTIFACT_AWS_SECRET_ACCESS_KEY }}
             aws-region: eu-central-1

         - name: Upload
           run: aws s3 cp .serverless/${{ env.PACKAGE_NAME }}.zip s3://${{ env.AWS_BUCKET_NAME }}/${{ steps.artifact-name.outputs.artifact }}

  ```

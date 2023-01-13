# Deploy using GitHub Actions
This section will describe how to setup your CI/CD pipeline using GitHub Actions
including some examples.

## MACH stack deployment
This page describes how to deploy your MACH composer configuration to your
cloud provider. We recommend to use a PR workflow whereby changes are done via
a PR which allows additional rules (e.g. approval flows).


### Authenticating with your cloud provider
For deploying to a cloud provider (e.g. AWS, Azure or GCP) as we recommend to
use OpenID Connect. See how this works with GitHub at
[security hardening your deployments](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments)


### Example

=== "AWS"

    ```yaml
    name: Deploy Test
    on:
      pull_request:
        types: [opened, synchronize, closed]
        paths:
          - 'main.yml'

    concurrency: test

    env:
      TF_PLUGIN_CACHE_DIR: ${{ github.workspace }}/.terraform.d/plugin-cache
      MACH_COMPOSER_VERSION: 2.5.4
      AWS_DEPLOY_ROLE: arn:aws:iam::<AWS_ACCOUNT_ID>:role/mach-deploy-role
      AWS_DEPLOY_REGION: "eu-west-1"
      AWS_PLAN_BUCKET: "bucket name to store terraform plans"
      CONFIG_FILE: "main.yml"

    jobs:
      plan:
        runs-on: ubuntu-latest
        if: github.event.pull_request.merged != true && github.base_ref == 'main'
        permissions:
          id-token: write # This is required for requesting the JWT
          contents: read  # This is required for actions/checkout
          pull-requests: write
        steps:
          - uses: actions/checkout@v3

          - name: Prepare credentials
            run: |
              echo ${{ secrets.ORG_REPO_TOKEN }}
              git config --global url."https://oauth2:${{ secrets.ORG_REPO_TOKEN }}@github.com".insteadOf https://github.com

          - name: Install terraform
            uses: hashicorp/setup-terraform@v2
            with:
              terraform_version: 1.1.7

          - name: Create Terraform Plugins Cache Dir
            run: mkdir --parents $TF_PLUGIN_CACHE_DIR

          - name: Cache Terraform Plugins
            uses: actions/cache@v2
            with:
              path: ${{ env.TF_PLUGIN_CACHE_DIR }}
              key: ${{ runner.os }}-terraform-${{ hashFiles('**/.terraform.lock.hcl') }}

          - name: Install MACH composer
            uses: mach-composer/setup-mach-composer@main
            with:
              version: ${{ env.MACH_COMPOSER_VERSION }}

          - name: Install sops
            uses: mdgreenwald/mozilla-sops-action@v1.4.1

          - name: Configure AWS Credentials
            uses: aws-actions/configure-aws-credentials@v1
            with:
              role-to-assume: ${{ env.AWS_DEPLOY_ROLE }}
              aws-region: ${{ env.AWS_DEPLOY_REGION }}

          - name: MACH composer plan
            uses: mach-composer/plan-action@main
            with:
              filename: ${{ env.CONFIG_FILE }}
              github-token: ${{ secrets.GITHUB_TOKEN }}

          - name: Store terraform plan
            run: aws s3 cp --sse AES256 --recursive --exclude '*' --include "deployments/*/terraform.plan" . s3://${{ env.AWS_PLAN_BUCKET }}/${{ github.event.pull_request.number }}/


      deploy:
        runs-on: ubuntu-latest
        if: github.event.pull_request.merged == true && github.base_ref == 'main'
        environment:
          name: test
          url: "<your env url>"
        permissions:
          id-token: write # This is required for requesting the JWT
          contents: read  # This is required for actions/checkout
        steps:
          - uses: actions/checkout@v3

          - name: Prepare credentials
            run: |
              echo ${{ secrets.ORG_REPO_TOKEN }}
              git config --global url."https://oauth2:${{ secrets.ORG_REPO_TOKEN }}@github.com".insteadOf https://github.com

          - name: Install terraform
            uses: hashicorp/setup-terraform@v2
            with:
              terraform_version: 1.1.7

          - name: Create Terraform Plugins Cache Dir
            run: mkdir --parents $TF_PLUGIN_CACHE_DIR

          - name: Cache Terraform Plugins
            uses: actions/cache@v2
            with:
              path: ${{ env.TF_PLUGIN_CACHE_DIR }}
              key: ${{ runner.os }}-terraform-${{ hashFiles('**/.terraform.lock.hcl') }}

          - name: Install MACH Composer
            uses: mach-composer/setup-mach-composer@main
            with:
              version: ${{ env.MACH_COMPOSER_VERSION }}

          - name: Install sops
            uses: mdgreenwald/mozilla-sops-action@v1.4.1

          - name: Configure AWS Credentials
            uses: aws-actions/configure-aws-credentials@v1
            with:
              role-to-assume: ${{ env.AWS_DEPLOY_ROLE }}
              aws-region: ${{ env.AWS_DEPLOY_REGION }}

          - name: Retrieve terraform plan
            run: aws s3 cp --sse AES256 --recursive s3://${{ env.AWS_PLAN_BUCKET }}/${{ github.event.pull_request.number }}/ .

          - name: Run MACH Composer apply
            run: mach-composer apply --auto-approve -f ${{ env.CONFIG_FILE }}

          - name: Cleanup terraform plan
            run: aws s3 rm --recursive s3://${{ env.AWS_PLAN_BUCKET }}/${{ github.event.pull_request.number }}/
    ```

=== "Azure"

    ```yaml
    name: MACH composer rollout

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

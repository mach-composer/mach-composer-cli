# Step 5. Create and deploy component

{%
   include-markdown "../../overhaul_warning.md"
%}

## 1. Create a MACH component

Create a MACH component which we can use in our MACH stack later on.

```bash
$ mach-composer bootstrap component
```

And follow the wizard to create a new component.

!!! note "Simple API component"
    We will be creating a simple API component that exposes 1 endpoint.

    The bootstrapper will ask for some input to create the component, for the
    sake of this tutorial let's use the following values and replace
    `your-project-lambdas` with the S3 bucket name created in [step 3](step-3-setup-aws-services.md).

    ```
    Cloud environment (aws, azure) [aws]: aws
    Language (python, node) [node]: node
    Name [example-name]: api
    Description [Api component]:
    Directory name [api-component]:
    Uses an HTTP endpoint? [Y/n]: y
    Include GraphQL support? [Y/n]: y
    Uses commercetools? [Y/n]: n
    Use Sentry? [y/N]: n
    Lambda repository S3 bucket: your-project-lambdas
    New component api-component created 🎉
    ```

## 2. Initiate and test component

In the previous step we've created a Node component.

To install the environment and test it locally, run the following commands:

```bash
$ yarn
$ yarn test
```

Now that we've verified that the component works, we can deploy it to our component repository.
## 3. Deploy your component

It depends on what component you have, but if you've created a component
containing a serverless function, that function needs to be built, packaged and
uploaded to a **component registry**.

### Package

For a component created with the `mach bootstrap` command, the component can be packaged using;

```bash
$ ./build.sh package
```

### Upload
Now that we have our packaged component, we need to upload this to the component
repository that we [created in step 3](./step-3-setup-aws-services.md).

Make sure that you're using the correct AWS profile to upload the ZIP file.

We can use the same role that we've used to perform the infra rollout. So make
sure the following is added to your `~/.aws/config` file:

```conf
[profile your-project-srv]
source_profile = default
role_arn = arn:aws:iam::<service-account-id>:role/admin
```

Using this AWS profile, we can upload the component.

```bash
$ ./build.sh upload
```

!!! tip "Providing AWS credentials"
    To ensure the correct AWS profile is used you can set the `AWS_DEFAULT_PROFILE` environment variable.

    The command will be `AWS_DEFAULT_PROFILE=your-project-srv ./build.sh upload`.

    We recommend using a solution like [aws-vault](https://github.com/99designs/aws-vault):
    ```bash
    aws-vault exec your-project-srv -- ./build.sh upload
    ```

    Another thing to note is that we are using `role/admin` here. Better would
    be to use the dedicated 'upload-role' that has been created using the
    `terraform-aws-mach-shared` module.<br>
    **todo** Describe how to set up these policies in a how-to article

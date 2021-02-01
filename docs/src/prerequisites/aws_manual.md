# Manual AWS setup

As described in the [tenenacy model](../guidance/tenancy.md#aws-tenancy), we advice to setup your MACH environment by creating **one service AWS account** containing shared resources and create an **AWS account per stack**.

This way, all resources are strictly seperated from eachother.


## Setup the service account

The most basic setup of the service account contains:

1. A Terraform state backend
2. A component registry
3. Route53 zone to route all other accounts from

### 1. Create S3 state backend
Create a S3 bucket which will be used as Terraform state backend.

!!! info "Setting up AWS"
    For more information on how to setup a S3 state backend including the correct IAM roles, see the [Terraform documentation](https://www.terraform.io/docs/backends/types/s3.html#s3-bucket-permissions)


### 2. Create component registry

This will be a S3 bucket should be created where the **packaged lambda** function code is stored.


!!! important "Component registry access"
      **All site-specific accounts** should have access to this bucket.

1. Create a S3 bucket, for example named `mach-lambda-repository` which is set to *private*
2. Set a bucket policy so that certain users and/or roles can upload to
   ```
   {
       "Version": "2012-10-17",
       "Statement": [
           {
               "Sid": "",
               "Effect": "Allow",
               "Principal": {
                   "AWS": [
                       "arn:aws:iam::1234567890:root",
                       "arn:aws:iam::1234567891:role/deploy"
                   ]
               },
               "Action": [
                   "s3:PutObject",
                   "s3:GetObject"
               ],
               "Resource": "arn:aws:s3:::mach-lambda-repository/*"
           }
       ]
   }
   ```
   Along the way, as new sites gets added, extra deploy roles need to be added as well, see [IAM deploy roles](#iam-deploy-role).


### 3. Setup Route53 zone

If a custom domain is going to be used for the API components a common practise is to have one (or more) top-level domain configured on the service account that can link to the site-specific Route53 zones using their NS records.

In the [example given below](#2-setup-route53-zone) a `my-company.com` and `internal-services.net` is needed on which other subdomains can be registrered on.
## Setup the site-specific account

These steps must be repeated **per site**

Start off with **createing a new account in AWS**.

Each account will contain at least the following:

1. A Terraform state backend
2. The Route53 hosted zones needed for the endpoints
3. A '*deploy*' IAM role for MACH to manage your resources

### 1. Create S3 state backend
Create a S3 bucket which will be used as Terraform state backend as well as the MACH state.

!!! info "Setting up AWS"
    For more information on how to setup a S3 state backend including the correct IAM roles, see the [Terraform documentation](https://www.terraform.io/docs/backends/types/s3.html#s3-bucket-permissions)

### 2. Setup Route53 zone

In case you are planning to deploy APIs that need custom routing, one or more Route53 zones needs to be configured.

If, for example, you need to have three differend endpoints in your MACH stack, say;

```yaml
 public: api.site-a.my-company.com
 management: manage.site-a.my-company.com
 internal: my-brand.internal-services.net
```

You'll need to setup two Route53 zones:

- `site-a.my-company.com`
- `my-brand.internal-services.net`

!!! note "API routings"
    MACH will make sure the API Gateway is created and a SSL certificate is created.<br>
    Each component is responsible for creating the correct routing to the Lambda endpoints.

    More on that in the [deployment section](../deployment/config/aws.md#custom-endpoint)

### 3. Create IAM deploy role

An IAM role should be created with which MACH can perform all necessary infra operations.

!!! tip ""
    For the sake of simplicity, name this IAM role `deploy`.

It should be possible for the main IAM user MACH is running is can assume this role.

#### Policies
Some examples of necessary actions that needs to be allowed on the deploy role are: [^1]

!!! warning ""
    Note that these are very simplified and loose policies used as an example. 
    In practise you might want to configure more strict policies per resource.

##### Access to code repository
```
statement {
   resources = [
      "arn:aws:s3:::mach-lambda-repository",
      "arn:aws:s3:::mach-lambda-repository/*",
   ]

   actions = [
      "s3:GetObject",
    ]
}

statement {
   resources = [
      "arn:aws:s3:::mach-lambda-repository",
   ]

   actions = [
      "s3:GetBucketLocation",
      "s3:ListBucket",
    ]
}
```

##### Public facing APIs

In case you have a public facing API:

```
statement {
    resources = ["*"]
    actions = [
      "route53:*",
      "apigateway:*",
      "acm:*",
    ]
}
```

##### Lambdas

```
statement {
    resources = ["*"]
    actions = [
      "lambda:*"
      "logs:*",
    ]
}
```

[^1]: Terraform syntax for a [`iam_policy_document`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) is used here.<br>
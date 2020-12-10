# Setting up AWS

!!! todo
    This section is not up-to-date yet

## Using Terraform

!!! tip "Example"
      See the [examples directory](https://github.com/labd/mach-composer/tree/master/examples/aws/infra/) for an example of a Terraform setup

## Manual setup

### Create S3 state backend
Create a S3 bucket which will be used as Terraform state backend.

!!! info "Setting up AWS"
    For more information on how to setup a S3 state backend including the correct IAM roles, see the [Terraform documentation](https://www.terraform.io/docs/backends/types/s3.html#s3-bucket-permissions)


### Create lambda repository

A S3 bucket should be created where the packaged lambda function code is stored.

**All site-specific accounts should have access to this bucket.**

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


### Setup AWS account per site

It is recommended to use a dedicated AWS account **per site**.<br>
This way, all resources are strictly seperated from eachother.

So the following steps need to be done per site:

1. Create an AWS account
2. Create the [Route53 hosted zones](#route53-zone) needed for the endpoints
3. Create a ['*deploy*' IAM role](#iam-deploy-role) for MACH to manage your resources

#### Route53 zone

In case you are planning to deploy APIs that need custom routing, one or more Route53 zones needs to be configured.

The full URL of this hosted zone can be configured in the MACH configuration.

!!! note "API routings"
    MACH will make sure the API Gateway is created and a SSL certificate is created.<br>
    Each component is responsible for creating the correct routing to the Lambda endpoints.

#### IAM deploy role

An IAM role should be created with which MACH can perform all necessary infra operations.

!!! tip ""
    For the sake of simplicity, name this IAM role `deploy`.

It should be possible for the main IAM user MACH is running is can assume this role.

##### Policies
Some examples of necessary actions that needs to be allowed on the deploy role are: [^1]

!!! warning ""
    Note that these are very simplified and loose policies used as an example. 
    In practise you might want to configure more strict policies per resource.

###### Access to code repository
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

###### Public facing APIs

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

###### Lambdas

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
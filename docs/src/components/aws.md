# AWS components

## AWS Lambda

AWS Lambda functions need to be uploaded to a S3 bucket. From there AWS Lambda will run the functions for you once instructed by the Terraform deployment.

### HTTP routing

MACH will provide the correct HTTP routing for you.  
To do so, the following has to be configured:

- [base_url](../syntax.md#sites) settings in the Site configuration
- The component has to be marked as [`has_public_api`](../syntax.md#components)

More information in the [deployment section](../deployment/aws.md).



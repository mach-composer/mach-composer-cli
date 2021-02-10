# Step 3. Add AWS services account

In AWS we need two accounts:

1. [**Service account**](#setup-service-account) for any shared resources amongst all MACH stacks
2. [**Site-specific account**](./step-4-setup-aws-site.md) for resources specific to a single MACH stack

In this step we'll create the first one, the service account.

!!! tip "Tenancy model"
      As described in the [tenenacy model](../../topics/architecture/tenancy.md#aws-tenancy), we advice to setup your MACH environment by creating **one service AWS account** containing shared resources and create an **AWS account per stack**.

      This way, all resources are strictly seperated from eachother.


## Setup service account

For this account we will create a;

1. **Terraform state backend** to store the infrastructure state
2. **Component registry** to deploy all MACH components to
3. **Route53 zone** to route all other site-specific accounts from

### 1. Create AWS account

- In your AWS console, go to **My Organization** and choose **Add accounts**
- For your new account choose a name like `yourproject-services` or `yourproject-shared`
  
!!! info "No root AWS account yet?"
    Go to [AWS support](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/) for instructions on how to setup your AWS root account.

### 2. Setup your Terraform configuration

### 3. Terraform infra Roll-out

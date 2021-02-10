# Step 4. Add AWS site account

Now that we've created our [service account](./step-3-setup-aws-services.md), we can create an AWS account for a specific site (or; MACH stack).

## Setup site-specific

For this account we will create:

1. Terraform state backend
2. The Route53 hosted zones needed for the endpoints
3. `deploy` IAM role for MACH to manage your resources

!!! Todo
    Describe steps
# Introduction

## How does it work?

MACH composer takes a [YAML configuration](../reference/syntax/index.md) as input, and will translate this into a Terraform configuration. It will then execute the terraform configuration, which will deploy all resources for the site architecture.

[![MACH diagram](../_img/mach.png)](../_img/mach.png)

MACH composer is intended for managing multiple instances of the architecture.

## Get started with MACH

In this tutorial we'll walk you through the steps required to get started with MACH.

- Step 1: [Install necessary tools](./step-1-installation.md)
- Step 2: [Setup your commercetools project](./step-2-setup-ct.md)
- For **AWS**:
    - Step 3: [Setup your AWS services account](./aws/step-3-setup-aws-services.md)
    - Step 4: [Setup your site-specific AWS account](./aws/step-4-setup-aws-site.md)
    - Step 5: [Create your first MACH component](./aws/step-5-create-component.md)
    - Step 6: [Setup and deploy your MACH stack](./aws/step-6-create-mach-stack.md)
- For **Azure**:
    - Step 3: [Setup your Azure environment](./azure/step-3-setup-azure.md)
    - Step 4: [Create your first MACH component](./azure/step-4-create-component.md)
    - Step 5: [Setup and deploy your MACH stack](./azure/step-5-create-mach-stack.md)
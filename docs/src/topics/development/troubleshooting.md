# Troubleshooting

## Fixing Terraform issues

In some cases you might encounter issues using Terraform like;

- A locked state file
- Having to import an existing resource into your state file or
- Having to remove a Terraform-managed resource from your state
- Having to move/rename a resource

For all of these cases, you can simply drop into the directory containing the
Terraform configuration and work with Terraform to perform the necessary
operations.

### Navigating to Terraform output

The easiest way to navigate to your Terraform output is by running

```bash
# in case of for example main.yml
cd deployments/main/<site-name>
```

From there you can start working with Terraform as usual.

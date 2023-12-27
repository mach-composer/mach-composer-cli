# Migration

Mach Composer does not support tools for migrating a component from site-managed
to site-component managed. However, it is possible to manually do the migrations
through using terraform directly.

To migrate a component from site-managed to site-component the following steps
are necessary:

1. Change the site-component configuration you want to move out by setting the
   correct deployment type:
   ```yaml
    # my-site.yml
   site:
     identifier: my-site
     components:
       - name: my-component
         deployment: 
            type: site-component
         # Etc...
   # Etc...
   ```
2. `mach-composer generate -f main.yaml` to generate the Terraform locally
3. `mach-composer init -f main.yaml` to initialize the Terraform CLI
4. Assuming a site named `my-site` move to generated `my-site`
   directory (`cd ./deployments/main/my-site`) and
   run `terraform state pull > terraform.tfstate`
5. Do the same for the component you want to move out, assuming a component
   named `my-component` move to generated `my-component` directory
   (`cd ./deployments/main/my-site/my-component`) and
   run `terraform state pull > terraform.tfstate`
6. Move back to the site directory (`cd ..`) and
   run `terraform state mv --state=terraform.tfstate --state-out=./my-component/terraform.tfstate module.my-component module.my-component`
7. Push the changes to the remote
   state (`terraform state push terraform.tfstate`)
8. Move to the component directory (`cd ./my-component`) and do the
   same: `terraform state push terraform.tfstate`
9. Move back to the root directory (`cd ../../..`) and
   run `terraform plan` to check if the changes will be applied correctly
10. Finally, run `terraform apply` to apply the changes

## Convenience script (use at own risk)

```bash
#!/bin/bash

if [ -n "$1" ]; then
  echo "Using $1 as the site directory."
else
  echo "Provide the site directory. This generally is in the format of ./deployments/main/<site-name>."
  exit 1
fi

if [ -n "$2" ]; then
  echo "Using $2 as the component"
else
  echo "Provide the component name"
  exit 1
fi


mach-composer generate -f main.yaml
mach-composer init -f main.yaml

cd ./$1 || exit

terraform state pull > terraform.tfstate

cd ./$2 || exit

terraform state pull > terraform.tfstate

cd ../

terraform state mv --state=terraform.tfstate --state-out=./"$2"/terraform.tfstate "module.$2" "module.$2"

terraform state push terraform.tfstate

cd ./"$2" || exit

terraform state push terraform.tfstate
```

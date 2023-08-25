# CLI Reference


MACH composer's command line interface allows you to perform the following actions:

```bash
$ mach-composer

Usage:
  mach-composer [command]

Available Commands:
  apply       Apply the configuration.
  cloud       Manage your Mach Composer Cloud
  completion  Generate the autocompletion script for the specified shell
  components  List all components.
  generate    Generate the Terraform files.
  help        Help about any command
  init        Initialize site directories Terraform files.
  plan        Plan the configuration.
  schema      Generate a JSON schema for your config based on the plugins.
  sites       List all sites.
  terraform   Execute terraform commands directly
  update      Update all (or a given) component.
  version     Return version information of the mach-composer cli

Flags:
  -h, --help      help for mach-composer
      --verbose   Verbose output.

Use "mach-composer [command] --help" for more information about a command.
```


## `apply`

Apply the configuration.

```bash
mach-composer apply --auto-approve -f main.yml
```

**Options**

- `--with-sp-login` If az login with service principal environment variables should be done.
- `--auto-approve` Auto-approve the Terraform plan
- `--file` or `-f TEXT` YAML file to apply. If not set apply all *.yml files.
- `--var-file` YAML file with variables to be used in the configuration file.
- `--site` or `-s TEXT` Site to apply. If not set apply all sites.
- `--component` or `-c TEXT` Specific component to target.
- `--output-path TEXT` Output path, defaults to `cwd`/deployments`.
- `--reuse` Suppress a terraform init for improved speed (not recommended for production usage)
- `--ignore-version` Skip MACH composer version check
- `--destroy` Destroy option is a convenient way to destroy all remote objects managed by this mach config



## `init`
Initialize the Terraform directory.

```bash
mach-composer terraform init -f main.yml
```

**Options**

- `--file` or `-f TEXT` YAML file to parse. If not set parse all *.yml files.
- `--var-file` YAML file with variables to be used in the configuration file.
- `--site` or `-s TEXT` Site to parse. If not set parse all sites.
- `--output-path TEXT` Output path, defaults to `cwd`/deployments.
- `--ignore-version` Skip MACH composer version check


## `plan`
Output the deployment plan.

```bash
mach-composer plan -f main.yml
```

**Options**

- `--file` or `-f TEXT` YAML file to parse. If not set parse all *.yml files.
- `--var-file` YAML file with variables to be used in the configuration file.
- `--site` or `-s TEXT` Site to generate plan of. If not set generate plans for all sites.
- `--component` or `-c TEXT` Specific component to target.
- `--output-path TEXT` Output path, defaults to `cwd`/deployments.
- `--reuse` Suppress a terraform init for improved speed (not recommended for production usage)
- `--ignore-version` Skip MACH composer version check
- `--destroy` Destroy option is a convenient way to destroy all remote objects managed by this mach config


## `update`

Usage: `mach-composer update [OPTIONS] [COMPONENT] [VERSION]`

Update all (or a given) component.

When no component and version is given, it will check the git repositories
for any updates. This command can also be used to manually update a single
component by specifying a component and version.

```bash
# To check for updates on all components
mach-composer update --check

# To update a specific component and create a commit message
mach-composer update pim-importer v1.2.0 -c
```

**Options**

- `--file` or `-f TEXT` YAML file to update. If not set update all *.yml files.
- `--verbose` or `-v` Verbose output.
- `--check` Only checks for updates, doesn't change files.
- `--commit` or `-c` Automatically commits the change.

## `cloud`

The `cloud` subcommand contains all the actions you can perform against [MACH composer ☁️](../cloud/index.md).

```bash
% mach-composer cloud     
Manage your Mach Composer Cloud

Usage:
  mach-composer cloud [flags]
  mach-composer cloud [command]

Available Commands:
  add-organization-user       Invite a user to a specific organization
  config                      Configure mach composer cloud
  create-api-client           Manage your components
  create-component            Register a new component
  create-organization         Create a new organization
  create-project              Create a new Project
  describe-component-versions List all version for an existing component
  list-api-clients            Manage your components
  list-component-versions     List all version for an existing component
  list-components             List your components
  list-organization-users     List all users in an organization
  list-organizations          List all organizations
  list-projects               List all Projects
  login                       Login to mach composer cloud
  register-component-version  Register a new version for an existing component

Flags:
  -h, --help   help for cloud

Global Flags:
      --verbose   Verbose output.

Use "mach-composer cloud [command] --help" for more information about a command.
```

!!! tip "Set default organization and project values"
  
    Most of the `cloud` commands require you to provide a `organization` and `project` parameter, so that it knows which entities to interact with. You set default settings for this as shown below.

    ```bash
    mach-composer cloud config set-organization my-org
    mach-composer cloud config set-project my-project
    ```


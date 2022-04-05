# CLI Reference


MACH composer's command line interface allows you to peform the following actions:

```bash
$ mach-composer

Usage:
  mach-composer [command]

Available Commands:
  apply       Apply the configuration.
  completion  Generate the autocompletion script for the specified shell
  generate    Generate the Terraform files.
  help        Help about any command
  init        Initialize site directories Terraform files.
  plan        Plan the configuration.
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
- `--reuse` Supress a terraform init for improved speed (not recommended for production usage)
- `--ignore-version` Skip MACH composer version check
- `--destroy` Destroy option is a convenient way to destroy all remote objects managed by this mach config



## `init`
Initiliaze the Terraform directory.

```bash
mach init -f main.yml
```

**Options**

- `--file` or `-f TEXT` YAML file to parse. If not set parse all *.yml files.
- `--var-file` YAML file with variables to be used in the configuration file.
- `--site` or `-s TEXT` Site to parse. If not set parse all sites.
- `--output-path TEXT` Output path, defaults to `cwd`/deployments.
- `--ignore-version` Skip MACH composer version check


## `plan`
Output the deploy plan.

```bash
mach plan -f main.yml
```

**Options**

- `--file` or `-f TEXT` YAML file to parse. If not set parse all *.yml files.
- `--var-file` YAML file with variables to be used in the configuration file.
- `--site` or `-s TEXT` Site to generate plan of. If not set generate plans for all sites.
- `--component` or `-c TEXT` Specific component to target.
- `--output-path TEXT` Output path, defaults to `cwd`/deployments.
- `--reuse` Supress a terraform init for improved speed (not recommended for production usage)
- `--ignore-version` Skip MACH composer version check
- `--destroy` Destroy option is a convenient way to destroy all remote objects managed by this mach config


## `update`

Usage: `mach update [OPTIONS] [COMPONENT] [VERSION]`

Update all (or a given) component.

When no component and version is given, it will check the git repositories
for any updates. This command can also be used to manually update a single
component by specifying a component and version.

```bash
# To check for updates on all components
mach update --check

# To update a specific component and create a commit message
mach update pim-importer v1.2.0 -c
```

**Options**

- `--file` or `-f TEXT` YAML file to update. If not set update all *.yml files.
- `--verbose` or `-v` Verbose output.
- `--check` Only checks for updates, doesnt change files.
- `--commit` or `-c` Automatically commits the change.


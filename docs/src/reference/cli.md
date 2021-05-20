# CLI Reference

MACH composer's command line interface allows you to peform the following actions:

```bash
$ mach
Usage: mach [OPTIONS] COMMAND [ARGS]...

Options:
  --help  Show this message and exit.

Commands:
  apply       Apply the configuration.
  bootstrap   Bootstraps a configuration or component.
  components  List all components.
  generate    Generate the Terraform files.
  plan        Output the deploy plan.
  sites       List all sites.
  update      Update all (or a given) component.
```


## `apply`

Apply the configuration.

```bash
mach apply --auto-approve -f main.yml
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


## `bootstrap`

Usage: `mach bootstrap [OPTIONS] [config|component]`

Bootstraps a configuration or component.

```bash
# To start a wizard to create a new configuration file
mach bootstrap config

# To start a wizard to create a new component
mach bootstrap component
```

**Options**

- `--output` or `-o TEXT` Output file or directory.
- `--cookiecutter` or `-c TEXT` cookiecutter repository to generate from.


## `components`
List all components.

```bash
mach components -f main.yml
```

**Options**

- `--file` or `-f TEXT` YAML file to read. If not set parse all *.yml files.

## `generate`
Generate the Terraform files.

```bash
mach generate -f main.yml
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


## `sites`
List all sites.

```bash
mach sites -f main.yml
```

**Options**

- `--file` or `-f TEXT` YAML file to read. If not set parse all *.yml files.

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

# You can also use the "@" delimiter for component version
mach update pim-importer@v1.2.0 -c
```

**Options**

- `--file` or `-f TEXT` YAML file to update. If not set update all *.yml files.
- `--verbose` or `-v` Verbose output.
- `--check` Only checks for updates, doesnt change files.
- `--commit` or `-c` Automatically commits the change.


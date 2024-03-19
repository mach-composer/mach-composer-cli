## mach-composer init

Initialize site directories Terraform files.

```
mach-composer init [flags]
```

### Options

```
  -f, --file string          YAML file to parse. (default "main.yml")
      --force                Force the apply to run even if the components are considered up to date
  -h, --help                 help for init
      --ignore-version       Skip MACH composer version check
      --output-path string   Outputs path to store the generated files. (default "deployments")
  -s, --site string          Site to parse. If not set parse all sites.
      --var-file string      Use a variable file to parse the configuration with.
  -w, --workers int          The number of workers to use (default 1)
```

### Options inherited from parent commands

```
      --verbose   Verbose output.
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


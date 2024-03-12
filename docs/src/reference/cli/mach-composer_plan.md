## mach-composer plan

Plan the configuration.

```
mach-composer plan [flags]
```

### Options

```
  -c, --component stringArray     
  -f, --file string               YAML file to parse. (default "main.yml")
  -h, --help                      help for plan
      --ignore-change-detection   Ignore change detection to run even if the components are considered up to date
      --ignore-version            Skip MACH composer version check
      --lock                      Acquire a lock on the state file before running terraform plan (default true)
      --output-path string        Outputs path to store the generated files. (default "deployments")
      --reuse                     Suppress a terraform init for improved speed (not recommended for production usage)
  -s, --site string               Site to parse. If not set parse all sites.
      --var-file string           Use a variable file to parse the configuration with.
  -w, --workers int               The number of workers to use (default 1)
```

### Options inherited from parent commands

```
      --verbose   Verbose output.
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


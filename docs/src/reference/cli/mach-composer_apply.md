## mach-composer apply

Apply the configuration.

```
mach-composer apply [flags]
```

### Options

```
      --auto-approve            Suppress a terraform init for improved speed (not recommended for production usage)
  -c, --component stringArray   
      --destroy                 Destroy option is a convenient way to destroy all remote objects managed by this mach config
  -f, --file string             YAML file to parse. (default "main.yml")
  -h, --help                    help for apply
      --ignore-version          Skip MACH composer version check
      --output-path string      Outputs path to store the generated files. (default "deployments")
      --reuse                   Suppress a terraform init for improved speed (not recommended for production usage)
  -s, --site string             Site to parse. If not set parse all sites.
      --var-file string         Use a variable file to parse the configuration with.
  -w, --workers int             The number of workers to use (default 1)
```

### Options inherited from parent commands

```
      --verbose   Verbose output.
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


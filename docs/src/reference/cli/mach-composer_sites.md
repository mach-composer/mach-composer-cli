## mach-composer sites

List all sites.

```
mach-composer sites [flags]
```

### Options

```
  -f, --file string          YAML file to parse. (default "main.yml")
  -h, --help                 help for sites
      --ignore-version       Skip MACH composer version check
      --output-path string   Outputs path to store the generated files. (default "deployments")
  -s, --site string          Site to parse. If not set parse all sites.
      --var-file string      Use a variable file to parse the configuration with.
  -w, --workers int          The number of workers to use (default 1)
```

### Options inherited from parent commands

```
  -q, --quiet     Quiet output. This is equal to setting log levels to error and higher
  -v, --verbose   Verbose output. This is equal to setting log levels to debug and higher
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


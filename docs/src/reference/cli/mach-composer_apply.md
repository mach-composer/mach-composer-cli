## mach-composer apply

Apply the configuration.

```
mach-composer apply [flags]
```

### Options

```
      --auto-approve              Suppress a terraform init for improved speed (not recommended for production usage)
  -b, --buffer                    Whether logs should be buffered and printed at the end of the run
  -c, --component stringArray     
      --destroy                   Destroy option is a convenient way to destroy all remote objects managed by this mach config
  -f, --file string               YAML file to parse. (default "main.yml")
      --force-init                Force terraform initialization. By default mach-composer will reuse existing terraform resources
  -g, --github                    Whether logs should be decorated with github-specific formatting
  -h, --help                      help for apply
      --ignore-change-detection   Ignore change detection to run even if the components are considered up to date
      --ignore-version            Skip MACH composer version check
      --output-path string        Outputs path to store the generated files. (default "deployments")
  -s, --site string               Site to parse. If not set parse all sites.
      --var-file string           Use a variable file to parse the configuration with.
  -w, --workers int               The number of workers to use (default 1)
```

### Options inherited from parent commands

```
      --output string   The output type. One of: console, json (default "console")
  -q, --quiet           Quiet output. This is equal to setting log levels to error and higher
      --strip-logs      Strip all context from the logs
  -v, --verbose         Verbose output. This is equal to setting log levels to debug and higher
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


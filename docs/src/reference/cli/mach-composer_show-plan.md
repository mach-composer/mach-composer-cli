## mach-composer show-plan

Show the planned configuration.

```
mach-composer show-plan [flags]
```

### Options

```
  -b, --buffer                    Whether logs should be buffered and printed at the end of the run
  -f, --file string               YAML file to parse. (default "main.yml")
      --force-init                Force terraform initialization. By default mach-composer will reuse existing terraform resources
  -g, --github                    Whether logs should be decorated with github-specific formatting
  -h, --help                      help for show-plan
      --ignore-change-detection   Ignore change detection to run even if the components are considered up to date
      --ignore-version            Skip MACH composer version check
      --no-color                  Disable color output
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


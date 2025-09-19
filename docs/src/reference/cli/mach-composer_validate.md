## mach-composer validate

Validate the generated terraform configuration. See [the documentation](/howto/cli/filtering-commands) for filtering options.

### Synopsis

This command validates the generated terraform configuration. It will check the provided configuration file for any errors, and will run `terraform validate` on the generated configuration. This will check for any syntax errors in the generated configuration without accessing the actual infrastructure.

By default, the generated configuration is stored in the `validations` directory in the current working directory. This can be changed by providing the `--validation-path` flag.

See [the terraform validation docs](https://www.terraform.io/docs/commands/validate.html) for more information on `terraform validate`.

```
mach-composer validate [flags]
```

### Options

```
  -b, --buffer                   Whether logs should be buffered and printed at the end of the run
  -f, --file string              YAML file to parse. (default "main.yml")
      --filter stringArray       Run only nodes matching the filter expression
  -g, --github                   Whether logs should be decorated with github-specific formatting
  -h, --help                     help for validate
      --ignore-version           Skip MACH composer version check
      --output-path string       Outputs path to store the generated files. (default "deployments")
  -s, --site string              Site to parse. If not set parse all sites.
      --validation-path string   Directory path to store files required for configuration validation. (default "validations")
      --var-file string          Use a variable file to parse the configuration with.
  -w, --workers int              The number of workers to use (default 1)
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


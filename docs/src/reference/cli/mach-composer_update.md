## mach-composer update

Update all (or a given) component.

```
mach-composer update [flags]
```

### Options

```
      --check                   Only checks for updates, doesnt change files.
      --cloud                   Use MACH composer cloud to check for updates.
  -c, --commit                  Automatically commits the change.
  -m, --commit-message string   Use a custom message for the commit.
      --component stringArray   
  -f, --file string             YAML file to update. (default "main.yml")
  -h, --help                    help for update
```

### Options inherited from parent commands

```
  -q, --quiet     Quiet output. This is equal to setting log levels to error and higher
  -v, --verbose   Verbose output. This is equal to setting log levels to debug and higher
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


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
      --output string   The output type. One of: console, json (default "console")
      --verbose         Verbose output.
```

### SEE ALSO

* [mach-composer](mach-composer.md)	 - MACH composer is an orchestration tool for modern MACH ecosystems


# Development environment

This section describes some tips on how to configure your development
environment so you can use MACH Composer to the fullest.

## Installing the CLI

MACH Composer is written in Go. For macOs and Linux users the easiest way to
install Mach Composer is via brew:

```bash
brew tap labd/mach-composer
brew install mach-composer
```

## Using the JSON schema for IntelliSense autocompletion

### On Visual Studio Code

Register the MACH Composer schema per project by adding a .vscode/settings.json with the
following configuration:
```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/labd/mach-composer/main/internal/config/schemas/schema-1.yaml": "*.yml"
  }
}
```

# Development environment

This section describes some tips on how to configure your development
environment so you can use MACH composer to the fullest.

## Installing the CLI

MACH composer is written in Go. For macOS and Linux users the easiest way to
install MACH composer is via brew:

```bash
brew tap mach-composer/mach-composer
brew install mach-composer
```

## Using the JSON schema for IntelliSense autocompletion

### On Visual Studio Code

Register the MACH composer schema per project by adding a .vscode/settings.json with the
following configuration:
```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/mach-composer/mach-composer-cli/main/internal/config/schemas/schema-1.yaml": "*.yml"
  }
}
```

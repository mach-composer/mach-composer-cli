# Development environment

This section describes some tips on how to configure your development environment so you can use MACH composer to the fullest.

## Installing the CLI

MACH composer is a Python project. In order to make this work on your local machine, make sure you have **Python 3.8** installed on your system.

You can install MACH directly from the [Python package index](https://pypi.org/project/mach-composer/). 

!!! tip "Using pipx"
    We recommend using [pipx](https://pipxproject.github.io/pipx/) as an installer.<br>
    This will make sure MACH composer gets installed in it's own sandboxed virtual environment and lets you run `mach` from your command line.

```bash
$ pipx install mach-composer
```

## Using the JSON schema for IntelliSense autocompletion

### On Visual Studio Code

Register the MACH schema per project by adding a .vscode/settings.json with the following configuration:
```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/labd/mach-composer/master/schema.json": "*.yml"
  }
}
```
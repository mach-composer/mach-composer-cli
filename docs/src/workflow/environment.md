# Development environment

This section describes some tips on how to configure your development environment so you can use the MACH composer to the fullest.

## Installing the CLI

The MACH composer is a Python project. In order to make this work on your local machine, make sure you have **Python 3.8** installed on your system.

- [ ] TODO: Describe installation steps

## Using the JSON schema

### On Visual Studio Code

Register the MACH schema per project by adding a .vscode/settings.json with the following configuration:
```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/labd/mach-composer/master/schema.json": "*.yml"
  }
}
```

!!! warning "Almost public.."
    At the moment, the mach-composer repository is not public yet, so directly refering to a file without the use of a (temporary) token is not possible.  
    For now, you can download the schema and place it in your configuration repo.
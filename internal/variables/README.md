# Variables

This package providers variables support. Currently supported variable types
are:
 - ${var.foo} -> Will use variables from a file
 - ${env.foo} -> Will use the environment variable FOO
 - ${component.foo.bar} -> Will use output from other component


The `var` and `env` vars are processed during the config parsing phase and
replaced before being passed to the generator. The `component` var is
interpreted during the rendering phase.

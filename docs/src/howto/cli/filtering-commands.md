# Filtering commands

> For site-component targeting with `--filter` the deployment type of the component must be marked as `site-component`!

By default, mach composer cli will run all nodes in the configuration in the order inferred from the configuration. You
can filter which nodes to run by using one or more `--filter` flags on the commands that support it. With this you can
more selectively target specific nodes in your configuration.

The filter expression is a simple string that is matched against the node names:

- Run a component: `mach-composer apply --filter component-1`
- Run a site: `mach-composer apply --filter site-1`
- Run a specific component in a specific site: `mach-composer apply --filter site-1/component-1`

Some additional microsyntaxes are available to further define which nodes to run:

- `site-1...` will run site-1 and all dependent components
- `site-1/component-1...` will run component-1 in site-1 and all its dependencies

Syntaxes can also be combined by using multiple `--filter` flags. Note that when multiple filters are provided that
could match the same node, priority is given to the most generic filter. This means that any less generic filters are
ignored:

- Ignores `site-1/component-1` because `site-1...` already matches it:
  `mach-composer apply --filter site-1... --filter site-1/component-1`




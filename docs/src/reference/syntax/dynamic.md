Based on loaded plugins in the `mach_composer` block, more configuration
blocks might be loaded here. For example, if the `sentry` plugin is loaded a
block named `sentry` will be loaded here. Refer to the
[plugin documentation](../../plugins/index.md)
to see what configuration blocks will be loaded.

```yaml
sentry:
  auth_token: <your-sentry-auth-token>
  organization: <your-sentry-organization>
  project: <your-sentry-project>
  rate_limit_window: <your-sentry-rate-limit-window>
  rate_limit_count: <your-sentry-rate-limit-count>
```

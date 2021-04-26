# Step 3. Prepare you Google Cloud environment

In Google Cloud you need to create two new [**Google Cloud projects**](https://cloud.google.com/resource-manager/docs/creating-managing-projects)

1. **Site-specific project** - for resources specific to a single MACH stack
2. **Shared project** - for any shared resources amongst all MACH stacks


## 1. Prepare

Login to your Google subscription through the CLI of [gcloud](https://cloud.google.com/sdk/gcloud/).

```bash
gcloud auth login
```
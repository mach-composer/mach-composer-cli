# Step 1. Installation

## Install Terraform

Make sure you have Terraform installed on your machine.

We recommend using [tfenv](https://github.com/tfutils/tfenv) to easily switch
between Terraform versions when needed.

1. Install [tfenv](https://github.com/tfutils/tfenv): Follow the instructions on [https://github.com/tfutils/tfenv](https://github.com/tfutils/tfenv)
2. Make sure you have Terraform installed
   ```bash
   $ tfenv install latest
   ```

## Install MACH composer
The easiest way to get started is by installing MACH composer locally to be able
to perform a MACH composer deploy and make use of the
[developer tools](../topics/development/workflow.md).

```bash
brew tap labd/mach-composer
brew install mach-composer
```

!!! tip "Next: step 2"
    Next we'll [setup our commercetools project](./step-2-setup-ct.md)

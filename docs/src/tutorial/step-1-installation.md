# Step 1. Installation

## Install MACH composer
The easiest way to get started is by installing MACH composer locally to be able
to perform a MACH composer deploy and make use of the
[developer tools](../topics/development/workflow.md).


### MacOS and linux
```bash
brew tap labd/mach-composer
brew install mach-composer
```

### Windows
```ps
choco install mach-composer
```

### Alternatives
You can always download the latest release via the [GitHub release page](https://github.com/labd/mach-composer/releases)


## Install Terraform
Make sure you have Terraform installed on your machine.

On MacOS and Linux we recommend using [tfenv](https://github.com/tfutils/tfenv)
to easily switch between Terraform versions when needed.

1. Install [tfenv](https://github.com/tfutils/tfenv): Follow the instructions on [https://github.com/tfutils/tfenv](https://github.com/tfutils/tfenv)
2. Make sure you have Terraform installed
   ```bash
   $ tfenv install latest
   ```

For windows you can install it via chocolatey.

```ps
choco install terraform
```


!!! tip "Next: step 2"
    Next we'll [setup our commercetools project](./step-2-setup-ct.md)

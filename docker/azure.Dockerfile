#syntax=docker/dockerfile:1.4.0

FROM base

ARG TERRAFORM_AZURE_VERSION=2.99.0
ARG AZURE_CLI_VERSION=2.34.1

# For Azure
RUN apk add --no-cache python3-dev py3-pip py3-bcrypt py3-pynacl

# Update pip so that we can install a wheel of cryptography
RUN python3 -m pip install --upgrade pip

# Install Azure CLI
RUN pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}

# Install azure provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-azurerm/${TERRAFORM_AZURE_VERSION}/terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*


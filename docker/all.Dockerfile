#syntax=docker/dockerfile:1.4.0

FROM base

ARG TERRAFORM_AWS_VERSION=3.74.1
ARG TERRAFORM_AZURE_VERSION=2.99.0
ARG TERRAFORM_GOOGLE_VERSION=4.83.0
ARG AZURE_CLI_VERSION=2.34.1

# Azure
RUN apk add --no-cache python3-dev py3-pip py3-bcrypt py3-pynacl

# Update pip so that we can install a wheel of cryptography
RUN python3 -m pip install --upgrade pip
RUN pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-azurerm/${TERRAFORM_AZURE_VERSION}/terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*


# AWS
RUN apk add --no-cache --update aws-cli
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-aws/${TERRAFORM_AWS_VERSION}/terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*

# Install google provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-google/${TERRAFORM_GOOGLE_VERSION}/terraform-provider-google_${TERRAFORM_GOOGLE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-google_${TERRAFORM_GOOGLE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*


FROM goreleaser/goreleaser:v1.7.0 AS builder

COPY . /code/
WORKDIR /code/
RUN goreleaser build --single-target --skip-validate

FROM alpine:3.14

ENV SOPS_VERSION=3.7.2
ENV AZURE_CLI_VERSION=2.34.1
ENV TERRAFORM_VERSION=1.1.3
ENV TERRAFORM_EXTERNAL_VERSION=2.2.2
ENV TERRAFORM_AZURE_VERSION=2.86.0
ENV TERRAFORM_AWS_VERSION=3.66.0
ENV TERRAFORM_NULL_VERSION=2.1.2
ENV TERRAFORM_COMMERCETOOLS_VERSION=1.0.0-pre.2
ENV TERRAFORM_CONTENTFUL_VERSION=0.1.0
ENV TERRAFORM_AMPLIENCE_VERSION=0.3.7
ENV TERRAFORM_SENTRY_VERSION=0.7.0

RUN apk add --no-cache --virtual .build-deps g++ libffi-dev openssl-dev wget unzip make \
    && apk add bash ca-certificates git libc6-compat openssl openssh-client jq curl

# For AWS
RUN apk add --update aws-cli

# For Azure
RUN apk add python3-dev py3-pip py3-bcrypt py3-pynacl

# Update pip so that we can install a wheel of cryptography
RUN python3 -m pip install --upgrade pip

# Install Azure CLI
RUN pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}

# Install SOPS
RUN cd /tmp && \
    wget https://github.com/mozilla/sops/releases/download/v${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux.amd64 && \
    /usr/bin/install sops-v3.7.2.linux.amd64 /usr/local/bin/sops

# Install terraform
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip -n terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Create non-root user
RUN adduser mach-composer \
    --disabled-password

RUN mkdir /code /deployments && chown mach-composer /code /deployments

USER mach-composer
WORKDIR /home/mach-composer

# Pre-install Terreform plugins
ENV TF_PLUGIN_CACHE_DIR=/home/mach-composer/.terraform.d/plugin-cache
ENV TERRAFORM_PLUGINS_PATH=/home/mach-composer/.terraform.d/plugins/linux_amd64
RUN mkdir -p ${TF_PLUGIN_CACHE_DIR}
RUN mkdir -p ${TERRAFORM_PLUGINS_PATH}

# Install null provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-null/${TERRAFORM_NULL_VERSION}/terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install external provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-external/${TERRAFORM_EXTERNAL_VERSION}/terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install aws provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-aws/${TERRAFORM_AWS_VERSION}/terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install azure provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-azurerm/${TERRAFORM_AZURE_VERSION}/terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install commercetools provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-commercetools/releases/download/v${TERRAFORM_COMMERCETOOLS_VERSION}/terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install contentful provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-contentful/releases/download/v${TERRAFORM_CONTENTFUL_VERSION}/terraform-provider-contentful_${TERRAFORM_CONTENTFUL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-contentful_${TERRAFORM_CONTENTFUL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install amplience provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-amplience/releases/download/v${TERRAFORM_AMPLIENCE_VERSION}/terraform-provider-amplience_${TERRAFORM_AMPLIENCE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-amplience_${TERRAFORM_AMPLIENCE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install sentry provider
RUN cd /tmp && \
    wget https://github.com/jianyuan/terraform-provider-sentry/releases/download/v${TERRAFORM_SENTRY_VERSION}/terraform-provider-sentry_${TERRAFORM_SENTRY_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-sentry_${TERRAFORM_SENTRY_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

COPY --from=builder /code/dist/mach-composer_linux_amd64/mach-composer /usr/local/bin

ENTRYPOINT ["mach-composer"]

#syntax=docker/dockerfile:1.4.0

FROM alpine:3.14 AS base

ARG TERRAFORM_VERSION=1.3.5

ENV SOPS_VERSION=3.7.2
ENV TERRAFORM_EXTERNAL_VERSION=2.2.2
ENV TERRAFORM_NULL_VERSION=2.1.2
ENV TERRAFORM_COMMERCETOOLS_VERSION=0.30.0
ENV TERRAFORM_CONTENTFUL_VERSION=0.1.0
ENV TERRAFORM_AMPLIENCE_VERSION=0.3.7
ENV TERRAFORM_SENTRY_VERSION=0.7.0

RUN apk add --no-cache --virtual .build-deps g++ libffi-dev openssl-dev wget unzip make \
    && apk add bash ca-certificates git libc6-compat openssl openssh-client jq curl

RUN mkdir /code /deployments


# Install SOPS
RUN cd /tmp && \
    wget https://github.com/mozilla/sops/releases/download/v${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux.amd64 && \
    /usr/bin/install sops-v3.7.2.linux.amd64 /usr/local/bin/sops && \
    rm -rf /tmp/*

# Install terraform
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip -n terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin && \
    rm -rf /tmp/*

# Pre-install Terreform plugins
ENV TF_PLUGIN_CACHE_DIR=/home/mach-composer/.terraform.d/plugin-cache
ENV TERRAFORM_PLUGINS_PATH=/home/mach-composer/.terraform.d/plugins/linux_amd64
RUN mkdir -p ${TF_PLUGIN_CACHE_DIR}
RUN mkdir -p ${TERRAFORM_PLUGINS_PATH}

# Install null provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-null/${TERRAFORM_NULL_VERSION}/terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*

# Install external provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-external/${TERRAFORM_EXTERNAL_VERSION}/terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*

# Install commercetools provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-commercetools/releases/download/v${TERRAFORM_COMMERCETOOLS_VERSION}/terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*

# Install contentful provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-contentful/releases/download/v${TERRAFORM_CONTENTFUL_VERSION}/terraform-provider-contentful_${TERRAFORM_CONTENTFUL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-contentful_${TERRAFORM_CONTENTFUL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*

# Install amplience provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-amplience/releases/download/v${TERRAFORM_AMPLIENCE_VERSION}/terraform-provider-amplience_${TERRAFORM_AMPLIENCE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-amplience_${TERRAFORM_AMPLIENCE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*

# Install sentry provider
RUN cd /tmp && \
    wget https://github.com/jianyuan/terraform-provider-sentry/releases/download/v${TERRAFORM_SENTRY_VERSION}/terraform-provider-sentry_${TERRAFORM_SENTRY_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-sentry_${TERRAFORM_SENTRY_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*


COPY . /mach-composer
RUN ln -s /mach-composer/dist/mach-composer_linux_amd64_v1/bin/mach-composer /usr/local/bin/mach-composer

ENTRYPOINT ["mach-composer"]

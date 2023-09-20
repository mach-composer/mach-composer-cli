#syntax=docker/dockerfile:1.4.0

FROM base

ARG TERRAFORM_AWS_VERSION=3.76.1

RUN echo $TERRAFORM_AWS_VERSION

RUN apk add --no-cache --update aws-cli

# Install aws provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-aws/${TERRAFORM_AWS_VERSION}/terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*
